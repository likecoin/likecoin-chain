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

func CmdDeleteBlindBoxContent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-blind-box-content [class-id] [content-id]",
		Short: "Delete blind box content",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteBlindBoxContent(
				clientCtx.GetFromAddress().String(),
				argClassId,
				argId,
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
