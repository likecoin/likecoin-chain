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

type CmdUpdateClassInput struct {
	types.ClassInput
}

func CmdUpdateClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-class [class-id] [json-file-input]",
		Short: "Update existing nft class. Only allowed when there is no token minted",
		Example: `JSON file content:
{
	"name": "",
	"symbol": "",
	"description": "",
	"uri": "",
	"uriHash": "",
	"metadata": {},
	"config": {
		"burnable": true,
		"maxSupply": 0, // 0 = unlimited
		"enablePayToMint": true,
		"mintPrice": 0 // 0 = free
	}
}`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			input, err := readCmdUpdateClassInput(args[1])
			if input == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateClass(
				clientCtx.GetFromAddress().String(),
				argClassId,
				input.ClassInput,
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
