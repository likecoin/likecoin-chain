package deposit

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/likecoin/likechain/abci/query"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"

	"github.com/likecoin/likechain/services/eth"
	logger "github.com/likecoin/likechain/services/log"
	"github.com/likecoin/likechain/services/utils"
)

const blockBatchSize = 10000

var (
	log         = logger.L
	proposeLock = &sync.Mutex{}
)

func fillSig(tx *txs.DepositTransaction, privKey *ecdsa.PrivateKey) {
	tx.Proposal.Sort()
	jsonMap := tx.GenerateJSONMap()
	hash, err := jsonMap.Hash()
	if err != nil {
		log.
			WithField("tx", tx).
			WithError(err).
			Panic("Cannot hash deposit transaction")
	}
	sig, err := crypto.Sign(hash, privKey)
	if err != nil {
		log.
			WithField("tx", tx).
			WithError(err).
			Panic("Cannot sign deposit transaction")
	}
	sig[64] += 27
	jsonSig := txs.DepositJSONSignature{}
	copy(jsonSig.JSONSignature[:], sig)
	tx.Sig = &jsonSig
}

func propose(tmClient *tmRPC.HTTP, tmPrivKey *ecdsa.PrivateKey, proposal deposit.Proposal) {
	proposeLock.Lock()
	defer proposeLock.Unlock()
	log.
		WithField("block_number", proposal.BlockNumber).
		WithField("event_count", len(proposal.Inputs)).
		Info("Proposing new proposal")
	ethAddr := crypto.PubkeyToAddress(tmPrivKey.PublicKey)
	addr, err := types.NewAddress(ethAddr[:])
	if err != nil {
		log.
			WithField("eth_addr", ethAddr.Hex()).
			WithError(err).
			Panic("Cannot convert Ethereum address to LikeChain address")
	}
	queryResult, err := tmClient.ABCIQuery("account_info", []byte(addr.String()))
	if err != nil {
		log.
			WithField("addr", addr.String()).
			WithError(err).
			Panic("Cannot query account info from ABCI")
	}
	accInfo := query.GetAccountInfoRes(queryResult.Response.Value)
	if accInfo == nil {
		log.
			WithField("account_info_res", string(queryResult.Response.Value)).
			WithField("account_info_res_raw", queryResult.Response.Value).
			Panic("Cannot parse account info result")
	}
	log.
		WithField("nonce", accInfo.NextNonce).
		Debug("Got account info")
	tx := &txs.DepositTransaction{
		Proposer: addr,
		Proposal: proposal,
		Nonce:    accInfo.NextNonce,
	}
	fillSig(tx, tmPrivKey)
	rawTx := txs.EncodeTx(tx)
	txHash := tmhash.Sum(rawTx)
	txResult, err := tmClient.Tx(txHash, false)
	if err == nil {
		log.
			WithField("tx_hash", common.Bytes2Hex(txHash)).
			WithField("tx_height", txResult.Height).
			Info("Deposit tx is already processed, skipping")
		return
	}
	log.
		WithField("raw_tx", common.Bytes2Hex(rawTx)).
		Debug("Broadcasting transaction onto LikeChain")
	result, err := tmClient.BroadcastTxCommit(rawTx)
	if err != nil {
		log.
			WithField("raw_tx", common.Bytes2Hex(rawTx)).
			WithError(err).
			Panic("Broadcast transaction onto LikeChain failed")
	}
	if result.CheckTx.Code != response.Success.Code {
		switch result.CheckTx.Code {
		case response.DepositDoubleApproval.ToResponseCheckTx().Code:
			fallthrough
		case response.DepositAlreadyExecuted.ToResponseCheckTx().Code:
			log.
				WithField("code", result.CheckTx.Code).
				WithField("info", result.CheckTx.Info).
				WithField("log", result.CheckTx.Log).
				Info("Deposit transaction unnecessary and rejected in CheckTx, skipping")
		default:
			log.
				WithField("code", result.CheckTx.Code).
				WithField("info", result.CheckTx.Info).
				WithField("log", result.CheckTx.Log).
				Panic("Deposit transaction failed in CheckTx")
		}
	} else if result.DeliverTx.Code != response.Success.Code {
		switch result.DeliverTx.Code {
		case response.DepositDoubleApproval.ToResponseDeliverTx().Code:
			fallthrough
		case response.DepositAlreadyExecuted.ToResponseDeliverTx().Code:
			log.
				WithField("code", result.DeliverTx.Code).
				WithField("info", result.DeliverTx.Info).
				WithField("log", result.DeliverTx.Log).
				Info("Deposit transaction unnecessary and rejected in DeliverTx, skipping")
		default:
			log.
				WithField("code", result.DeliverTx.Code).
				WithField("info", result.DeliverTx.Info).
				WithField("log", result.DeliverTx.Log).
				Panic("Deposit transaction failed in DeliverTx")
		}
	} else {
		log.Info("Successfully broadcasted deposit transaction onto LikeChain")
	}
}

type runState struct {
	LastEthBlock       int64
	PendingBlockRanges [][]int64
	lock               *sync.Mutex
}

func loadState(path string) (*runState, error) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	state := runState{}
	err = json.Unmarshal(jsonBytes, &state)
	if err != nil {
		return nil, err
	}
	state.lock = &sync.Mutex{}
	return &state, nil
}

func (state *runState) save(path string) error {
	jsonBytes, err := json.Marshal(&state)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, jsonBytes, 0644)
	return err
}

func scanAndProposeForRange(config *Config, from, to uint64) {
	log.
		WithField("begin_block", from).
		WithField("end_block", to).
		Debug("Searching blocks in range")
	proposals := eth.GetTransfersFromBlocks(
		config.LoadBalancer,
		config.TokenAddr,
		config.RelayAddr,
		uint64(from),
		uint64(to),
	)
	if len(proposals) == 0 {
		log.
			WithField("begin_block", from).
			WithField("end_block", to).
			Info("No transfer events in range")
	} else {
		for _, proposal := range proposals {
			log.
				WithField("block", proposal.BlockNumber).
				Info("Proposing proposal")
			utils.RetryIfPanic(5, func() {
				propose(config.TMClient, config.TMPrivKey, proposal)
			})
		}
	}
}

// Config is the configuration about deposit
type Config struct {
	TMClient       *tmRPC.HTTP
	LoadBalancer   *eth.LoadBalancer
	TokenAddr      common.Address
	RelayAddr      common.Address
	TMPrivKey      *ecdsa.PrivateKey
	BlockDelay     int64
	StatePath      string
	StartFromBlock int64
	HTTPLogHook    *logger.HTTPHook
}

// Run starts the subscription to the deposits on Ethereum into the relay contract and commits proposal onto LikeChain
func Run(config *Config) {
	httpHookCleanupFunc := func() {
		if config.HTTPLogHook != nil {
			config.HTTPLogHook.Cleanup()
		}
	}
	defer httpHookCleanupFunc()
	state, err := loadState(config.StatePath)
	blockNumber := eth.GetHeight(config.LoadBalancer)
	if err != nil {
		log.
			WithField("state_path", config.StatePath).
			WithError(err).
			Info("Failed to load state, creating empty state")
		state = &runState{lock: &sync.Mutex{}}
		state.LastEthBlock = blockNumber
		if config.StartFromBlock > 0 && blockNumber >= config.StartFromBlock+config.BlockDelay {
			state.PendingBlockRanges = [][]int64{{config.StartFromBlock, blockNumber - config.BlockDelay}}
		}
	}
	if state.LastEthBlock < blockNumber {
		state.PendingBlockRanges = append(
			state.PendingBlockRanges,
			[]int64{state.LastEthBlock - config.BlockDelay + 1, blockNumber - config.BlockDelay},
		)
	}
	state.LastEthBlock = blockNumber
	state.save(config.StatePath)
	go func() {
		defer httpHookCleanupFunc()
		log.
			WithField("ranges", state.PendingBlockRanges).
			Info("Clearing pending block ranges previously left and accumulated during service halt")
		for len(state.PendingBlockRanges) > 0 {
			i := len(state.PendingBlockRanges) - 1
			start := state.PendingBlockRanges[i][0]
			end := state.PendingBlockRanges[i][1]
			for to := end; to >= start; to -= blockBatchSize {
				from := to - blockBatchSize + 1
				if from < start {
					from = start
				}
				scanAndProposeForRange(config, uint64(from), uint64(to))
				func() {
					state.lock.Lock()
					defer state.lock.Unlock()
					if from == start {
						state.PendingBlockRanges = state.PendingBlockRanges[:i]
					} else {
						state.PendingBlockRanges[i][1] = from - 1
					}
					state.save(config.StatePath)
				}()
			}
		}
		log.Info("All pending blocks cleared")
	}()
	eth.SubscribeHeader(config.LoadBalancer, func(header *ethTypes.Header) bool {
		newBlock := header.Number.Int64()
		if newBlock <= 0 {
			return true
		}
		if newBlock < config.BlockDelay {
			return true
		}
		log.
			WithField("last_block", state.LastEthBlock).
			WithField("new_block", newBlock).
			Info("Received new Ethereum block")
		for from := state.LastEthBlock - config.BlockDelay + 1; from <= newBlock-config.BlockDelay; from += blockBatchSize {
			to := from + blockBatchSize - 1
			if to > newBlock-config.BlockDelay {
				to = newBlock - config.BlockDelay
			}
			scanAndProposeForRange(config, uint64(from), uint64(to))
			func() {
				state.lock.Lock()
				defer state.lock.Unlock()
				state.LastEthBlock = to + config.BlockDelay
				state.save(config.StatePath)
			}()
		}
		return true
	})
}
