package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/api/routes"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
)

func main() {
	// TODO: Put in config file
	host := os.Getenv("LIKECHAIN_API_CLIENT_HOST")
	if host == "" {
		host = "localhost:26657"
	}

	port := os.Getenv("LIKECHAIN_API_PORT")
	if port == "" {
		port = "3000"
	}

	client := rpcClient.NewHTTP("tcp://"+host, "/websocket")

	router := gin.Default()

	routes.Initialize(router, client)

	router.Run(":" + port)
}
