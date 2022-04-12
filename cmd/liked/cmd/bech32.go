package cmd

// Based on: https://github.com/osmosis-labs/osmosis/blob/v7.1.0/cmd/osmosisd/cmd/bech32.go

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

var (
	flagToPrefix = "prefix"
)

// get cmd to convert any bech32 address to an like prefix
func ConvertPrefixCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-prefix [bech32 address]",
		Short: "Convert any bech32 address to the like prefix",
		Long: `Convert any bech32 address to the like prefix
Especially useful for converting cosmos prefix to like prefix
Example:
	liked convert-prefix cosmos1ey69r37gfxvxg62sh4r0ktpuc46pzjrmz29g45
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			toPrefix, err := cmd.Flags().GetString(flagToPrefix)
			if err != nil {
				return err
			}

			_, bz, err := bech32.DecodeAndConvert(args[0])
			if err != nil {
				return err
			}

			bech32Addr, err := bech32.ConvertAndEncode(toPrefix, bz)
			if err != nil {
				panic(err)
			}

			cmd.Println(bech32Addr)

			return nil
		},
	}

	cmd.Flags().StringP(flagToPrefix, "p", "like", "Bech32 prefix to encode to")

	return cmd
}
