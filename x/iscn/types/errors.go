package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidIscnRecord  = sdkerrors.Register(ModuleName, 1, "invalid ISCN record")
	ErrInvalidIscnId      = sdkerrors.Register(ModuleName, 2, "invalid ISCN ID")
	ErrReusingIscnId      = sdkerrors.Register(ModuleName, 3, "reusing ISCN ID")
	ErrCidAlreadyExist    = sdkerrors.Register(ModuleName, 4, "CID already exist")
	ErrEncodingJsonLd     = sdkerrors.Register(ModuleName, 5, "error when encoding JSON-LD record")
	ErrAddingIscnRecord   = sdkerrors.Register(ModuleName, 6, "error when adding ISCN record")
	ErrInvalidIscnVersion = sdkerrors.Register(ModuleName, 7, "invalid ISCN ID version")
	ErrDeductIscnFee      = sdkerrors.Register(ModuleName, 8, "error when deducting fee for ISCN record")
	ErrCidNotFound        = sdkerrors.Register(ModuleName, 9, "CID not found")
	ErrRecordNotFound     = sdkerrors.Register(ModuleName, 10, "record not found")
)
