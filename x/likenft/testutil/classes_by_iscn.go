package testutil

import (
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func MakeDummyNFTClasses(msg types.ClassesByISCN) []nft.Class {
	classes := make([]nft.Class, len(msg.ClassIds))
	for i, classId := range msg.ClassIds {
		classes[i] = nft.Class{
			Id: classId,
		}
	}
	return classes
}

func BatchMakeDummyNFTClasses(msgs []types.ClassesByISCN) [][]nft.Class {
	output := make([][]nft.Class, len(msgs))
	for i, msg := range msgs {
		output[i] = MakeDummyNFTClasses(msg)
	}
	return output
}
