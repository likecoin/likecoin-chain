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

type CmdMintNFTInput struct {
	Uri      string          `json:"uri"`
	UriHash  string          `json:"uriHash`
	Metadata types.JsonInput `json:"metadata"`
}

func CmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-nft [class-id] [id] [json-file-input]",
		Short: "Mint NFT under a class",
		Example: `JSON file content:
{
	"uri": "",
	"uriHash": "",
	"metadata": {}
}`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argId := args[1]
			input, err := readCmdMintNFTInput(args[2])
			if input == nil || err != nil {
				return err
			}
			argUri := input.Uri
			argUriHash := input.UriHash
			argMetadata := input.Metadata

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(
				clientCtx.GetFromAddress().String(),
				argClassId,
				argId,
				argUri,
				argUriHash,
				argMetadata,
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
