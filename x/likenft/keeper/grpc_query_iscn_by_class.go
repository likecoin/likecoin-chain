package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	iscntypes "github.com/likecoin/likecoin-chain/v3/x/iscn/types"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ISCNByClass(goCtx context.Context, req *types.QueryISCNByClassRequest) (*types.QueryISCNByClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get iscn id from class data
	class, classData, err := k.GetClass(ctx, req.ClassId)
	if err != nil {
		return nil, err
	}
	// Check parent is iscn
	if classData.Parent.Type != types.ClassParentType_ISCN {
		return nil, types.ErrNftClassNotRelatedToAnyIscn.Wrapf("NFT Class is related to a %s, not ISCN", classData.Parent.Type.String())
	}
	if err := k.validateClassParentRelation(ctx, class.Id, classData.Parent); err != nil {
		return nil, err
	}

	// Return related iscn data
	iscnId, contentIdRecord, err := k.resolveIscnIdAndRecord(ctx, classData.Parent.IscnIdPrefix)
	if err != nil {
		return nil, err
	}
	latestVersion := contentIdRecord.LatestVersion
	iscnId.Version = latestVersion
	seq := k.iscnKeeper.GetIscnIdSequence(ctx, iscnId)
	storeRecord := k.iscnKeeper.GetStoreRecord(ctx, seq)
	record := iscntypes.QueryResponseRecord{
		Ipld: storeRecord.Cid().String(),
		Data: storeRecord.Data,
	}
	return &types.QueryISCNByClassResponse{
		IscnIdPrefix:  classData.Parent.IscnIdPrefix,
		Owner:         contentIdRecord.OwnerAddress().String(),
		LatestVersion: latestVersion,
		LatestRecord:  record,
	}, nil
}
