package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListOffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offer-index",
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
		Use:   "offer [class-id] [nft-id] [buyer]",
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

func CmdOffersByClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "class-offers [class-id]",
		Short: "Query offers by class",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqClassId := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryOffersByClassRequest{

				ClassId: reqClassId,
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params.Pagination = pageReq

			res, err := queryClient.OffersByClass(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdOffersByNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-offers [class-id] [nft-id]",
		Short: "Query offers by nft",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqClassId := args[0]
			reqNftId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryOffersByNFTRequest{

				ClassId: reqClassId,
				NftId:   reqNftId,
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params.Pagination = pageReq

			res, err := queryClient.OffersByNFT(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
