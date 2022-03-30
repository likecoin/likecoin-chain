package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

type CmdNewClassInput struct {
	ParentType   string          `json:"parentType"`
	IscnIdPrefix string          `json:"iscnIdPrefix,omitempty"`
	Name         string          `json:"name"`
	Symbol       string          `json:"symbol"`
	Description  string          `json:"description"`
	Uri          string          `json:"uri"`
	UriHash      string          `json:"uriHash"`
	Metadata     types.JsonInput `json:"metadata"`
	Burnable     bool            `json:"burnable"`
}

func CmdNewClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-class (--account | --iscnIdPrefix=iscn://...) [json-file]",
		Short: "Create new NFT class related to an ISCN record",
		Example: `JSON file content:
{
	"name": "",
	"symbol": "",
	"description": "",
	"uri": "",
	"uriHash": "",
	"metadata": {},
	"burnable": true
}
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			input, err := readCmdNewClassInput(args[0])
			if input == nil || err != nil {
				return err
			}
			argAccount, err := cmd.Flags().GetBool("account")
			if err != nil {
				return err
			}
			argIscnIdPrefix, err := cmd.Flags().GetString("iscnIdPrefix")
			if err != nil {
				return err
			}
			if (argAccount == false && argIscnIdPrefix == "") || (argAccount == true && argIscnIdPrefix != "") {
				return fmt.Errorf("Either one of --account or --iscnIdPrefix should be set")
			}
			var argParent types.ClassParentInput
			if argAccount {
				argParent = types.ClassParentInput{
					Type: types.ClassParentType_ACCOUNT,
				}
			} else if argIscnIdPrefix != "" {
				argParent = types.ClassParentInput{
					Type:         types.ClassParentType_ISCN,
					IscnIdPrefix: argIscnIdPrefix,
				}
			}
			argName := input.Name
			argSymbol := input.Symbol
			argDescription := input.Description
			argUri := input.Uri
			argUriHash := input.UriHash
			argMetadata := input.Metadata
			argBurnable := input.Burnable

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNewClass(
				clientCtx.GetFromAddress().String(),
				argParent,
				argName,
				argSymbol,
				argDescription,
				argUri,
				argUriHash,
				argMetadata,
				argBurnable,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Bool("account", false, "Relate NFT Class to Account")
	cmd.Flags().String("iscnIdPrefix", "", "Relate NFT Class to ISCN")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
