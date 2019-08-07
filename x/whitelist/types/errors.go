package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName
)

func ErrInvalidApprover(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, sdk.CodeInvalidAddress, "approver address is invalid")
}

func ErrValidatorNotInWEhitelist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, staking.CodeInvalidValidator, "validator not in whitelist")
}
