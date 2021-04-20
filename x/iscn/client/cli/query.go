package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/likecoin/likechain/x/iscn/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the ISCN module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmdQueryIscnRecord(),
		GetCmdQueryFingerprintIscn(),
		GetCmdQueryParams(),
	)
	return cmd
}

func GetCmdQueryIscnRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [iscn_id_url]",
		Short: "Query the given ISCN record.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the given ISCN record. If the version part of the ISCN ID URL is not given, then the query will return the newest version of the record.

Example:
  $ %s query %s record iscn://likecoin-chain/yc53s4qfazn4z7doh4clxj7rugzkb2runruv4go6qsbix3vt5g2q/2
  $ %s query %s record iscn://likecoin-chain/yc53s4qfazn4z7doh4clxj7rugzkb2runruv4go6qsbix3vt5g2q
`,
				version.AppName, types.ModuleName, version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			iscnId, err := types.ParseIscnId(args[0])
			if err != nil {
				return err
			}
			params := types.NewQueryIscnRecordsRequestByIscnId(iscnId)
			res, err := queryClient.IscnRecords(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryFingerprintIscn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fingerprint [fingerprint_url]",
		Short: "Query the ISCN records for the given fingerprint.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the ISCN records for the given fingerprint. There could be more than one record associated with the given fingerprint.

Example:
  $ %s query %s fingerprint hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e
  $ %s query %s fingerprint ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi
`,
				version.AppName, types.ModuleName, version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := types.NewQueryIscnRecordsRequestByFingerprint(args[0])
			res, err := queryClient.IscnRecords(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	// TODO: pagination?
	return cmd
}

func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the chain parameters of the ISCN module, including the ISCN registry ID and fee.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := types.NewQueryParamsRequest()
			res, err := queryClient.Params(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
