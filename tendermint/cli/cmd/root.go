package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{}

func init() {
	rootCmd.PersistentFlags().String("profile_dir", "./tendermint/nodes", "output directory for profiles containing generated config and key files")
	rootCmd.PersistentFlags().Uint("node_count", 1, "number of nodes")
	rootCmd.AddCommand(initCmd)
}

// Execute executes the root command, which defines the subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
