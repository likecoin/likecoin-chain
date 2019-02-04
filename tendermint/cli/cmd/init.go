package cmd

import (
	"github.com/spf13/cobra"

	likeInit "github.com/likecoin/likechain/tendermint/cli/init"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config, key and data files for Tendermint",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		profileDir, err := cmd.Flags().GetString("profile_dir")
		if err != nil {
			panic(err)
		}
		dockerDir, err := cmd.Flags().GetString("docker_dir")
		if err != nil {
			panic(err)
		}
		nodeCount, err := cmd.Flags().GetUint("node_count")
		if err != nil {
			panic(err)
		}
		if nodeCount == 0 {
			panic("Node count must be greater than 0")
		}
		likeInit.Run(profileDir, dockerDir, nodeCount)
	},
}

func init() {
	initCmd.PersistentFlags().String("docker_dir", ".", "directory containing docker-compose.sample.yml and docker-compose.production.sample.yml")
}
