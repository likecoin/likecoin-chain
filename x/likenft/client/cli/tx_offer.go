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

func CmdCreateOffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-offer [class-id] [nft-id] [price] [expiration]",
		Short: "Create a new offer",
		// todo add example
		Args: cobra.ExactArgs(4),
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

			msg := types.NewMsgCreateOffer(
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

func CmdUpdateOffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-offer [class-id] [nft-id] [price] [expiration]",
		Short: "Update a offer",
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

			msg := types.NewMsgUpdateOffer(
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

func CmdDeleteOffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-offer [class-id] [nft-id]",
		Short: "Delete a offer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexClassId := args[0]
			indexNftId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteOffer(
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
