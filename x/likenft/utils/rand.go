package utils

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RandSeedFromLastBlock(ctx sdk.Context) int64 {
	appHash := ctx.BlockHeader().AppHash
	seed, read := binary.Varint(appHash[:8]) // only use first 64 bit / 8 bytes
	if seed == 0 || read <= 0 {
		panic(fmt.Errorf("Failed seeding random due to bad last block hash length, read %d bytes", read))
	}
	return seed
}
