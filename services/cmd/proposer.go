package cmd

import (
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/likecoin/likechain/services/proposer"
)

var proposerCmd = &cobra.Command{
	Use:   "proposer",
	Short: "run the deposit proposer service",
	Run: func(cmd *cobra.Command, args []string) {
		tmClient := tmRPC.NewHTTP(viper.GetString("tmEndPoint"), "/websocket")
		tokenAddr := common.HexToAddress(viper.GetString("tokenContractAddr"))
		relayAddr := common.HexToAddress(viper.GetString("relayContractAddr"))
		ethClient, err := ethclient.Dial(viper.GetString("ethEndPoint"))
		if err != nil {
			panic(err)
		}
		privKeyBytes := common.Hex2Bytes(viper.GetString("tmPrivKey"))
		privKey, err := ethCrypto.ToECDSA(privKeyBytes)
		if err != nil {
			panic(err)
		}
		proposerDelay := viper.GetInt64("proposerDelay")
		if proposerDelay <= 0 {
			panic("Wrong proposer delay")
		}
		proposer.Run(tmClient, ethClient, tokenAddr, relayAddr, privKey, uint64(proposerDelay))
	},
}
