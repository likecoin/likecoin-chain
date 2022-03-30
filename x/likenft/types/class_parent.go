package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ClassParentAndOwner struct {
	ClassParent
	Owner sdk.AccAddress
}

func (m ClassParent) ToInput() ClassParentInput {
	return ClassParentInput{
		Type:         m.Type,
		IscnIdPrefix: m.IscnIdPrefix,
	}
}
