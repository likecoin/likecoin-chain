package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *ClassesByAccount) ToStoreRecord() ClassesByAccountStoreRecord {
	acc, err := sdk.AccAddressFromBech32(m.Account)
	if err != nil {
		panic(err)
	}

	return ClassesByAccountStoreRecord{
		AccAddress: acc,
		ClassIds:   m.ClassIds,
	}
}

func (m *ClassesByAccountStoreRecord) ToPublicRecord() ClassesByAccount {
	return ClassesByAccount{
		Account:  sdk.AccAddress(m.AccAddress).String(),
		ClassIds: m.ClassIds,
	}
}
