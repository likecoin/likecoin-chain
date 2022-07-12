package testutil

import (
	"github.com/likecoin/likecoin-chain/v3/backport/cosmos-sdk/v0.46.0-rc1/x/nft"
	"github.com/likecoin/likecoin-chain/v3/x/likenft/types"
)

func MakeDummyNFTClassesForISCN(msg types.ClassesByISCN) []nft.Class {
	classes := make([]nft.Class, len(msg.ClassIds))
	for i, classId := range msg.ClassIds {
		classes[i] = nft.Class{
			Id: classId,
		}
	}
	return classes
}

func BatchMakeDummyNFTClassesForISCN(msgs []types.ClassesByISCN) [][]nft.Class {
	output := make([][]nft.Class, len(msgs))
	for i, msg := range msgs {
		output[i] = MakeDummyNFTClassesForISCN(msg)
	}
	return output
}

func MakeDummyNFTClassesForAccount(msg types.ClassesByAccount) []nft.Class {
	classes := make([]nft.Class, len(msg.ClassIds))
	for i, classId := range msg.ClassIds {
		classes[i] = nft.Class{
			Id: classId,
		}
	}
	return classes
}

func BatchMakeDummyNFTClassesForAccount(msgs []types.ClassesByAccount) [][]nft.Class {
	output := make([][]nft.Class, len(msgs))
	for i, msg := range msgs {
		output[i] = MakeDummyNFTClassesForAccount(msg)
	}
	return output
}
