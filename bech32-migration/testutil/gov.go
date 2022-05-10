package testutil

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/likecoin/likecoin-chain/v2/bech32-migration/utils"
)

func AssertGovAddressBech32(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) bool {
	ok := true
	utils.IterateStoreByPrefix(ctx, storeKey, types.VotesKeyPrefix, func(bz []byte) []byte {
		vote := types.Vote{}
		cdc.MustUnmarshal(bz, &vote)
		if !strings.HasPrefix(vote.Voter, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: vote.Voter", vote.Voter))
			ok = false
		}
		return bz
	})
	utils.IterateStoreByPrefix(ctx, storeKey, types.DepositsKeyPrefix, func(bz []byte) []byte {
		deposit := types.Deposit{}
		cdc.MustUnmarshal(bz, &deposit)
		if !strings.HasPrefix(deposit.Depositor, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: deposit.Depositor", deposit.Depositor))
			ok = false
		}
		return bz
	})
	utils.IterateStoreByPrefix(ctx, storeKey, types.ProposalsKeyPrefix, func(bz []byte) []byte {
		proposal := types.Proposal{}
		cdc.MustUnmarshal(bz, &proposal)
		content := proposal.GetContent()
		communityPoolSpendProposal, ok := content.(*distrtypes.CommunityPoolSpendProposal)
		if !ok {
			return bz
		}
		if !strings.HasPrefix(communityPoolSpendProposal.Recipient, "like") {
			ctx.Logger().Info(fmt.Sprintf("Bad prefix found: %s, interface type: communityPoolSpendProposal.Recipient", communityPoolSpendProposal.Recipient))
			ok = false
		}
		return bz
	})
	return ok
}
