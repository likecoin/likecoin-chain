package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (r RoyaltyStakeholderInput) ToStakeholder() RoyaltyStakeholder {
	accAddress, err := sdk.AccAddressFromBech32(r.Account)
	if err != nil {
		panic(err)
	}
	return RoyaltyStakeholder{
		Account: accAddress,
		Weight:  r.Weight,
	}
}

func (r RoyaltyConfigInput) ToConfig() RoyaltyConfig {
	stakeholders := make([]RoyaltyStakeholder, len(r.Stakeholders))
	for i, stakeholderInput := range r.Stakeholders {
		stakeholders[i] = stakeholderInput.ToStakeholder()
	}
	return RoyaltyConfig{
		RateBasisPoints: r.RateBasisPoints,
		Stakeholders:    stakeholders,
	}
}
