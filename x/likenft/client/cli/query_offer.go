package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListOffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-offer",
		Short: "list all offer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryOfferIndexRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.OfferIndex(context.Background(), params)
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

func CmdShowOffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-offer [class-id] [nft-id] [buyer]",
		Short: "shows a offer",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argClassId := args[0]
			argNftId := args[1]
			argBuyer := args[2]

			params := &types.QueryOfferRequest{
				ClassId: argClassId,
				NftId:   argNftId,
				Buyer:   argBuyer,
			}

			res, err := queryClient.Offer(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
