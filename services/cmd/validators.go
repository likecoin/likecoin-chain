package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

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
		validatorFile := viper.GetString("validatorFile")
		if validatorFile == "" {
			tmClient := tmRPC.NewHTTP(viper.GetString("tmEndPoint"), "/websocket")
			validators := tendermint.GetValidators(tmClient)
			for i, v := range validators {
				pubKey := v.PubKey.(secp256k1.PubKeySecp256k1)
				ethAddr := tendermint.PubKeyToEthAddr(&pubKey)
				fmt.Printf("Validator %d: %v\n", i, ethAddr.Hex())
			}
		} else {
			data, err := ioutil.ReadFile(validatorFile)
			if err != nil {
				panic(err)
			}
			pubKeysBase64 := []string{}
			err = json.Unmarshal(data, &pubKeysBase64)
			if err != nil {
				panic(err)
			}
			for i, pubKeyBase64 := range pubKeysBase64 {
				pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyBase64)
				if err != nil {
					panic(err)
				}
				pubKey := secp256k1.PubKeySecp256k1{}
				copy(pubKey[:], pubKeyBytes)
				ethAddr := tendermint.PubKeyToEthAddr(&pubKey)
				fmt.Printf("Validator %d: %v\n", i, ethAddr.Hex())
			}
		}
	},
}

func init() {
	validatorsCmd.PersistentFlags().String("from-file", "", "Load validators from file instead of tendermint endpoint")
	viper.BindPFlag("validatorFile", validatorsCmd.PersistentFlags().Lookup("from-file"))
}
