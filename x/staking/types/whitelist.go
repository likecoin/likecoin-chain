package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValidatorsWhitelist struct {
	Validators []sdk.ValAddress `json:"validators"`
}

func (whitelist ValidatorsWhitelist) String() string {
	bz, err := json.Marshal(whitelist)
	if err != nil {
		panic(err)
	}
	return string(bz)
}
