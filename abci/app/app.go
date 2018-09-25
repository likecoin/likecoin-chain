package app

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/query"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

var log = logger.L

type LikeChainApplication struct {
	abci.BaseApplication

	ctx *context.ApplicationContext
}

func NewLikeChainApplication(dbPath string) *LikeChainApplication {
	return &LikeChainApplication{
		ctx: context.New(dbPath),
	}
}

func (app *LikeChainApplication) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	log.Info("APP BeginBlock")
	app.ctx.GetMutableState().SetBlockHash(req.Hash)
	return abci.ResponseBeginBlock{}
}

func (app *LikeChainApplication) CheckTx(rawTx []byte) abci.ResponseCheckTx {
	log.Info("APP CheckTx")
	tx := &types.Transaction{}
	if err := proto.Unmarshal(rawTx, tx); err != nil {
		log.WithError(err).Debug("APP CheckTx cannot parse transaction")
		return abci.ResponseCheckTx{
			Code: 1,
			Info: "Cannot parse transaction",
		}
	}
	return handlers.CheckTx(app.ctx.GetImmutableState(), tx)
}

func (app *LikeChainApplication) DeliverTx(rawTx []byte) abci.ResponseDeliverTx {
	log.Info("APP DeliverTx")
	tx := &types.Transaction{}
	if err := proto.Unmarshal(rawTx, tx); err != nil {
		log.WithError(err).Debug("APP DeliverTx Cannot parse transaction")
		return abci.ResponseDeliverTx{
			Code: 1,
			Info: "Cannot parse transaction",
		}
	}
	return handlers.DeliverTx(app.ctx.GetMutableState(), tx, utils.HashRawTx(rawTx))
}

func (app *LikeChainApplication) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	log.WithField("height", req.Height).Info("APP EndBlock")
	return abci.ResponseEndBlock{} // TODO
}

func (app *LikeChainApplication) Commit() abci.ResponseCommit {
	state := app.ctx.GetMutableState()
	height := state.GetHeight() + 1
	state.SetHeight(height)
	hash := state.Save()

	withdrawTreeVersion := state.MutableWithdrawTree().Version64()
	state.SetWithdrawVersionAtHeight(height, withdrawTreeVersion)

	log.
		WithField("hash", cmn.HexBytes(hash)).
		WithField("height", height).
		WithField("withdraw_tree_version", withdrawTreeVersion).
		Info("APP Commit")

	return abci.ResponseCommit{Data: hash}
}

func (app *LikeChainApplication) InitChain(params abci.RequestInitChain) abci.ResponseInitChain {
	app.ctx.GetMutableState().Init()

	log.
		Info("APP InitChain")

	return abci.ResponseInitChain{
		ConsensusParams: params.ConsensusParams,
		Validators:      params.Validators,
	}
}

func (app *LikeChainApplication) Query(reqQuery abci.RequestQuery) abci.ResponseQuery {
	log.Info("APP Query")
	return query.Query(app.ctx.GetMutableState(), reqQuery)
}

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
