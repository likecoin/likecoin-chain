package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName
)

func ErrInvalidSender(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, sdk.CodeInvalidAddress, "sender address is invalid")
}
