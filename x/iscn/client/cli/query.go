package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/likecoin/likecoin-chain/v2/x/iscn/types"
)

const (
	flagFromVersion = "from-version"
	flagToVersion   = "to-version"
	flagFromSeq     = "from-sequence"
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
			fmt.Sprintf(`Query the given ISCN record. If the version part of the ISCN ID URL is not given, then the query will return the newest version of the record, or the versions controlled by the --%s and --%s flags if the flags are provided.

Example:
  $ %s query %s record iscn://likecoin-chain/yc53s4qfazn4z7doh4clxj7rugzkb2runruv4go6qsbix3vt5g2q/2
  $ %s query %s record iscn://likecoin-chain/yc53s4qfazn4z7doh4clxj7rugzkb2runruv4go6qsbix3vt5g2q
`,
				flagFromVersion, flagToVersion, version.AppName, types.ModuleName, version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromVersion, _ := cmd.Flags().GetUint64(flagFromVersion)
			toVersion, _ := cmd.Flags().GetUint64(flagToVersion)
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			iscnId, err := types.ParseIscnId(args[0])
			if err != nil {
				return err
			}
			params := types.NewQueryRecordsByIdRequest(iscnId, fromVersion, toVersion)
			res, err := queryClient.RecordsById(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	cmd.Flags().Uint64(flagFromVersion, 0, "minimum version of the records to be queried, 0 means either follow the ISCN URL or provide the latest version")
	cmd.Flags().Uint64(flagToVersion, 0, "maximum version of the records to be queried, 0 means either follow the ISCN URL or provide the latest version")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryFingerprintIscn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fingerprint [fingerprint_url]",
		Short: "Query the ISCN records for the given fingerprint.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the ISCN records for the given fingerprint. There could be more than one record associated with the given fingerprint. Note that the request is paginated, you may provide the --%s flag by the value in the previous response for querying the next page.

Example:
  $ %s query %s fingerprint hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e
  $ %s query %s fingerprint ipfs://QmNrgEMcUygbKzZeZgYFosdd27VE9KnWbyUD73bKZJ3bGi
`,
				flagFromSeq, version.AppName, types.ModuleName, version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromSeq, _ := cmd.Flags().GetUint64(flagFromSeq)
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := types.NewQueryRecordsByFingerprintRequest(args[0], fromSeq)
			res, err := queryClient.RecordsByFingerprint(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Uint64(flagFromSeq, 0, "returns the page starting from the given sequence number, for pagination together with the next_sequence field from the previous response")
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
