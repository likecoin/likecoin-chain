package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

type CmdNewClassInput struct {
	Name        string          `json:"name"`
	Symbol      string          `json:"symbol"`
	Description string          `json:"description"`
	Uri         string          `json:"uri"`
	UriHash     string          `json:"uriHash"`
	Metadata    types.JsonInput `json:"metadata"`
	Burnable    bool            `json:"burnable"`
}

func CmdNewClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-class [iscn-id-prefix] [json-file]",
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
}`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argIscnIdPrefix := args[0]
			input, err := readCmdNewClassInput(args[1])
			if input == nil || err != nil {
				return err
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
				argIscnIdPrefix,
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

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
