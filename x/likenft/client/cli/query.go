package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/likecoin/likechain/x/likenft/types"
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

	cmd.AddCommand(CmdListMintableNFT())
	cmd.AddCommand(CmdShowMintableNFT())
	cmd.AddCommand(CmdMintableNFTs())

	cmd.AddCommand(CmdListClassRevealQueue())
	cmd.AddCommand(CmdShowClassRevealQueue())
	// this line is used by starport scaffolding # 1

	return cmd
}
