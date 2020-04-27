package cli

import (
	"fmt"
	"strings"

	gocid "github.com/ipfs/go-cid"
	"github.com/polydawn/refmt/json"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/likecoin/likechain/x/iscn/types"
	"github.com/multiformats/go-multibase"
)

type printableIscnMap types.RawIscnMap

func (m printableIscnMap) String() string {
	bz, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(bz)
}

type printableCID types.CID

func (cid printableCID) String() string {
	s, err := types.CID(cid).StringOfBase(types.CidMbaseEncoder.Encoding())
	if err != nil {
		return ""
	}
	return s
}

type printableString string

func (s printableString) String() string {
	return string(s)
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	iscnQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the ISCN module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	iscnQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryIscnKernel(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryCID(queryRoute, cdc),
	)...)

	return iscnQueryCmd

}

func GetCmdQueryIscnKernel(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "kernel [iscn-id]",
		Short: "Query an ISCN kernel by ID",
		Long: strings.TrimSpace(`Query an ISCN kernel by ID:

$ likecli query iscn kernel xxxxxxx
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			idStr := args[0]
			// TODO: parse by ISCN ID format: 1/xxxxx...
			_, id, err := multibase.Decode(idStr)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", storeName, types.QueryIscnKernel), id)
			if err != nil {
				return err
			}

			_, kernelCID, err := gocid.CidFromBytes(res)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(printableCID(kernelCID))
		},
	}
}

func GetCmdQueryCID(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cid [cid]",
		Short: "Query a CID",
		Long: strings.TrimSpace(`Query a CID:

$ likecli query iscn cid xxxxxxx
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cidStr := args[0]
			cid, err := gocid.Decode(cidStr)
			if err != nil {
				return err
			}
			bz, err := cdc.MarshalJSON(cid.Bytes())
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", storeName, types.QueryCidBlockGet), bz)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(printableString(res))
		},
	}
}

func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the ISCN module params",
		Long: strings.TrimSpace(`Query the ISCN module params:

$ likecli query iscn params
`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", storeName, types.QueryParams))
			if err != nil {
				return err
			}

			params := types.Params{}
			if len(res) > 0 {
				cdc.UnmarshalJSON(res, &params)
			}

			return cliCtx.PrintOutput(params)
		},
	}
}
