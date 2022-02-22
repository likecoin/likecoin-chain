package testutil

import (
	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/x/likenft/types"
)

func DummyConcretizeClassesByISCN(msg types.ClassesByISCN) types.ConcreteClassesByISCN {
	var classes []*nft.Class
	for _, classId := range msg.ClassIds {
		class := nft.Class{}
		class.Id = classId
		classes = append(classes, &class)
	}
	return types.ConcreteClassesByISCN{
		IscnIdPrefix: msg.IscnIdPrefix,
		Classes:      classes,
	}
}

func BatchDummyConcretizeClassesByISCN(msgs []types.ClassesByISCN) []types.ConcreteClassesByISCN {
	var output []types.ConcreteClassesByISCN
	for _, msg := range msgs {
		output = append(output, DummyConcretizeClassesByISCN(msg))
	}
	return output
}
