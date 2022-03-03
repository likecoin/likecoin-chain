package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

func CmdListClassesByISCN() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Enumerate all ISCN to NFT classes relation records",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryClassesByISCNIndexRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ClassesByISCNIndex(context.Background(), params)
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

func CmdShowClassesByISCN() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "classes [iscn-id-prefix]",
		Short: "Query NFT classes related to an ISCN",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIscnIdPrefix := args[0]

			params := &types.QueryClassesByISCNRequest{
				IscnIdPrefix: argIscnIdPrefix,
			}

			res, err := queryClient.ClassesByISCN(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
