package app

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

type LikeChainApplication struct {
	abci.BaseApplication
	// TODO
}

func (app *LikeChainApplication) CheckTx(rawTx []byte) abci.ResponseCheckTx {
	return abci.ResponseCheckTx{} // TODO
}

func (app *LikeChainApplication) DeliverTx(rawTx []byte) abci.ResponseDeliverTx {
	return abci.ResponseDeliverTx{} // TODO
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
