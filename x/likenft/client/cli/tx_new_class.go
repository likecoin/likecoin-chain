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
	"uri_hash": "",
	"metadata": {},
	"config": {
		"burnable": true,
		"max_supply": 0, // 0 = unlimited
		"blind_box_config": { // null = not using blind box feature
			"mint_periods": [
				{
					"start_time": "2022-01-01T00:00:00Z",
					"allowed_addresses": ["like1..."], // [] = public
					"mint_price": 0 // 0 = free
				}
			],
			"reveal_time": "2022-02-01T00:00:00Z"
		},
		"royalty_basis_points": 0 // each base point is 0.01%, max 10% / 1000 bps
	}
}
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			classInput, err := readJsonFile[types.ClassInput](args[0])
			if classInput == nil || err != nil {
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

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNewClass(
				clientCtx.GetFromAddress().String(),
				argParent,
				*classInput,
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
