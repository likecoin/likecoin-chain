package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Whitelist []sdk.ValAddress

func (whitelist Whitelist) String() string {
	bz, err := json.Marshal(whitelist)
	if err != nil {
		panic(err)
	}
	return string(bz)
}