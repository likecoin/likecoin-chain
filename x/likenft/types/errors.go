package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/likenft module sentinel errors
var (
	ErrInvalidIscnId                     = sdkerrors.Register(ModuleName, 1, "invalid ISCN ID")
	ErrIscnRecordNotFound                = sdkerrors.Register(ModuleName, 2, "ISCN record not found")
	ErrFailedToSaveClass                 = sdkerrors.Register(ModuleName, 3, "Failed to save class")
	ErrFailedToMarshalData               = sdkerrors.Register(ModuleName, 4, "Failed to marshal data")
	ErrNftClassNotFound                  = sdkerrors.Register(ModuleName, 5, "NFT Class not found")
	ErrFailedToUnmarshalData             = sdkerrors.Register(ModuleName, 6, "Failed to unmarshal data")
	ErrNftClassNotRelatedToAnyIscn       = sdkerrors.Register(ModuleName, 7, "NFT Class not related to any ISCN")
	ErrFailedToQueryIscnRecord           = sdkerrors.Register(ModuleName, 8, "Failed to query iscn record")
	ErrCannotUpdateClassWithMintedTokens = sdkerrors.Register(ModuleName, 9, "Cannot update class with minted tokens")
	ErrFailedToUpdateClass               = sdkerrors.Register(ModuleName, 10, "Failed to update class")
	ErrFailedToMintNFT                   = sdkerrors.Register(ModuleName, 11, "Failed to mint NFT")
	ErrInvalidTokenId                    = sdkerrors.Register(ModuleName, 12, "Invalid Token ID")
	ErrNftNotFound                       = sdkerrors.Register(ModuleName, 13, "NFT not found")
	ErrNftNotBurnable                    = sdkerrors.Register(ModuleName, 14, "NFT not burnable")
	ErrFailedToBurnNFT                   = sdkerrors.Register(ModuleName, 15, "Failed to burn NFT")
	ErrNftClassNotRelatedToAnyAccount    = sdkerrors.Register(ModuleName, 16, "NFT Class not related to any account")
	ErrNftNoSupply                       = sdkerrors.Register(ModuleName, 17, "No supply left for the NFT Class")
	ErrInsufficientFunds                 = sdkerrors.Register(ModuleName, 18, "Insufficient funds")
	ErrClaimableNftAlreadyExists         = sdkerrors.Register(ModuleName, 19, "Claimable NFT already exists")
	ErrClaimableNftNotFound              = sdkerrors.Register(ModuleName, 20, "Claimable NFT not found")
)
