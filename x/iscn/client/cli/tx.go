package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/likecoin/likechain/x/iscn/types"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "ISCN module related transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		NewCreateIscnTxCmd(),
		NewUpdateIscnTxCmd(),
		NewChangeIscnOwnershipTxCmd(),
	)
	return txCmd
}

func readIscnRecordFile(path string) (*types.IscnRecord, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	record := types.IscnRecord{}
	err = json.Unmarshal(contents, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func NewCreateIscnTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-iscn [iscn_record_json_file]",
		Short: `Create an ISCN record from the given file. The ISCN record will be registered on chain and assigned an ISCN ID.`,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create an ISCN record on the chain. The record and the related parameters are given in the JSON file.

Example:
$ %s tx iscn create-iscn record.json --from mykey

Content of record.json:

{
  "recordNotes": "Some Notes",
  "contentFingerprints": [
    "hash://sha256/..."
  ],
  "stakeholders": [
    ...
  ],
  "contentMetadata": {
    ...
  }
}

Where: 

"recordNotes" is a string representing the notes of this record (like Git commit logs).
"contentFingerprints" must contain URLs representing the fingerprints of the content.
"stakeholders" must contains valid JSON values.
"contentMetadata" must be a valid JSON value.
`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			record, err := readIscnRecordFile(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgCreateIscnRecord(clientCtx.GetFromAddress(), record)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewUpdateIscnTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-iscn [iscn_id_url] [iscn_record_json_file]",
		Short: `Update an ISCN record using the record from the given file.`,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update an ISCN record on the chain. The ISCN ID to be updated needs to be specified in the CLI parameter, where the new record and the related parameters are given in the JSON file.

Example:
$ %s tx iscn update-iscn "iscn://likecoin-chain/yc53s4qfazn4z7doh4clxj7rugzkb2runruv4go6qsbix3vt5g2q/1" record.json --from mykey

The ISCN ID needs to be a URL representing the newest version of the record to be updated, i.e. the scheme must be "iscn://", the numeric part at the end must be the existing latest version of the record on the chain.

Content of record.json:

{
  "recordNotes": "Updating new stakeholders",
  "contentFingerprints": [
    "hash://sha256/..."
  ],
  "stakeholders": [
    ...
  ],
  "contentMetadata": {
    ...
  }
}

Where: 

"recordNotes" is a string representing the notes of this update (like Git commit logs).
"contentFingerprints" must contain URLs representing the fingerprints of the content.
"stakeholders" must contains valid JSON values.
"contentMetadata" must be a valid JSON value.
`, version.AppName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			iscnId, err := types.ParseIscnId(args[0])
			if err != nil {
				return err
			}
			record, err := readIscnRecordFile(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateIscnRecord(clientCtx.GetFromAddress(), iscnId, record)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewChangeIscnOwnershipTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-iscn-ownership [iscn_id_url] [new_owner_address]",
		Short: `Change the ownership of an ISCN record, so that the new owner can update the record afterwards.`,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Change the ownership of an ISCN record on the chain, so that the new owner can update the record afterwards.

Example:
$ %s tx iscn change-iscn-ownership "iscn://likecoin-chain/yc53s4qfazn4z7doh4clxj7rugzkb2runruv4go6qsbix3vt5g2q/1" cosmos1ww3qews2y5jxe8apw2zt8stqqrcu2tptejfwaf --from mykey

The ISCN ID needs to be a URL representing the newest version of the record to be updated, i.e. the scheme must be "iscn://", the numeric part at the end must be the existing latest version of the record on the chain.`, version.AppName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			iscnId, err := types.ParseIscnId(args[0])
			if err != nil {
				return err
			}
			newOwner, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgChangeIscnRecordOwnership(clientCtx.GetFromAddress(), iscnId, newOwner)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
