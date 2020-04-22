package cli

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/iscn/types"
)

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
		GetCmdQueryIscnRecord(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryAuthor(queryRoute, cdc),
	)...)

	return iscnQueryCmd

}

func GetCmdQueryIscnRecord(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "record [iscn-id]",
		Short: "Query an ISCN record by ID",
		Long: strings.TrimSpace(`Query an ISCN record by ID:

$ likecli query iscn record xxxxxxx
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			idStr := args[0]
			id, err := base64.URLEncoding.DecodeString(idStr)
			if err != nil {
				return err
			}
			queryData := types.QueryRecordParams{
				Id: id,
			}
			bz, err := cdc.MarshalJSON(queryData)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", storeName, types.QueryIscnRecord), bz)
			if err != nil {
				return err
			}

			record := types.IscnRecord{}
			if len(res) > 0 {
				cdc.UnmarshalJSON(res, &record)
			}

			return cliCtx.PrintOutput(record)
		},
	}
}

func GetCmdQueryAuthor(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "author [iscn-id]",
		Short: "Query an author by CID",
		Long: strings.TrimSpace(`Query an author CID:

$ likecli query iscn author xxxxxxx
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			idStr := args[0]
			cid, err := base64.URLEncoding.DecodeString(idStr)
			if err != nil {
				return err
			}
			queryData := types.QueryAuthorParams{
				Cid: cid,
			}
			bz, err := cdc.MarshalJSON(queryData)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", storeName, types.QueryAuthor), bz)
			if err != nil {
				return err
			}

			author := types.Author{}
			if len(res) > 0 {
				cdc.UnmarshalJSON(res, &author)
			}

			return cliCtx.PrintOutput(author)
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

			approver := sdk.AccAddress{}
			if len(res) > 0 {
				cdc.UnmarshalJSON(res, &approver)
			}

			return cliCtx.PrintOutput(approver)
		},
	}
}
