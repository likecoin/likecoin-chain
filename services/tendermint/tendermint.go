package tendermint

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmRPC "github.com/tendermint/tendermint/rpc/client"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"

	logger "github.com/likecoin/likechain/services/log"
)

var log = logger.L

// BroadcastTxResult wraps the result of BroadcastTxCommit in different cases (e.g. Tx already committed, Tx timeout)
type BroadcastTxResult struct {
	Code   uint32
	Info   string
	Log    string
	Hash   cmn.HexBytes
	Height int64
}

// BroadcastTimeout wraps Tendermint broadcast transaction error with a type of timeout semantic
type BroadcastTimeout struct {
	error
}

// BroadcastTxCommit broadcast a Tendermint transaction if needed, then query and return the results
func BroadcastTxCommit(tmClient *tmRPC.HTTP, rawTx types.Tx) (*BroadcastTxResult, error) {
	txHash := tmhash.Sum(rawTx)
	txResult, err := tmClient.Tx(txHash, false)
	if err == nil {
		log.
			WithField("tx_hash", common.Bytes2Hex(txHash)).
			WithField("tx_height", txResult.Height).
			Info("Deposit tx is already processed, skipping")
		return &BroadcastTxResult{
			Code:   txResult.TxResult.Code,
			Info:   txResult.TxResult.Info,
			Log:    txResult.TxResult.Log,
			Hash:   txResult.Hash,
			Height: txResult.Height,
		}, nil
	}
	broadcastTxResult, err := tmClient.BroadcastTxCommit(rawTx)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "Tx already exists in cache") || strings.Contains(errMsg, "Timed out waiting for tx to be included in a block") {
			return nil, BroadcastTimeout{error: err}
		}
		return nil, err
	}
	result := &BroadcastTxResult{
		Hash:   broadcastTxResult.Hash,
		Height: broadcastTxResult.Height,
	}
	if broadcastTxResult.CheckTx.Code != 0 {
		result.Code = broadcastTxResult.CheckTx.Code
		result.Info = broadcastTxResult.CheckTx.Info
		result.Log = broadcastTxResult.CheckTx.Log
	} else {
		result.Code = broadcastTxResult.DeliverTx.Code
		result.Info = broadcastTxResult.DeliverTx.Info
		result.Log = broadcastTxResult.DeliverTx.Log
	}
	return result, nil
}

// GetSignedHeader returns the signed header at the given height
func GetSignedHeader(tmClient *tmRPC.HTTP, height int64) types.SignedHeader {
	commit, err := tmClient.Commit(&height)
	if err != nil {
		log.
			WithField("height", height).
			WithError(err).
			Panic("Cannot get Tendermint commit with height")
	}
	return commit.SignedHeader
}

// GetValidators returns the current validators
func GetValidators(tmClient *tmRPC.HTTP) []types.Validator {
	rawConsensusState, err := tmClient.DumpConsensusState()
	if err != nil {
		log.
			WithError(err).
			Panic("Cannot dump Tendermint consensus state")
	}

	jsonRes := struct {
		Validators struct {
			Validators []types.Validator `json:"validators"`
		} `json:"validators"`
	}{}

	err = AminoCodec().UnmarshalJSON(rawConsensusState.RoundState, &jsonRes)
	if err != nil {
		log.
			WithField("round_state", rawConsensusState.RoundState).
			WithError(err).
			Panic("Cannot unmarshal consensus round state into JSON")
	}

	return jsonRes.Validators.Validators
}

// GetHeight returns the current height
func GetHeight(tmClient *tmRPC.HTTP) int64 {
	abciInfo, err := tmClient.ABCIInfo()
	if err != nil {
		log.
			WithError(err).
			Panic("Cannot get Tendermint ABCI info")
	}
	return abciInfo.Response.GetLastBlockHeight()
}

// TxSearch returns all transactions with specific tag valued within range
func TxSearch(tmClient *tmRPC.HTTP, tag string, from, to int64) []*core_types.ResultTx {
	queryString := fmt.Sprintf("%s>=%d AND %s<=%d", tag, from, tag, to)
	var result []*core_types.ResultTx
	doneCount := 0
	for page := 1; ; page++ {
		searchResult, err := tmClient.TxSearch(queryString, true, page, 100)
		if err != nil {
			log.
				WithField("tag", tag).
				WithField("from", from).
				WithField("to", to).
				WithField("page", page).
				WithError(err).
				Panic("Search transaction failed")
		}
		if searchResult.TotalCount <= 0 {
			return nil
		}
		if result == nil {
			result = make([]*core_types.ResultTx, 0, searchResult.TotalCount)
		}
		for _, tx := range searchResult.Txs {
			result = append(result, tx)
		}
		doneCount += len(searchResult.Txs)
		if doneCount >= searchResult.TotalCount {
			return result
		}
	}
}
