package app

import (
	"encoding/json"
	"fmt"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/query"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

var log = logger.L

// LikeChainApplication implements Tendermint ABCI interface
type LikeChainApplication struct {
	abci.BaseApplication

	ctx *context.ApplicationContext
}

// NewLikeChainApplication creates a new LikeChainApplication
func NewLikeChainApplication(ctx *context.ApplicationContext) *LikeChainApplication {
	return &LikeChainApplication{ctx: ctx}
}

// BeginBlock implements ABCI BeginBlock
func (app *LikeChainApplication) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	log.Info("APP BeginBlock")
	app.ctx.GetMutableState().SetBlockHash(req.Hash)
	return abci.ResponseBeginBlock{}
}

// CheckTx implements ABCI CheckTx
func (app *LikeChainApplication) CheckTx(rawTx []byte) abci.ResponseCheckTx {
	log.Info("APP CheckTx")
	var tx txs.Transaction
	err := types.AminoCodec().UnmarshalBinary(rawTx, &tx)
	if err != nil {
		log.WithError(err).Debug("APP CheckTx cannot parse transaction")
		return abci.ResponseCheckTx{
			Code: 1,
			Info: "Cannot parse transaction",
		}
	}
	return tx.CheckTx(app.ctx.GetImmutableState()).ToResponseCheckTx()
}

// DeliverTx implements ABCI DeliverTx
func (app *LikeChainApplication) DeliverTx(rawTx []byte) abci.ResponseDeliverTx {
	txHash := utils.HashRawTx(rawTx)
	log.
		WithField("hash", cmn.HexBytes(txHash)).
		Info("APP DeliverTx")

	var tx txs.Transaction
	err := types.AminoCodec().UnmarshalBinary(rawTx, &tx)
	if err != nil {
		log.WithError(err).Debug("APP DeliverTx cannot parse transaction")
		return abci.ResponseDeliverTx{
			Code: 1,
			Info: "Cannot parse transaction",
		}
	}
	state := app.ctx.GetMutableState()
	r := tx.DeliverTx(state, txHash)
	oldStatus := txstatus.GetStatus(state, txHash)
	if oldStatus != txstatus.TxStatusSuccess {
		txstatus.SetStatus(app.ctx.GetMutableState(), txHash, r.Status)
	}
	return r.ToResponseDeliverTx()
}

// EndBlock implements ABCI EndBlock
func (app *LikeChainApplication) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	log.WithField("height", req.Height).Info("APP EndBlock")
	return abci.ResponseEndBlock{} // TODO
}

// Commit implements ABCI Commit
func (app *LikeChainApplication) Commit() abci.ResponseCommit {
	state := app.ctx.GetMutableState()
	height := state.GetHeight() + 1
	state.SetHeight(height)
	hash := state.Save()

	stateTreeVersion := state.MutableStateTree().Version64()
	withdrawTreeVersion := state.MutableWithdrawTree().Version64()
	state.SetMetadataAtHeight(height, context.TreeMetadata{
		StateTreeVersion:    stateTreeVersion,
		WithdrawTreeVersion: withdrawTreeVersion,
	})
	state.GC(height)

	log.
		WithField("hash", cmn.HexBytes(hash)).
		WithField("height", height).
		WithField("state_tree_version", stateTreeVersion).
		WithField("withdraw_tree_version", withdrawTreeVersion).
		Info("APP Commit")

	return abci.ResponseCommit{Data: hash}
}

// InitChain implements ABCI InitChain
func (app *LikeChainApplication) InitChain(params abci.RequestInitChain) abci.ResponseInitChain {
	log.
		Info("APP InitChain")
	app.ctx.GetMutableState().Init()
	if len(params.AppStateBytes) > 0 {
		log.
			WithField("app_state_bytes", string(params.AppStateBytes)).
			Info("Initializing app state")
		appInitState := types.AppInitState{}
		err := json.Unmarshal(params.AppStateBytes, &appInitState)
		if err != nil {
			log.
				WithError(err).
				Panic("Cannot load app initial state")
		}
		state := app.ctx.GetMutableState()
		usedIDs := make(map[types.LikeChainID]bool)
		usedAddrs := make(map[types.Address]bool)
		for i, accInfo := range appInitState.Accounts {
			id := accInfo.ID
			if usedIDs[id] {
				log.
					WithField("entry_number", i).
					WithField("id", id).
					Panic("Duplicated LikeChainID")
			}
			usedIDs[id] = true
			addr := accInfo.Addr
			if usedAddrs[addr] {
				log.
					WithField("entry_number", i).
					WithField("addr", addr).
					Panic("Duplicated address")
			}
			usedAddrs[addr] = true
			balance := accInfo.Balance
			if !balance.IsWithinRange() {
				log.
					WithField("entry_number", i).
					WithField("balance", balance).
					Panic("Invalid initial balance")
			}
			account.NewAccountFromID(state, &id, &accInfo.Addr)
			account.SaveBalance(state, &id, balance.Int)
		}
	}
	return abci.ResponseInitChain{
		ConsensusParams: params.ConsensusParams,
		Validators:      params.Validators,
	}
}

// Query implements ABCI Query
func (app *LikeChainApplication) Query(reqQuery abci.RequestQuery) abci.ResponseQuery {
	log.WithField("path", reqQuery.Path).Info("APP Query")
	return query.Query(app.ctx.GetMutableState(), reqQuery)
}

// Info implements ABCI Info
func (app *LikeChainApplication) Info(req abci.RequestInfo) abci.ResponseInfo {
	state := app.ctx.GetImmutableState()
	height := state.GetHeight()
	appHash := cmn.HexBytes(state.GetAppHash())

	log.
		WithField("height", height).
		WithField("hash", appHash).
		Info("APP Info")

	return abci.ResponseInfo{
		Data:             fmt.Sprintf("{\"hash\":\"%v\"}", appHash),
		Version:          fmt.Sprintf("LikeChain v1.0 (Tendermint v%s", req.Version),
		LastBlockHeight:  height,
		LastBlockAppHash: appHash,
	}
}

// Stop handles quit app
func (app *LikeChainApplication) Stop() {
	app.ctx.CloseDb()
}
