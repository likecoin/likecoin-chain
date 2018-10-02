package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/likecoin/likechain/abci/cmd/api/routes"
	customvalidator "github.com/likecoin/likechain/abci/cmd/api/validator"
	appConf "github.com/likecoin/likechain/abci/config"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var (
	config = appConf.GetConfig()
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

	if config.IsProduction() {
		// Disable Console Color, you don't need console color when writing the logs to file
		gin.DisableConsoleColor()

		// Logging to a file
		file, err := os.Create("api.log")
		if err != nil {
			log.Panic("Unable to create log file")
		}

		// Write the logs to file and console at the same time
		gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	}

	router := gin.Default()

	customvalidator.Bind()

	// Allows all origins
	router.Use(cors.Default())

	routes.Initialize(router, client)

	log.Println("Listen and serve on 0.0.0.0:" + port)
	router.Run(":" + port)
}
