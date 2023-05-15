package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group likenft queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListClassesByISCN())
	cmd.AddCommand(CmdShowClassesByISCN())
	cmd.AddCommand(CmdISCNByClass())

	cmd.AddCommand(CmdListClassesByAccount())
	cmd.AddCommand(CmdShowClassesByAccount())
	cmd.AddCommand(CmdAccountByClass())

	cmd.AddCommand(CmdListBlindBoxContent())
	cmd.AddCommand(CmdShowBlindBoxContent())
	cmd.AddCommand(CmdBlindBoxContents())

	cmd.AddCommand(CmdListOffer())
	cmd.AddCommand(CmdShowOffer())
	cmd.AddCommand(CmdOffersByClass())

	cmd.AddCommand(CmdOffersByNFT())

	cmd.AddCommand(CmdListListing())
	cmd.AddCommand(CmdShowListing())
	cmd.AddCommand(CmdListingsByClass())

	cmd.AddCommand(CmdListingsByNFT())

	cmd.AddCommand(CmdListRoyaltyConfig())
	cmd.AddCommand(CmdShowRoyaltyConfig())
	// this line is used by starport scaffolding # 1

	return cmd
}
