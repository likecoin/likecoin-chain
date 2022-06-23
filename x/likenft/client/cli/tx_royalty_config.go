package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdCreateRoyaltyConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-royalty-config [class-id] [json-file-input]",
		Short: "Create royalty config for class",
		Example: `JSON file content:
{
	"rate_basis_points": 1000, // royalty rate in basis points (1 bps = 0.01%), max 1000 (10%)
	"stakeholders": [
		{
			"account": "like1...", // address of stakeholder
			"weight": 100
		}
	]
}`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			classId := args[0]

			// Get value arguments
			argInput := args[1]
			royaltyConfigInput, err := readJsonFile[types.RoyaltyConfigInput](argInput)
			if royaltyConfigInput == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateRoyaltyConfig(
				clientCtx.GetFromAddress().String(),
				classId,
				*royaltyConfigInput,
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

func CmdUpdateRoyaltyConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-royalty-config [class-id] [json-file-input]",
		Short: "Update royalty config for class",
		Example: `JSON file content:
{
	"rate_basis_points": 1000, // royalty rate in basis points (1 bps = 0.01%), max 1000 (10%)
	"stakeholders": [
		{
			"account": "like1...", // address of stakeholder
			"weight": 100
		}
	]
}`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			classId := args[0]

			// Get value arguments
			argInput := args[1]
			royaltyConfigInput, err := readJsonFile[types.RoyaltyConfigInput](argInput)
			if royaltyConfigInput == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRoyaltyConfig(
				clientCtx.GetFromAddress().String(),
				classId,
				*royaltyConfigInput,
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

func CmdDeleteRoyaltyConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-royalty-config [class-id]",
		Short: "Delete royalty config for class",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexClassId := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteRoyaltyConfig(
				clientCtx.GetFromAddress().String(),
				indexClassId,
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
