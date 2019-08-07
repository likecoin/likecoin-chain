package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/whitelist/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	whitelistQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the whitelist module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	whitelistQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryWhitelist(queryRoute, cdc),
		GetCmdQueryApprover(queryRoute, cdc),
	)...)

	return whitelistQueryCmd

}

// GetCmdQueryWhitelist implements the validator whitelist query command.
func GetCmdQueryWhitelist(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "whitelist",
		Short: "Query the current validator whitelist",
		Long: strings.TrimSpace(`Query the current validator whitelist:

$ likecli query whitelist whitelist
`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", storeName, types.QueryWhitelist))
			if err != nil {
				return err
			}

			whitelist := types.Whitelist{}
			if len(res) > 0 {
				cdc.UnmarshalJSON(res, &whitelist)
			}

			return cliCtx.PrintOutput(whitelist)
		},
	}
}

// GetCmdQueryApprover implements the validator whitelist approver query command.
func GetCmdQueryApprover(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "approver",
		Short: "Query the validator whitelist approver",
		Long: strings.TrimSpace(`Query the validator whitelist approver:

$ likecli query whitelist approver
`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", storeName, types.QueryApprover))
			if err != nil {
				return err
			}

			approver := sdk.AccAddress{}
			if len(res) > 0 {
				cdc.UnmarshalJSON(res, &approver)
			}

			return cliCtx.PrintOutput(approver)
		},
	}
}
