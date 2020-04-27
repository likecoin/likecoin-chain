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
		GetCmdAddEntity(cdc),
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

			// TODO: read from file?
			kernelBz := []byte{} // TODO
			from := cliCtx.GetFromAddress()

			msg := types.NewMsgCreateIscn(from, kernelBz)

			/* vvv debug vvv */
			s, err := cdc.MarshalJSONIndent(msg, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(s))
			/* ^^^ debug ^^^ */

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}

func GetCmdAddEntity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-entity",
		Short: "add an entity",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// TODO: read from file?
			entityBz := []byte{}
			from := cliCtx.GetFromAddress()

			msg := types.NewMsgAddEntity(from, entityBz)
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
