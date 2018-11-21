package proposer

import (
	"crypto/ecdsa"
	"fmt"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/likecoin/likechain/abci/query"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"

	"github.com/likecoin/likechain/services/abi/token"
	"github.com/likecoin/likechain/services/eth"
)

func fillSig(tx *txs.DepositTransaction, privKey *ecdsa.PrivateKey) {
	tx.Proposal.Sort()
	jsonMap := tx.GenerateJSONMap()
	hash, err := txs.JSONMapToHash(jsonMap)
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

func propose(tmClient *tmRPC.HTTP, tmPrivKey *ecdsa.PrivateKey, blockNumber uint64, events []token.TokenTransfer) {
	if len(events) == 0 {
		return
	}
	fmt.Printf("Proposing for blockNumber %d\n", blockNumber)
	ethAddr := crypto.PubkeyToAddress(tmPrivKey.PublicKey)
	addr := types.NewAddress(ethAddr[:])
	queryResult, err := tmClient.ABCIQuery("account_info", []byte(addr.String()))
	if err != nil {
		panic(err)
	}
	accInfo := query.GetAccountInfoRes(queryResult.Response.Value)
	if accInfo == nil {
		panic("Cannot parse account_info result")
	}
	fmt.Printf("Nonce: %d\n", accInfo.NextNonce)
	inputs := make([]deposit.Input, 0, len(events))
	for _, e := range events {
		inputs = append(inputs, deposit.Input{
			FromAddr: *types.NewAddress(e.From[:]),
			Value:    types.BigInt{Int: e.Value},
		})
	}
	tx := &txs.DepositTransaction{
		Proposer: addr,
		Proposal: deposit.Proposal{
			BlockNumber: blockNumber,
			Inputs:      inputs,
		},
		Nonce: accInfo.NextNonce,
	}
	fillSig(tx, tmPrivKey)
	rawTx := txs.EncodeTx(tx)
	fmt.Printf("Now broadcasting, rawTx: %v\n", rawTx)
	_, err = tmClient.BroadcastTxCommit(rawTx)
	fmt.Printf("After broadcast\n")
	if err != nil {
		panic(err)
	}
}

// Run starts the subscription to the deposits on Ethereum into the relay contract and commits proposal onto LikeChain
func Run(tmClient *tmRPC.HTTP, ethClient *ethclient.Client, tokenAddr, relayAddr common.Address, tmPrivKey *ecdsa.PrivateKey, blockDelay uint64) {
	lastHeight := uint64(0) // TODO: load from DB
	eth.SubscribeHeader(ethClient, func(header *ethTypes.Header) bool {
		blockNumber := header.Number.Int64()
		if blockNumber <= 0 {
			return true
		}
		newHeight := uint64(blockNumber)
		fmt.Println(newHeight)
		if newHeight < blockDelay {
			return true
		}
		for h := lastHeight; h <= newHeight-blockDelay; h++ {
			events := eth.GetTransfersFromBlock(ethClient, tokenAddr, relayAddr, h)
			if len(events) == 0 {
				continue
			}
			propose(tmClient, tmPrivKey, h, events)
		}
		lastHeight = newHeight
		return true
	})
}
