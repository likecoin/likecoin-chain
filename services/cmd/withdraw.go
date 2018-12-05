package cmd

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/likecoin/likechain/services/withdraw"
)

var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "run the withdraw-relay service",
	Run: func(cmd *cobra.Command, args []string) {
		tmClient := tmRPC.NewHTTP(viper.GetString("tmEndPoint"), "/websocket")
		contractAddr := common.HexToAddress(viper.GetString("relayContractAddr"))
		ethEndPoint := viper.GetString("ethEndPoint")
		ethClient, err := ethclient.Dial(ethEndPoint)
		if err != nil {
			log.
				WithField("eth_endpoint", ethEndPoint).
				WithError(err).
				Panic("Cannot initialize Ethereum endpoint")
		}
		privKeyBytes := common.Hex2Bytes(viper.GetString("ethPrivKey"))
		privKey, err := ethCrypto.ToECDSA(privKeyBytes)
		if err != nil {
			log.
				WithError(err).
				Panic("Cannot initialize ECDSA private key for Ethereum")
		}
		auth := bind.NewKeyedTransactor(privKey)

		withdraw.Run(tmClient, ethClient, auth, contractAddr)
	},
}
