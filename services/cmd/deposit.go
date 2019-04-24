package cmd

import (
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/likecoin/likechain/services/deposit"
	"github.com/likecoin/likechain/services/eth"
)

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "run the deposit service",
	Run: func(cmd *cobra.Command, args []string) {
		tmEndPoint := viper.GetString("tmEndPoint")
		ethEndPoints := viper.GetStringSlice("ethEndPoints")
		tokenAddr := common.HexToAddress(viper.GetString("tokenContractAddr"))
		relayAddr := common.HexToAddress(viper.GetString("relayContractAddr"))
		blockDelay := viper.GetInt64("blockDelay")
		statePath := viper.GetString("depositStatePath")
		minTrialPerClient := viper.GetInt("ethMinTrialPerClient")
		maxTrialCount := viper.GetInt("ethMaxTrialCount")
		log.
			WithField("tm_endpoint", tmEndPoint).
			WithField("eth_endpoints", ethEndPoints).
			WithField("token_addr", tokenAddr).
			WithField("relay_addr", relayAddr).
			WithField("block_delay", blockDelay).
			WithField("state_path", statePath).
			WithField("min_trial_per_client", minTrialPerClient).
			WithField("max_trial_count", maxTrialCount).
			Debug("Read deposit config and parameters")

		tmClient := tmRPC.NewHTTP(tmEndPoint, "/websocket")
		lb := eth.NewLoadBalancer(ethEndPoints, uint(minTrialPerClient), uint(maxTrialCount))
		privKeyBytes := common.Hex2Bytes(viper.GetString("tmPrivKey"))
		privKey, err := ethCrypto.ToECDSA(privKeyBytes)
		if err != nil {
			log.
				WithError(err).
				Panic("Cannot initialize ECDSA private key for LikeChain")
		}
		if blockDelay <= 0 {
			log.
				WithField("block_delay", blockDelay).
				Panic("Invalid block delay value")
		}
		deposit.Run(&deposit.Config{
			TMClient:     tmClient,
			LoadBalancer: lb,
			TokenAddr:    tokenAddr,
			RelayAddr:    relayAddr,
			TMPrivKey:    privKey,
			BlockDelay:   blockDelay,
			StatePath:    statePath,
		})
	},
}

func init() {
	depositCmd.PersistentFlags().String("token-addr", "", "Ethereum address of the token contract")
	viper.BindPFlag("tokenContractAddr", depositCmd.PersistentFlags().Lookup("token-addr"))

	depositCmd.PersistentFlags().String("tm-priv-key", "", "Tendermint private key")
	viper.BindPFlag("tmPrivKey", depositCmd.PersistentFlags().Lookup("tm-priv-key"))

	depositCmd.PersistentFlags().Int("block-delay", 12, "Ethereum block delay before confirm")
	viper.BindPFlag("blockDelay", depositCmd.PersistentFlags().Lookup("block-delay"))

	depositCmd.PersistentFlags().String("deposit-state-path", "./state_deposit.json", "State storage file path")
	viper.BindPFlag("depositStatePath", depositCmd.PersistentFlags().Lookup("deposit-state-path"))
}
