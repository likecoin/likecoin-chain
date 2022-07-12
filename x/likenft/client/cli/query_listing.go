package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listing-index",
		Short: "list all listing",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryListingIndexRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ListingIndex(context.Background(), params)
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
		Use:   "listing [class-id] [nft-id] [seller]",
		Short: "shows a listing",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argClassId := args[0]
			argNftId := args[1]
			argSeller := args[2]

			params := &types.QueryListingRequest{
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

func CmdListingsByClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "class-listings [class-id]",
		Short: "Query listings by class",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqClassId := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryListingsByClassRequest{

				ClassId: reqClassId,
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params.Pagination = pageReq

			res, err := queryClient.ListingsByClass(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListingsByNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-listings [class-id] [nft-id]",
		Short: "Query listings by nft",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqClassId := args[0]
			reqNftId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryListingsByNFTRequest{

				ClassId: reqClassId,
				NftId:   reqNftId,
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params.Pagination = pageReq

			res, err := queryClient.ListingsByNFT(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
