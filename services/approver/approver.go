package approver

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/likecoin/likechain/abci/query"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"

	"github.com/likecoin/likechain/services/eth"
	"github.com/likecoin/likechain/services/tendermint"
)

type depositTx struct {
	TxHash []byte
	Tx     txs.DepositTransaction
}

func fillSig(tx *txs.DepositApprovalTransaction, privKey *ecdsa.PrivateKey) {
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
	jsonSig := txs.DepositApprovalJSONSignature{}
	copy(jsonSig.JSONSignature[:], sig)
	tx.Sig = &jsonSig
}

func approve(tmClient *tmRPC.HTTP, tmPrivKey *ecdsa.PrivateKey, tx depositTx) {
	fmt.Printf("Approving txHash %s\n", common.Bytes2Hex(tx.TxHash))
	ethAddr := crypto.PubkeyToAddress(tmPrivKey.PublicKey)
	addr := types.NewAddress(ethAddr[:])
	if tx.Tx.Proposer.Equals(addr) {
		fmt.Printf("Deposit tx %v is by myself, skipping\n", tx.TxHash)
		return
	}

	queryResult, err := tmClient.ABCIQuery("tx_state", tx.TxHash)
	if err != nil {
		panic(err)
	}
	txState := query.GetTxStateRes(queryResult.Response.Value)
	if txState == nil {
		panic("Cannot parse tx_state result")
	}
	if txState.Status != "pending" {
		fmt.Printf("Deposit tx %v is not pending, skipping\n", tx.TxHash)
		return
	}

	queryResult, err = tmClient.ABCIQuery("account_info", []byte(addr.String()))
	if err != nil {
		panic(err)
	}
	accInfo := query.GetAccountInfoRes(queryResult.Response.Value)
	if accInfo == nil {
		panic("Cannot parse account_info result")
	}
	fmt.Printf("Nonce: %d\n", accInfo.NextNonce)
	approvalTx := &txs.DepositApprovalTransaction{
		Approver:      addr,
		DepositTxHash: tx.TxHash,
		Nonce:         accInfo.NextNonce,
	}
	fillSig(approvalTx, tmPrivKey)
	rawTx := txs.EncodeTx(approvalTx)
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
	for ; ; time.Sleep(time.Minute) {
		newHeight := uint64(tendermint.GetHeight(tmClient))
		if newHeight <= blockDelay {
			continue
		}
		fmt.Println(newHeight)
		fmt.Printf("Search deposits with %d < height <= %d\n", lastHeight, newHeight)
		queryString := fmt.Sprintf("deposit.height>%d AND deposit.height<=%d", lastHeight, newHeight)
		// TODO: may need pagination
		searchResult, err := tmClient.TxSearch(queryString, true, 1, 100)
		if err != nil {
			panic(err)
		}
		if searchResult.TotalCount <= 0 {
			fmt.Println("No search result")
			return
		}
		currentBlockNumber := eth.GetHeight(ethClient)
		for i := 0; i < searchResult.TotalCount; i++ {
			txHash := searchResult.Txs[i].Hash
			var tx txs.DepositTransaction
			err := types.AminoCodec().UnmarshalBinary(searchResult.Txs[i].Tx, &tx)
			if err != nil {
				fmt.Printf("Cannot unmarshal deposit tx %v\n", txHash)
				continue
			}
			if currentBlockNumber < int64(tx.Proposal.BlockNumber+blockDelay) {
				// TODO: may store it, check back later within subscription logic for Ethereum new block headers
				continue
			}
			events := eth.GetTransfersFromBlock(ethClient, tokenAddr, relayAddr, tx.Proposal.BlockNumber)
			if len(events) == 0 {
				continue
			}
			ok := true
			for i, e := range tx.Proposal.Inputs {
				// TODO: sort inputs
				if !e.FromAddr.Equals(types.NewAddress(events[i].From[:])) || e.Value.Cmp(events[i].Value) != 0 {
					ok = false
					break
				}
			}
			if ok {
				approve(tmClient, tmPrivKey, depositTx{txHash, tx})
			}
		}
		lastHeight = newHeight
		// TODO: store into DB
	}
}
