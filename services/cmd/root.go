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
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "likechain",
	Short: "likechain is a program for LikeChain validators to run LikeChain related background services.",
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			log.Level = logrus.DebugLevel
			log.Debug("Using debug mode")
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.json", "config file (default: ./config.json)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logs")
	rootCmd.AddCommand(withdrawCmd)
	rootCmd.AddCommand(validatorsCmd)
	rootCmd.AddCommand(proposerCmd)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.
			WithError(err).
			Panic("Cannot read config file")
	}
}

// Execute executes the root command, which defines the subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Panic("Execution failed")
	}
}
