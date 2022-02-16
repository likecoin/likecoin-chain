package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/likenft module sentinel errors
var (
	ErrInvalidIscnId               = sdkerrors.Register(ModuleName, 1, "invalid ISCN ID")
	ErrIscnRecordNotFound          = sdkerrors.Register(ModuleName, 2, "ISCN record not found")
	ErrFailedToSaveClass           = sdkerrors.Register(ModuleName, 3, "Failed to save class")
	ErrFailedToMarshalData         = sdkerrors.Register(ModuleName, 4, "Failed to marshal data")
	ErrNftClassNotFound            = sdkerrors.Register(ModuleName, 5, "NFT Class not found")
	ErrFailedToUnmarshalData       = sdkerrors.Register(ModuleName, 6, "Failed to unmarshal data")
	ErrNftClassNotRelatedToAnyIscn = sdkerrors.Register(ModuleName, 7, "NFT Class not related to any ISCN")
	ErrFailedToQueryIscnRecord     = sdkerrors.Register(ModuleName, 8, "Failed to query iscn record")
)
