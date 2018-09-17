package main

import (
	"github.com/likecoin/likechain/abci/app"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/tendermint/tendermint/abci/server"

	cmn "github.com/tendermint/tendermint/libs/common"
)

var log = logger.L

func main() {
	app := app.NewLikeChainApplication("/tmp")
	svr, err := server.NewServer("tcp://0.0.0.0:26658", "socket", app)
	if err != nil {
		log.WithError(err).Panic("Error when initializing server")
	}
	err = svr.Start()
	if err != nil {
		log.WithError(err).Panic("Error when starting server")
	}
	cmn.TrapSignal(func() {
		app.Stop()
		svr.Stop()
	})
}
