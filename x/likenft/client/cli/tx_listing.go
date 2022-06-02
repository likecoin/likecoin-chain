package cli

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdCreateListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-listing [class-id] [nft-id] [price] [expiration]",
		Short: "Create a new listing",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexClassId := args[0]
			indexNftId := args[1]

			// Get value arguments
			argPrice, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			argExpiration, err := time.Parse(time.RFC3339, args[3])
			if err != nil {
				return nil
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateListing(
				clientCtx.GetFromAddress().String(),
				indexClassId,
				indexNftId,
				argPrice,
				argExpiration,
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

func CmdUpdateListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-listing [class-id] [nft-id] [price] [expiration]",
		Short: "Update a listing",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			indexClassId := args[0]
			indexNftId := args[1]

			// Get value arguments
			argPrice, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			argExpiration, err := time.Parse(time.RFC3339, args[3])
			if err != nil {
				return nil
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateListing(
				clientCtx.GetFromAddress().String(),
				indexClassId,
				indexNftId,
				argPrice,
				argExpiration,
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

func CmdDeleteListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-listing [class-id] [nft-id]",
		Short: "Delete a listing",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexClassId := args[0]
			indexNftId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteListing(
				clientCtx.GetFromAddress().String(),
				indexClassId,
				indexNftId,
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
