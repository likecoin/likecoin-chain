package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListClassesByAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-index",
		Short: "Enumerate all ISCN to NFT classes relation records",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryClassesByAccountIndexRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ClassesByAccountIndex(context.Background(), params)
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

func CmdShowClassesByAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-classes [account]",
		Short: "Query NFT classes related to an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argAccount := args[0]

			params := &types.QueryClassesByAccountRequest{
				Account:    argAccount,
				Pagination: pageReq,
			}

			res, err := queryClient.ClassesByAccount(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	// Copy from sdk `flags.AddPaginationFlagsToCmd(cmd, cmd.Use)`
	// We needed to override default values of limit and count total
	cmd.Flags().Uint64(flags.FlagPage, 1, fmt.Sprintf("pagination page of %s to query for. This sets offset to a multiple of limit", cmd.Use))
	cmd.Flags().String(flags.FlagPageKey, "", fmt.Sprintf("pagination page-key of %s to query for", cmd.Use))
	cmd.Flags().Uint64(flags.FlagOffset, 0, fmt.Sprintf("pagination offset of %s to query for", cmd.Use))
	// TODO refactor this constant in oursky/likecoin-chain#98
	cmd.Flags().Uint64(flags.FlagLimit, 20, fmt.Sprintf("pagination limit of %s to query for", cmd.Use))
	cmd.Flags().Bool(flags.FlagCountTotal, true, fmt.Sprintf("count total number of records in %s to query for", cmd.Use))
	cmd.Flags().Bool(flags.FlagReverse, false, "results are sorted in descending order")

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
