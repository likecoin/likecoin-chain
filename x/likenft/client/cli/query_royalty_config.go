package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListRoyaltyConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "royalty-config-index",
		Short: "list all royalty config by class",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRoyaltyConfigIndexRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RoyaltyConfigIndex(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowRoyaltyConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "royalty-config [class-id]",
		Short: "shows royalty config of a class",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argClassId := args[0]

			params := &types.QueryRoyaltyConfigRequest{
				ClassId: argClassId,
			}

			res, err := queryClient.RoyaltyConfig(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
