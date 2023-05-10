package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [class-id] (--id [id] --input [json-file-input])",
		Short: "Mint NFT under a class",
		Example: `--id and --input required for minting under normal class, ignored for blind box class
JSON file content:
{
	"uri": "",
	"uri_hash": "",
	"metadata": {}
}`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argId, err := cmd.Flags().GetString("id")
			if err != nil {
				return err
			}

			argInput, err := cmd.Flags().GetString("input")
			if err != nil {
				return err
			}
			var nftInput *types.NFTInput
			if argInput != "" {
				nftInput, err = readJsonFile[types.NFTInput](argInput)
				if nftInput == nil || err != nil {
					return err
				}
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(
				clientCtx.GetFromAddress().String(),
				argClassId,
				argId,
				nftInput,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("id", "", "NFT ID")
	cmd.Flags().String("input", "", "Path to json file containing NFT Input data")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
