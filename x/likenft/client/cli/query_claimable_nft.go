package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListClaimableNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-claimable-nft",
		Short: "list all claimableNFT",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllClaimableNFTRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ClaimableNFTAll(context.Background(), params)
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

func CmdShowClaimableNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-claimable-nft [class-id] [claimable-id]",
		Short: "shows a claimableNFT",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argClassId := args[0]
			argId := args[1]

			params := &types.QueryGetClaimableNFTRequest{
				ClassId: argClassId,
				Id:      argId,
			}

			res, err := queryClient.ClaimableNFT(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
