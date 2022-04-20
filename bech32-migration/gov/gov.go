package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/likecoin/likechain/bech32-migration/utils"
)

func MigrateAddressBech32(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) {
	ctx.Logger().Info("Migration of address bech32 for gov module begin")
	voteCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.VotesKeyPrefix, func(bz []byte) []byte {
		vote := types.Vote{}
		cdc.MustUnmarshal(bz, &vote)
		vote.Voter = utils.ConvertAccAddr(vote.Voter)
		voteCount++
		return cdc.MustMarshal(&vote)
	})
	depositCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.DepositsKeyPrefix, func(bz []byte) []byte {
		deposit := types.Deposit{}
		cdc.MustUnmarshal(bz, &deposit)
		deposit.Depositor = utils.ConvertAccAddr(deposit.Depositor)
		depositCount++
		return cdc.MustMarshal(&deposit)
	})
	communityPoolSpendProposalCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.ProposalsKeyPrefix, func(bz []byte) []byte {
		proposal := types.Proposal{}
		cdc.MustUnmarshal(bz, &proposal)
		content := proposal.GetContent()
		communityPoolSpendProposal, ok := content.(*distrtypes.CommunityPoolSpendProposal)
		if !ok {
			return bz
		}
		communityPoolSpendProposal.Recipient = utils.ConvertAccAddr(communityPoolSpendProposal.Recipient)
		newContentAny, err := codectypes.NewAnyWithValue(communityPoolSpendProposal)
		if err != nil {
			panic(err)
		}
		proposal.Content = newContentAny
		communityPoolSpendProposalCount++
		return cdc.MustMarshal(&proposal)
	})
	ctx.Logger().Info(
		"Migration of address bech32 for gov module done",
		"vote_count", voteCount,
		"deposit_count", depositCount,
		"community_pool_spend_proposal_count", communityPoolSpendProposalCount,
	)
}
