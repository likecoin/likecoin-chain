package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/cmd/api/routes"
	customvalidator "github.com/likecoin/likechain/abci/cmd/api/validator"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
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

	client := rpcclient.NewHTTP("tcp://"+host, "/websocket")

	router := gin.Default()

	customvalidator.Bind()

	// Allows all origins
	router.Use(cors.Default())

	routes.Initialize(router, client)

	router.Run(":" + port)
}
