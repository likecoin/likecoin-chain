package deposit

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"sync"

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
	"github.com/likecoin/likechain/services/tendermint"
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
	log.
		WithField("raw_tx", common.Bytes2Hex(rawTx)).
		Debug("Broadcasting transaction onto LikeChain")
	result, err := tendermint.BroadcastTxCommit(tmClient, rawTx)
	if err != nil {
		log.
			WithField("raw_tx", common.Bytes2Hex(rawTx)).
			WithError(err).
			Panic("Broadcast transaction onto LikeChain failed")
	}
	if result.Code != response.Success.Code {
		switch result.Code {
		case response.DepositDoubleApproval.ToResponseCheckTx().Code:
			fallthrough
		case response.DepositDoubleApproval.ToResponseDeliverTx().Code:
			fallthrough
		case response.DepositAlreadyExecuted.ToResponseCheckTx().Code:
			fallthrough
		case response.DepositAlreadyExecuted.ToResponseDeliverTx().Code:
			log.
				WithField("code", result.Code).
				WithField("info", result.Info).
				WithField("log", result.Log).
				Info("Deposit transaction unnecessary and rejected, skipping")
		default:
			log.
				WithField("code", result.Code).
				WithField("info", result.Info).
				WithField("log", result.Log).
				Panic("Deposit transaction executed but failed")
		}
	} else {
		log.Info("Successfully executed deposit transaction onto LikeChain")
	}
}

type runState struct {
	LastEthBlock                   int64     `json:",omitempty"` // For compatibility
	LastProcessedSubscriptionBlock int64     // The last block processed by the subscriber goroutine
	PendingBlockRanges             [][]int64 // Ranges of blocks to be processed by the backward-scanner goroutine

	lock *sync.Mutex
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
	networkConfirmedBlock := eth.GetHeight(config.LoadBalancer) - config.BlockDelay
	if err != nil {
		log.
			WithField("state_path", config.StatePath).
			WithError(err).
			Info("Failed to load state, creating empty state")
		state = &runState{lock: &sync.Mutex{}}
		state.LastProcessedSubscriptionBlock = networkConfirmedBlock
		if config.StartFromBlock > 0 && networkConfirmedBlock >= config.StartFromBlock {
			state.PendingBlockRanges = [][]int64{{config.StartFromBlock, networkConfirmedBlock}}
		}
	}
	if state.LastEthBlock != 0 && state.LastProcessedSubscriptionBlock == 0 {
		state.LastProcessedSubscriptionBlock = state.LastEthBlock - config.BlockDelay
		log.
			WithField("last_eth_block", state.LastEthBlock).
			WithField("block_delay", config.BlockDelay).
			WithField("converted_last_processed_subscription_block", state.LastProcessedSubscriptionBlock).
			Info("Converted deprecated LastEthBlock to LastProcessedSubscriptionBlock in deposit state")
		state.LastEthBlock = 0
	}
	if state.LastProcessedSubscriptionBlock < networkConfirmedBlock {
		state.PendingBlockRanges = append(
			state.PendingBlockRanges,
			[]int64{state.LastProcessedSubscriptionBlock + 1, networkConfirmedBlock},
		)
		log.
			WithField("last_processed_subscription_block", state.LastProcessedSubscriptionBlock).
			WithField("network_confirmed_block", networkConfirmedBlock).
			WithField("new_pending_block_ranges", state.PendingBlockRanges).
			Info("Detected unprocessed blocks during service suspension, added to pending block ranges")
	}
	// The blocks before networkConfirmedBlock are processed by the backward-scanner goroutine
	state.LastProcessedSubscriptionBlock = networkConfirmedBlock
	state.save(config.StatePath)

	// The backward-scanner, which scans and processes pending blocks previously left
	go func() {
		defer httpHookCleanupFunc()
		log.
			WithField("ranges", state.PendingBlockRanges).
			Info("Clearing pending block ranges previously left and accumulated during service suspension")
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

	// The subscriber, which subscribes and processes newly confirmed Ethereum blocks
	eth.SubscribeHeader(config.LoadBalancer, func(header *ethTypes.Header) bool {
		newConfirmedBlock := header.Number.Int64() - config.BlockDelay
		if newConfirmedBlock <= 0 {
			return true
		}
		log.
			WithField("last_processed_block", state.LastProcessedSubscriptionBlock).
			WithField("new_confirmed_block", newConfirmedBlock).
			Info("Received new Ethereum block")
		for from := state.LastProcessedSubscriptionBlock + 1; from <= newConfirmedBlock; from += blockBatchSize {
			to := from + blockBatchSize - 1
			if to > newConfirmedBlock {
				to = newConfirmedBlock
			}
			scanAndProposeForRange(config, uint64(from), uint64(to))
			func() {
				state.lock.Lock()
				defer state.lock.Unlock()
				state.LastProcessedSubscriptionBlock = to
				state.save(config.StatePath)
			}()
		}
		return true
	})
}
