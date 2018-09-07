package app

import (
	"github.com/gogo/protobuf/proto"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	abci "github.com/tendermint/tendermint/abci/types"
)

type LikeChainApplication struct {
	abci.BaseApplication

	ctx *context.ApplicationContext
}

func (app *LikeChainApplication) CheckTx(rawTx []byte) abci.ResponseCheckTx {
	tx := &types.Transaction{}
	if err := proto.Unmarshal(rawTx, tx); err != nil {
		return abci.ResponseCheckTx{
			Code: 1,
			Info: "Cannot parse transaction",
		}
	}
	return handlers.CheckTx(app.ctx.GetImmutableState(), tx)
}

func (app *LikeChainApplication) DeliverTx(rawTx []byte) abci.ResponseDeliverTx {
	tx := &types.Transaction{}
	if err := proto.Unmarshal(rawTx, tx); err != nil {
		return abci.ResponseDeliverTx{
			Code: 1,
			Info: "Cannot parse transaction",
		}
	}
	return handlers.DeliverTx(app.ctx.GetMutableState(), tx, utils.HashRawTx(rawTx))
}

func (app *LikeChainApplication) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{} // TODO
}

func (app *LikeChainApplication) Commit() abci.ResponseCommit {
	return abci.ResponseCommit{} // TODO
}

func (app *LikeChainApplication) InitChain(params abci.RequestInitChain) abci.ResponseInitChain {
	return abci.ResponseInitChain{} // TODO
}

func (app *LikeChainApplication) Query(reqQuery abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{} // TODO
}

func (app *LikeChainApplication) Info(req abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{} // TODO
}
