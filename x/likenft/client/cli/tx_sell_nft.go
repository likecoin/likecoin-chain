package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSellNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell-nft [class-id] [nft-id] [buyer] [price] (--full-pay-to-royalty)",
		Short: "Broadcast message SellNFT",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argNftId := args[1]
			argBuyer := args[2]
			argPrice, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			flagFullPayToRoyalty, err := cmd.Flags().GetBool("full-pay-to-royalty")
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSellNFT(
				clientCtx.GetFromAddress().String(),
				argClassId,
				argNftId,
				argBuyer,
				argPrice,
				flagFullPayToRoyalty,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Bool("full-pay-to-royalty", false, "Pay full price to royalty")

	return cmd
}
