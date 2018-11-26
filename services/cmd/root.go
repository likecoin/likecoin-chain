package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "likechain",
	Short: "likechain is a program for LikeChain validators to run LikeChain related background services.",
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./config.json)")
	rootCmd.AddCommand(withdrawCmd)
	rootCmd.AddCommand(validatorsCmd)
	rootCmd.AddCommand(proposerCmd)
	rootCmd.AddCommand(approverCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("./config.json")
	}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// Execute executes the root command, which defines the subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
