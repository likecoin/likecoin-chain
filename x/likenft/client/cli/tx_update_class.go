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
}`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			classInput, err := readClassInputJsonFile(args[1])
			if classInput == nil || err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateClass(
				clientCtx.GetFromAddress().String(),
				argClassId,
				*classInput,
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
