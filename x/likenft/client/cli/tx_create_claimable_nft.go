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

func CmdCreateBlindBoxContent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-blind-box-content [class-id] [content-id] [json-file-input]",
		Short: "Create blind box content",
		Example: `JSON file content:
{
	"uri": "",
	"uri_hash": "",
	"metadata": {}
}`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argId := args[1]
			nftInput, err := readJsonFile[types.NFTInput](args[2])
			if nftInput == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateBlindBoxContent(
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
