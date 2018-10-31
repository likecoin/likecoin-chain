package types

// InitAccountInfo is a structure holding info for each account in the initial account list
type InitAccountInfo struct {
	ID                    LikeChainID `json:"id"`
	Addr                  Address     `json:"addr"`
	Balance               BigInt      `json:"balance"`
	DepositApproverWeight uint32      `json:"depositApproverWeight"`
}

// AppInitState is for genesis initial blockchain data, including initial account list
type AppInitState struct {
	Accounts []InitAccountInfo `json:"accounts"`
}
