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
				// claimed relation is valid
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
				// claimed relation is valid
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

func (k Keeper) resolveClassParentAndOwner(ctx sdk.Context, parentInput types.ClassParentInput, ownerBech32 string) (types.ClassParentAndOwner, error) {
	if parentInput.Type == types.ClassParentType_ISCN {
		iscnId, iscnRecord, err := k.resolveIscnIdAndRecord(ctx, parentInput.IscnIdPrefix)
		if err != nil {
			return types.ClassParentAndOwner{}, err
		}
		return types.ClassParentAndOwner{
			ClassParent: types.ClassParent{
				Type:              types.ClassParentType_ISCN,
				IscnIdPrefix:      iscnId.Prefix.String(),
				IscnVersionAtMint: iscnRecord.LatestVersion,
			},
			Owner: iscnRecord.OwnerAddress(),
		}, nil
	} else if parentInput.Type == types.ClassParentType_ACCOUNT {
		owner, err := sdk.AccAddressFromBech32(ownerBech32)
		if err != nil {
			return types.ClassParentAndOwner{}, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
		}
		return types.ClassParentAndOwner{
			ClassParent: types.ClassParent{
				Type:    types.ClassParentType_ACCOUNT,
				Account: owner.String(),
			},
			Owner: owner,
		}, nil
	} else {
		return types.ClassParentAndOwner{}, sdkerrors.ErrInvalidRequest.Wrapf("Unsupported parent type %s in nft class", parentInput.Type.String())
	}
}

func (k msgServer) validateAndGetClassParentAndOwner(ctx sdk.Context, classId string, classData *types.ClassData) (*types.ClassParentAndOwner, error) {
	if err := k.validateClassParentRelation(ctx, classId, classData.Parent); err != nil {
		return nil, err
	}

	// refresh parent info (e.g. iscn latest version) & check ownership
	parent, err := k.resolveClassParentAndOwner(ctx, classData.Parent.ToInput(), classData.Parent.Account)
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func (k msgServer) validateClaimPeriods(classConfig *types.ClassConfig) error {
	for _, claimPeriod := range classConfig.ClaimPeriods {
		// Ensure all claim period start time is before reveal time
		if claimPeriod.StartTime.After(*classConfig.RevealTime) {
			return types.ErrInvalidNftClassConfig.Wrapf("One of the claim periods' start time %s is after reveal time %s", claimPeriod.StartTime.String(), classConfig.RevealTime.String())
		}
		// Ensure all the addresses in allow list is valid
		for _, allowedAddress := range claimPeriod.AllowedAddressList {
			if _, err := sdk.AccAddressFromBech32(allowedAddress); err != nil {
				return sdkerrors.ErrInvalidAddress.Wrapf("One of the allowed addresses %s is invalid", allowedAddress)
			}
		}
	}
	return nil
}

func (k msgServer) validateClassConfig(classConfig *types.ClassConfig) error {
	// Ensure claim periods and reveal time are set when blind box mode is enabled
	if classConfig.EnableBlindBox {
		if classConfig.ClaimPeriods == nil {
			return types.ErrInvalidNftClassConfig.Wrapf("Claim periods are enabled but no claim periods are provided")
		}

		if classConfig.RevealTime == nil {
			return types.ErrInvalidNftClassConfig.Wrapf("Claim periods are enabled but no reveal time is provided")
		}

		if err := k.validateClaimPeriods(classConfig); err != nil {
			return err
		}
	}
	return nil
}
