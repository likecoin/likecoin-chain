package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-listing",
		Short: "list all listing",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllListingRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListingAll(context.Background(), params)
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

func CmdShowListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-listing [class-id] [nft-id] [seller]",
		Short: "shows a listing",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argClassId := args[0]
			argNftId := args[1]
			argSeller := args[2]

			params := &types.QueryGetListingRequest{
				ClassId: argClassId,
				NftId:   argNftId,
				Seller:  argSeller,
			}

			res, err := queryClient.Listing(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
