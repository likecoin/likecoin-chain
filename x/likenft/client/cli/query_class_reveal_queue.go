package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListClassRevealQueue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-class-reveal-queue",
		Short: "list all classRevealQueue",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllClassRevealQueueRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ClassRevealQueueAll(context.Background(), params)
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

func CmdShowClassRevealQueue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-class-reveal-queue [reveal-time] [class-id]",
		Short: "shows a classRevealQueue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argRevealTime := args[0]
			argClassId := args[1]

			params := &types.QueryGetClassRevealQueueRequest{
				RevealTime: argRevealTime,
				ClassId:    argClassId,
			}

			res, err := queryClient.ClassRevealQueue(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
