package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/likenft module sentinel errors
var (
	ErrInvalidIscnId       = sdkerrors.Register(ModuleName, 1, "invalid ISCN ID")
	ErrIscnRecordNotFound  = sdkerrors.Register(ModuleName, 2, "ISCN record not found")
	ErrFailedToSaveClass   = sdkerrors.Register(ModuleName, 3, "Failed to save class")
	ErrFailedToMarshalData = sdkerrors.Register(ModuleName, 4, "Failed to marshal data")
)
