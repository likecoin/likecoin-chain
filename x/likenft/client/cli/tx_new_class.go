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

func CmdNewClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-class [iscn-id-prefix] [name] [symbol] [description] [uri] [uri-hash] [metadata]",
		Short: "Broadcast message NewClass",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argIscnIdPrefix := args[0]
			argName := args[1]
			argSymbol := args[2]
			argDescription := args[3]
			argUri := args[4]
			argUriHash := args[5]
			argMetadata := args[6]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNewClass(
				clientCtx.GetFromAddress().String(),
				argIscnIdPrefix,
				argName,
				argSymbol,
				argDescription,
				argUri,
				argUriHash,
				types.JsonInput(argMetadata),
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
