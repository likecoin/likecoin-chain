package deposit

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/likecoin/likechain/services/eth"
	logger "github.com/likecoin/likechain/services/log"
)

const blockBatchSize = 10000

var log = logger.L

type runState struct {
	LastEthBlock                   int64     `json:",omitempty"` // For compatibility
	LastProcessedSubscriptionBlock int64     // The last block processed by the subscriber goroutine
	PendingBlockRanges             [][]int64 // Ranges of blocks to be processed by the backward-scanner goroutine
	AbandonedBlocks                *queue    // Blocks which are timeout and need retry in lower priority

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

func abandonedBlocksPoller(config *Config, state *runState) {
	for {
		time.Sleep(10 * time.Second)
		var block uint64
		var ok bool
		func() {
			state.lock.Lock()
			defer state.lock.Unlock()
			block, ok = state.AbandonedBlocks.peek()
		}()
		if !ok {
			continue
		}
		log.
			WithField("block", block).
			Info("Processing previously abandoned block")
		abandonedBlocks := scanAndProposeForRange(config, block, block)
		func() {
			state.lock.Lock()
			defer state.lock.Unlock()
			state.AbandonedBlocks.dequeue()
			for _, block := range abandonedBlocks {
				state.AbandonedBlocks.enqueue(block)
			}
			state.save(config.StatePath)
		}()
	}
}

func backwardScanner(config *Config, state *runState) {
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
			abandonedBlocks := scanAndProposeForRange(config, uint64(from), uint64(to))
			func() {
				state.lock.Lock()
				defer state.lock.Unlock()
				if from == start {
					state.PendingBlockRanges = state.PendingBlockRanges[:i]
				} else {
					state.PendingBlockRanges[i][1] = from - 1
				}
				for _, block := range abandonedBlocks {
					state.AbandonedBlocks.enqueue(block)
				}
				state.save(config.StatePath)
			}()
		}
	}
	log.Info("All pending blocks cleared")
}

func ethSubscriber(config *Config, state *runState) {
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
			abandonedBlocks := scanAndProposeForRange(config, uint64(from), uint64(to))
			func() {
				state.lock.Lock()
				defer state.lock.Unlock()
				state.LastProcessedSubscriptionBlock = to
				for _, block := range abandonedBlocks {
					state.AbandonedBlocks.enqueue(block)
				}
				state.save(config.StatePath)
			}()
		}
		return true
	})
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
		state = &runState{
			lock: &sync.Mutex{},
		}
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
	if state.AbandonedBlocks == nil {
		state.AbandonedBlocks = newQueue(32)
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

	// abandoned blocks poller, which polls blocks that are timeout in lower priority
	go func() {
		defer httpHookCleanupFunc()
		abandonedBlocksPoller(config, state)
	}()

	// The backward-scanner, which scans and processes pending blocks previously left
	go func() {
		defer httpHookCleanupFunc()
		backwardScanner(config, state)
	}()

	// The subscriber, which subscribes and processes newly confirmed Ethereum blocks
	ethSubscriber(config, state)
}
