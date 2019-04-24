package deposit

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"

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

var log = logger.L

func fillSig(tx *txs.DepositTransaction, privKey *ecdsa.PrivateKey) {
	tx.Proposal.Sort()
	jsonMap := tx.GenerateJSONMap()
	hash, err := jsonMap.Hash()
	if err != nil {
		panic(err)
	}
	sig, err := crypto.Sign(hash, privKey)
	if err != nil {
		panic(err)
	}
	sig[64] += 27
	jsonSig := txs.DepositJSONSignature{}
	copy(jsonSig.JSONSignature[:], sig)
	tx.Sig = &jsonSig
}

func propose(tmClient *tmRPC.HTTP, tmPrivKey *ecdsa.PrivateKey, proposal deposit.Proposal) {
	log.
		WithField("block_number", proposal.BlockNumber).
		WithField("event_count", len(proposal.Inputs)).
		Info("Proposing new proposal")
	ethAddr := crypto.PubkeyToAddress(tmPrivKey.PublicKey)
	addr, err := types.NewAddress(ethAddr[:])
	if err != nil {
		panic(err)
	}
	queryResult, err := tmClient.ABCIQuery("account_info", []byte(addr.String()))
	if err != nil {
		panic(err)
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
		log.
			WithField("code", result.CheckTx.Code).
			WithField("info", result.CheckTx.Info).
			WithField("log", result.CheckTx.Log).
			Error("Deposit transaction failed in CheckTx")
	} else if result.DeliverTx.Code != response.Success.Code {
		log.
			WithField("code", result.DeliverTx.Code).
			WithField("info", result.DeliverTx.Info).
			WithField("log", result.DeliverTx.Log).
			Error("Deposit transaction failed in DeliverTx")
	} else {
		log.Info("Successfully broadcasted deposit transaction onto LikeChain")
	}
}

type proposedBlockSet struct {
	Map       map[uint64]bool
	Queue     []uint64
	QueueHead int
	QueueTail int
	Capacity  int
}

func newProposedBlockSet(capacity int) proposedBlockSet {
	return proposedBlockSet{
		Map:       make(map[uint64]bool),
		Queue:     make([]uint64, capacity),
		QueueHead: 0,
		QueueTail: 0,
		Capacity:  capacity,
	}
}

func (set proposedBlockSet) Has(block uint64) bool {
	return set.Map[block]
}

func (set proposedBlockSet) Put(block uint64) {
	set.Map[block] = true
	if len(set.Map) > set.Capacity {
		toRemove := set.Queue[set.QueueHead]
		delete(set.Map, toRemove)
		set.QueueHead = (set.QueueHead + 1) % set.Capacity
	}
	set.Queue[set.QueueTail] = block
	set.QueueTail = (set.QueueTail + 1) % set.Capacity
}

type runState struct {
	LastEthBlock int64
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
}

// Run starts the subscription to the deposits on Ethereum into the relay contract and commits proposal onto LikeChain
func Run(config *Config) {
	state, err := loadState(config.StatePath)
	if err != nil {
		log.
			WithField("state_path", config.StatePath).
			WithError(err).
			Info("Failed to load state, creating empty state")
		state = &runState{}
		blockNumber := eth.GetHeight(config.LoadBalancer)
		state.LastEthBlock = config.StartFromBlock + config.BlockDelay
		if config.StartFromBlock < 0 || blockNumber < state.LastEthBlock {
			state.LastEthBlock = blockNumber
		}
		state.save(config.StatePath)
	}
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
		proposals := eth.GetTransfersFromBlocks(
			config.LoadBalancer,
			config.TokenAddr,
			config.RelayAddr,
			uint64(state.LastEthBlock-config.BlockDelay),
			uint64(newBlock-config.BlockDelay-1),
		)
		if len(proposals) == 0 {
			log.
				WithField("begin_block", state.LastEthBlock-config.BlockDelay).
				WithField("end_block", newBlock-config.BlockDelay-1).
				Info("No transfer events in range")
		} else {
			for _, proposal := range proposals {
				utils.RetryIfPanic(5, func() {
					propose(config.TMClient, config.TMPrivKey, proposal)
				})
			}
		}
		state.LastEthBlock = newBlock
		state.save(config.StatePath)
		return true
	})
}
