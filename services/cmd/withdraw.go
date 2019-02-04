package cmd

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/likecoin/likechain/services/eth"
	"github.com/likecoin/likechain/services/withdraw"
)

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "run the withdraw-relay service",
	Run: func(cmd *cobra.Command, args []string) {
		tmEndPoint := viper.GetString("tmEndPoint")
		ethEndPoints := viper.GetStringSlice("ethEndPoints")
		relayAddr := common.HexToAddress(viper.GetString("relayContractAddr"))
		statePath := viper.GetString("withdrawStatePath")
		log.
			WithField("tm_endpoint", tmEndPoint).
			WithField("eth_endpoints", ethEndPoints).
			WithField("relay_addr", relayAddr).
			WithField("state_path", statePath).
			Debug("Read withdraw config and parameters")

		tmClient := tmRPC.NewHTTP(tmEndPoint, "/websocket")
		lb := eth.NewLoadBalancer(ethEndPoints)
		privKeyBytes := common.Hex2Bytes(viper.GetString("ethPrivKey"))
		privKey, err := ethCrypto.ToECDSA(privKeyBytes)
		if err != nil {
			log.
				WithError(err).
				Panic("Cannot initialize ECDSA private key for Ethereum")
		}
		auth := bind.NewKeyedTransactor(privKey)

		withdraw.Run(tmClient, lb, auth, relayAddr, statePath)
	},
}

func init() {
	withdrawCmd.PersistentFlags().String("eth-priv-key", "", "Ethereum private key")
	viper.BindPFlag("ethPrivKey", withdrawCmd.PersistentFlags().Lookup("eth-priv-key"))

	withdrawCmd.PersistentFlags().String("withdraw-state-path", "./state_withdraw.json", "State storage file path")
	viper.BindPFlag("withdrawStatePath", withdrawCmd.PersistentFlags().Lookup("withdraw-state-path"))
}
