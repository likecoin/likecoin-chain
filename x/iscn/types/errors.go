package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidIscnRecord  = sdkerrors.Register(ModuleName, 1, "invalid ISCN record")
	ErrInvalidIscnId      = sdkerrors.Register(ModuleName, 2, "invalid ISCN ID")
	ErrReusingIscnId      = sdkerrors.Register(ModuleName, 3, "reusing ISCN ID")
	ErrRecordAlreadyExist = sdkerrors.Register(ModuleName, 4, "record already exist")
	ErrEncodingJsonLd     = sdkerrors.Register(ModuleName, 5, "error when encoding JSON-LD record")
	ErrInvalidIscnVersion = sdkerrors.Register(ModuleName, 6, "invalid ISCN ID version")
	ErrDeductIscnFee      = sdkerrors.Register(ModuleName, 7, "error when deducting fee for ISCN record")
	ErrRecordNotFound     = sdkerrors.Register(ModuleName, 8, "record not found")
)
