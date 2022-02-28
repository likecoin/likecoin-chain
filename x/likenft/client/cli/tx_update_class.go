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

type CmdUpdateClassInput struct {
	Name        string          `json:"name"`
	Symbol      string          `json:"symbol"`
	Description string          `json:"description"`
	Uri         string          `json:"uri"`
	UriHash     string          `json:"uriHash"`
	Metadata    types.JsonInput `json:"metadata"`
	Burnable    bool            `json:"burnable"`
}

func CmdUpdateClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-class [class-id] [json-file-input]",
		Short: "Update existing nft class. Only allowed when there is no token minted",
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
			argClassId := args[0]
			input, err := readCmdUpdateClassInput(args[1])
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

			msg := types.NewMsgUpdateClass(
				clientCtx.GetFromAddress().String(),
				argClassId,
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
