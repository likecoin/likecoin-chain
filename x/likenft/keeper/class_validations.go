package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k Keeper) validateClassParentRelation(ctx sdk.Context, classId string, parent types.ClassParent) error {
	if parent.Type == types.ClassParentType_ISCN {
		classesByISCN, found := k.GetClassesByISCN(ctx, parent.IscnIdPrefix)
		if !found {
			return types.ErrNftClassNotRelatedToAnyIscn.Wrapf("NFT claims it is related to ISCN %s but no mapping is found", parent.IscnIdPrefix)
		}
		isRelated := false
		for _, validClassId := range classesByISCN.ClassIds {
			if validClassId == classId {
				// minted relation is valid
				isRelated = true
				break
			}
		}
		if !isRelated {
			return types.ErrNftClassNotRelatedToAnyIscn.Wrapf("NFT claims it is related to ISCN %s but no mapping is found", parent.IscnIdPrefix)
		}
	} else if parent.Type == types.ClassParentType_ACCOUNT {
		acc, err := sdk.AccAddressFromBech32(parent.Account)
		if err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
		}
		classesByAccount, found := k.GetClassesByAccount(ctx, acc)
		if !found {
			return types.ErrNftClassNotRelatedToAnyAccount.Wrapf("NFT claims it is related to account %s but no mapping is found", parent.Account)
		}
		isRelated := false
		for _, validClassId := range classesByAccount.ClassIds {
			if validClassId == classId {
				// minted relation is valid
				isRelated = true
				break
			}
		}
		if !isRelated {
			return types.ErrNftClassNotRelatedToAnyAccount.Wrapf("NFT claims it is related to account %s but no mapping is found", parent.Account)
		}
	} else {
		return sdkerrors.ErrInvalidRequest.Wrapf("Unsupported parent type %s in nft class", parent.Type.String())
	}
	return nil
}

func (k msgServer) sanitizeBlindBoxConfig(blindBoxConfig *types.BlindBoxConfig) (*types.BlindBoxConfig, error) {
	if blindBoxConfig == nil {
		return nil, nil
	}
	if len(blindBoxConfig.MintPeriods) <= 0 {
		return nil, types.ErrInvalidNftClassConfig.Wrapf("Mint period cannot be empty")
	}
	// Sort the mint period by start time
	blindBoxConfig.MintPeriods = SortMintPeriod(blindBoxConfig.MintPeriods, true)
	for _, mintPeriod := range blindBoxConfig.MintPeriods {
		// Ensure all mint period start time is before reveal time
		if mintPeriod.StartTime.After(blindBoxConfig.RevealTime) {
			return nil, types.ErrInvalidNftClassConfig.Wrapf("One of the mint periods' start time %s is after reveal time %s", mintPeriod.StartTime.String(), blindBoxConfig.RevealTime.String())
		}
		// Ensure all the addresses in allow list is valid
		for _, allowedAddress := range mintPeriod.AllowedAddresses {
			if _, err := sdk.AccAddressFromBech32(allowedAddress); err != nil {
				return nil, sdkerrors.ErrInvalidAddress.Wrapf("One of the allowed addresses %s is invalid", allowedAddress)
			}
		}
	}
	return blindBoxConfig, nil
}

func (k msgServer) sanitizeClassConfig(ctx sdk.Context, classConfig types.ClassConfig, mintableCount uint64) (*types.ClassConfig, error) {
	// Ensure mint periods and reveal time are set when blind box mode is enabled
	cleanBlindBoxConfig, err := k.sanitizeBlindBoxConfig(classConfig.BlindBoxConfig)
	if err != nil {
		return nil, err
	}
	classConfig.BlindBoxConfig = cleanBlindBoxConfig

	// Assert new max supply >= mintable count
	if classConfig.IsBlindBox() && classConfig.MaxSupply < mintableCount {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("New max supply %d is less than mintable count %d", classConfig.MaxSupply, mintableCount)
	}

	return &classConfig, nil
}
