package eth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/types"

	"github.com/likecoin/likechain/services/abi/token"
	logger "github.com/likecoin/likechain/services/log"
)

var log = logger.L

// GetHeight returns the block number of the newest block header
func GetHeight(ethClient *ethclient.Client) int64 {
	header, _ := ethClient.HeaderByNumber(context.Background(), nil)
	return header.Number.Int64()
}

// SubscribeHeader subscribes to new block headers
func SubscribeHeader(ethClient *ethclient.Client, fn func(*ethTypes.Header) bool) {
	for ; ; time.Sleep(time.Minute) {
		log.Debug("Getting new Ethereum block")
		header, _ := ethClient.HeaderByNumber(context.Background(), nil)
		if !fn(header) {
			return
		}
	}
}

// SubscribeTransfer subscribes to Transfer events on the token contract to the relay contract address
func SubscribeTransfer(ethClient *ethclient.Client, tokenAddr, relayAddr common.Address, fn func(*deposit.Input, ethTypes.Log) bool) {
	tokenFilter, err := token.NewTokenFilterer(tokenAddr, ethClient)
	if err != nil {
		panic(err)
	}
	ch := make(chan *token.TokenTransfer)
	sub, err := tokenFilter.WatchTransfer(nil, ch, nil, []common.Address{relayAddr})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case e := <-ch:
			addr, err := types.NewAddress(e.From[:])
			if err != nil {
				panic(err)
			}
			input := deposit.Input{
				FromAddr: *addr,
				Value:    types.BigInt{Int: e.Value},
			}
			cont := fn(&input, e.Raw)
			if !cont {
				return
			}
		}
	}
}

// GetTransfersFromBlocks returns all the Transfer events on a token contract to a specific address
func GetTransfersFromBlocks(ethClient *ethclient.Client, tokenAddr, relayAddr common.Address, fromBlock, toBlock uint64) []deposit.Proposal {
	tokenFilter, err := token.NewTokenFilterer(tokenAddr, ethClient)
	if err != nil {
		panic(err)
	}
	opts := bind.FilterOpts{
		Start: fromBlock,
		End:   &toBlock,
	}
	it, err := tokenFilter.FilterTransfer(&opts, nil, []common.Address{relayAddr})
	if err != nil {
		panic(err)
	}
	proposals := []deposit.Proposal{}
	blockToProposal := map[uint64]*deposit.Proposal{}
	for it.Next() {
		e := it.Event
		blockNumber := e.Raw.BlockNumber
		proposal, ok := blockToProposal[blockNumber]
		if !ok {
			proposals = append(proposals, deposit.Proposal{})
			proposal = &proposals[len(proposals)-1]
			blockToProposal[blockNumber] = proposal
			proposal.BlockNumber = blockNumber
		}
		addr, err := types.NewAddress(e.From[:])
		if err != nil {
			panic(err)
		}
		proposal.Inputs = append(proposal.Inputs, deposit.Input{
			FromAddr: *addr,
			Value:    types.BigInt{Int: e.Value},
		})
	}
	return proposals
}
