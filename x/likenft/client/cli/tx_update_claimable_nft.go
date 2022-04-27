package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateMintableNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-mintable-nft [class-id] [id] [json-file-input]",
		Short: "Update mintable nft content",
		Example: `JSON file content:
{
	"uri": "",
	"uriHash": "",
	"metadata": {}
}`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argId := args[1]
			nftInput, err := readNFTInputJsonFile(args[2])
			if nftInput == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateMintableNFT(
				clientCtx.GetFromAddress().String(),
				argClassId,
				argId,
				*nftInput,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
