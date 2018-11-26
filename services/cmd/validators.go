package cmd

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/likecoin/likechain/services/tendermint"
)

var validatorsCmd = &cobra.Command{
	Use:   "validators",
	Short: "show current validators",
	Run: func(cmd *cobra.Command, args []string) {
		tmClient := tmRPC.NewHTTP(viper.GetString("tmEndPoint"), "/websocket")
		validators := tendermint.GetValidators(tmClient)
		for i, v := range validators {
			pubKey := v.PubKey.(secp256k1.PubKeySecp256k1)
			ethAddr := tendermint.PubKeyToEthAddr(&pubKey)
			fmt.Printf("Validator %d: %v\n", i, ethAddr.Hex())
		}
	},
}
