package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-rc1/x/nft"
)

var _ authz.Authorization = &SendNFTAuthorization{}
var _ authz.Authorization = &MintNFTAuthorization{}

func NewMintNFTAuthorization(classId string) *MintNFTAuthorization {
	return &MintNFTAuthorization{
		ClassId: classId,
	}
}

func (a MintNFTAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&MsgMintNFT{})
}

func (a MintNFTAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgMint, ok := msg.(*MsgMintNFT)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgMint.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a MintNFTAuthorization) ValidateBasic() error {
	return nil
}

// implement Authorization interface for SendNFTAuthorization

func NewSendNFTAuthorization(classId string, nftId string) *SendNFTAuthorization {
	return &SendNFTAuthorization{
		ClassId: classId,
		Id:      nftId,
	}
}

func (a SendNFTAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&nft.MsgSend{})
}

func (a SendNFTAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	msgSend, ok := msg.(*nft.MsgSend)
	if !ok {
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidType.Wrap("type mismatch")
	}
	if msgSend.ClassId != a.ClassId {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("class ID mismatch")
	}
	if a.Id != "" && msgSend.Id != a.Id {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrap("NFT ID mismatch")
	}
	return authz.AcceptResponse{Accept: true}, nil
}

func (a SendNFTAuthorization) ValidateBasic() error {
	return nil
}
