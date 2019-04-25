package eth

import (
	"context"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
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
func GetHeight(lb *LoadBalancer) int64 {
	height := int64(0)
	lb.Do(func(ethClient *ethclient.Client) error {
		c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		header, err := ethClient.HeaderByNumber(c, nil)
		if err != nil {
			log.WithError(err).Error("Cannot get latest Ethereum header")
			return err
		}
		height = header.Number.Int64()
		return nil
	})
	return height
}

// SubscribeHeader subscribes to new block headers
func SubscribeHeader(lb *LoadBalancer, fn func(*ethTypes.Header) bool) {
	for ; ; time.Sleep(time.Minute) {
		log.Debug("Getting new Ethereum block")
		var header *ethTypes.Header
		lb.Do(func(ethClient *ethclient.Client) error {
			c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			h, err := ethClient.HeaderByNumber(c, nil)
			if err != nil {
				log.WithError(err).Panic("Cannot get latest Ethereum header")
				return err
			}
			header = h
			return nil
		})
		if !fn(header) {
			return
		}
	}
}

// GetTransfersFromBlocks returns all the Transfer events on a token contract to a specific address
func GetTransfersFromBlocks(lb *LoadBalancer, tokenAddr, relayAddr common.Address, fromBlock, toBlock uint64) []deposit.Proposal {
	var proposals []deposit.Proposal
	from := fromBlock
	to := toBlock
	for from <= to {
		lb.Do(func(ethClient *ethclient.Client) error {
			tokenFilter, err := token.NewTokenFilterer(tokenAddr, ethClient)
			if err != nil {
				return err
			}
			c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			opts := bind.FilterOpts{
				Start:   from,
				End:     &to,
				Context: c,
			}
			it, err := tokenFilter.FilterTransfer(&opts, nil, []common.Address{relayAddr})
			if err != nil {
				if strings.Contains(err.Error(), "query returned more than") && from != to { // Infura refuses to return more than 1000 results
					log.
						WithField("begin_block", from).
						WithField("end_block", to).
						WithField("new_end_block", from+(to-from)/2).
						Info("Endpoint complained too many query results, backoff the range by half")
					to = from + (to-from)/2
					return nil
				}
				return err
			}
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
			from = to + 1
			to = toBlock
			return nil
		})
	}
	return proposals
}

// WaitForReceipt polls for the receipt of a transaction
func WaitForReceipt(lb *LoadBalancer, txHash common.Hash) *ethTypes.Receipt {
	var receipt *ethTypes.Receipt
	for {
		lb.Do(func(ethClient *ethclient.Client) error {
			c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			r, err := ethClient.TransactionReceipt(c, txHash)
			if r != nil {
				receipt = r
				return nil
			}
			if err != ethereum.NotFound {
				return err
			}
			return nil
		})
		if receipt != nil {
			return receipt
		}
		time.Sleep(15 * time.Second)
	}
}

// GetNonce return the current usable nonce for the given Ethereum address
func GetNonce(lb *LoadBalancer, addr common.Address) int64 {
	nonce := int64(0)
	lb.Do(func(ethClient *ethclient.Client) error {
		c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ethNonce, err := ethClient.NonceAt(c, addr, nil)
		if err != nil {
			log.
				WithField("addr", addr.Hex()).
				WithError(err).
				Error("Error when getting Ethereum nonce")
			return err
		}
		nonce = int64(ethNonce)
		return nil
	})
	return nonce
}
