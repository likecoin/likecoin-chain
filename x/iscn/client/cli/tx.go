package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/likecoin/likechain/x/iscn/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	iscnTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "ISCN transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	iscnTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateIscn(cdc),
	)...)

	return iscnTxCmd
}

// GetCmdCreateIscn implements the create ISCN command
func GetCmdCreateIscn(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-iscn",
		Short: "create an ISCN record",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// TODO: read record from file?
			// record := types.IscnRecord{}
			record := types.IscnRecord{
				Version: 1,
				Stakeholders: []types.Stakeholder{
					{
						Type:  "author",
						Id:    "testing-author-asdf",
						Stake: 1337,
					},
					{
						Type:  "whoever",
						Id:    "testing-whoever-asdf",
						Stake: 2345,
					},
				},
				Timestamp: 1234567890,
				Parent:    "cid-parent",
				Right: []types.Right{
					{
						Holder: "Chung",
						Type:   "license",
						Terms:  "cc-by-sa-4.0",
						Period: types.Period{
							To: "2030-01-01 12:34:56Z",
						},
						Territory: "Hong Kong",
					},
				},
				Content: types.IscnContent{
					Type:        "article",
					Source:      "https://nnkken.github.io/about",
					Fingerprint: "hash://sha256/3de89366df13254ee59a3a4ff1ab1471cd5372d7004eeb72718ab8e72e9168fa",
					Feature:     "",
					Edition:     "3",
					Title:       "About nnkken",
					Description: "About page of nnkken's blog",
					Tags:        []string{"about", "nnkken"},
				},
			}
			from := cliCtx.GetFromAddress()

			msg := types.NewMsgCreateIscn(from, record)
			s, err := cdc.MarshalJSONIndent(msg, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(s))

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}
