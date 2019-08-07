package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/likecoin/likechain/x/whitelist/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	whitelistTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Whitelist transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	whitelistTxCmd.AddCommand(client.PostCommands(
		GetCmdSetWhitelist(cdc),
	)...)

	return whitelistTxCmd
}

// GetCmdSetWhitelist implements the set validator whitelist command
func GetCmdSetWhitelist(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-whitelist",
		Short: "set validator whitelist",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddrs := []sdk.ValAddress{}
			for _, valAddrStr := range args {
				valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
				if err != nil {
					return err
				}
				valAddrs = append(valAddrs, valAddr)
			}
			approverAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgSetWhitelist(approverAddr, valAddrs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}
