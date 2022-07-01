package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdNewClass())
	cmd.AddCommand(CmdUpdateClass())
	cmd.AddCommand(CmdMintNFT())
	cmd.AddCommand(CmdBurnNFT())
	cmd.AddCommand(CmdCreateBlindBoxContent())
	cmd.AddCommand(CmdUpdateBlindBoxContent())
	cmd.AddCommand(CmdDeleteBlindBoxContent())
	cmd.AddCommand(CmdCreateOffer())
	cmd.AddCommand(CmdUpdateOffer())
	cmd.AddCommand(CmdDeleteOffer())
	cmd.AddCommand(CmdCreateListing())
	cmd.AddCommand(CmdUpdateListing())
	cmd.AddCommand(CmdDeleteListing())
	cmd.AddCommand(CmdSellNFT())
	cmd.AddCommand(CmdBuyNFT())
	cmd.AddCommand(CmdCreateRoyaltyConfig())
	cmd.AddCommand(CmdUpdateRoyaltyConfig())
	cmd.AddCommand(CmdDeleteRoyaltyConfig())
	// this line is used by starport scaffolding # 1

	return cmd
}
