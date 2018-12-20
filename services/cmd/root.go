package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	logger "github.com/likecoin/likechain/services/log"

	"github.com/sirupsen/logrus"
)

var (
	log     = logger.L
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "likechain",
	Short: "likechain is a program for LikeChain validators to run LikeChain related background services.",
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.json", "config file path")

	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logs")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().String("tm-endpoint", "tcp://localhost:26657", "Tendermint endpoint")
	viper.BindPFlag("tmEndPoint", rootCmd.PersistentFlags().Lookup("tm-endpoint"))

	rootCmd.PersistentFlags().String("eth-endpoint", "http://localhost:8545", "Ethereum endpoint")
	viper.BindPFlag("ethEndPoint", rootCmd.PersistentFlags().Lookup("eth-endpoint"))

	rootCmd.PersistentFlags().String("relay-addr", "", "Ethereum address of the relay contract")
	viper.BindPFlag("relayContractAddr", rootCmd.PersistentFlags().Lookup("relay-addr"))

	rootCmd.AddCommand(withdrawCmd)
	rootCmd.AddCommand(validatorsCmd)
	rootCmd.AddCommand(depositCmd)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.
			WithError(err).
			Panic("Cannot read config file")
	}
	if viper.GetBool("debug") {
		log.Level = logrus.DebugLevel
		log.Debug("Using debug mode")
	}
}

// Execute executes the root command, which defines the subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Panic("Execution failed")
	}
}
