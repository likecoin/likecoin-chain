package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/likecoin/likechain/x/likenft/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ISCNByClass(goCtx context.Context, req *types.QueryISCNByClassRequest) (*types.QueryISCNByClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get iscn id from class data
	class, found := k.nftKeeper.GetClass(ctx, req.ClassId)
	if !found {
		return nil, types.ErrNftClassNotFound.Wrapf("Class id %s not found", req.ClassId)
	}

	// Validate iscn and nft class is related
	var classData types.ClassData
	if err := k.cdc.Unmarshal(class.Data.Value, &classData); err != nil {
		return nil, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}
	classesByISCN, found := k.GetClassesByISCN(ctx, classData.IscnIdPrefix)
	if !found {
		return nil, types.ErrNftClassNotRelatedToAnyIscn.Wrapf("NFT claims it is related to ISCN %s but no mapping is found", classData.IscnIdPrefix)
	}
	isRelated := false
	for _, validClassId := range classesByISCN.ClassIds {
		if validClassId == class.Id {
			// claimed relation is valid
			isRelated = true
			break
		}
	}
	if !isRelated {
		return nil, types.ErrNftClassNotRelatedToAnyIscn.Wrapf("NFT claims it is related to ISCN %s but no mapping is found", classData.IscnIdPrefix)
	}

	// Return related iscn data
	iscnId, err := iscntypes.ParseIscnId(classesByISCN.IscnIdPrefix)
	if err != nil {
		return nil, types.ErrInvalidIscnId.Wrapf(err.Error())
	}
	contentIdRecord := k.iscnKeeper.GetContentIdRecord(ctx, iscnId.Prefix)
	if contentIdRecord == nil {
		return nil, types.ErrFailedToQueryIscnRecord
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
		IscnIdPrefix:  classesByISCN.IscnIdPrefix,
		Owner:         contentIdRecord.OwnerAddress().String(),
		LatestVersion: latestVersion,
		LatestRecord:  record,
	}, nil
}
