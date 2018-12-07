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
)

// GetHeight returns the block number of the newest block header
func GetHeight(ethClient *ethclient.Client) int64 {
	header, _ := ethClient.HeaderByNumber(context.Background(), nil)
	return header.Number.Int64()
}

// SubscribeHeader subscribes to new block headers
func SubscribeHeader(ethClient *ethclient.Client, fn func(*ethTypes.Header) bool) {
	for ; ; time.Sleep(time.Minute) {
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

// GetTransfersFromBlock returns all the Transfer events on a token contract to a specific address
func GetTransfersFromBlock(ethClient *ethclient.Client, tokenAddr, relayAddr common.Address, blockNumber uint64) []token.TokenTransfer {
	tokenFilter, err := token.NewTokenFilterer(tokenAddr, ethClient)
	if err != nil {
		panic(err)
	}
	opts := bind.FilterOpts{
		Start: blockNumber,
		End:   &blockNumber,
	}
	it, err := tokenFilter.FilterTransfer(&opts, nil, []common.Address{relayAddr})
	if err != nil {
		panic(err)
	}
	events := []token.TokenTransfer{}
	for it.Next() {
		events = append(events, *it.Event)
	}
	return events
}
