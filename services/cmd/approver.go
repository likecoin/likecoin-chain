package cmd

import (
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/likecoin/likechain/services/approver"
)

var approverCmd = &cobra.Command{
	Use:   "approver",
	Short: "run the deposit approver service",
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
		approverDelay := viper.GetInt64("approverDelay")
		if approverDelay <= 0 {
			panic("Wrong approver delay")
		}
		approver.Run(tmClient, ethClient, tokenAddr, relayAddr, privKey, uint64(approverDelay))
	},
}
