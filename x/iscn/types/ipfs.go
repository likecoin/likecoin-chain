package types

import (
	iscnblock "github.com/likecoin/iscn-ipld/plugin/block"
	"github.com/likecoin/iscn-ipld/plugin/block/content"
	"github.com/likecoin/iscn-ipld/plugin/block/entity"
	"github.com/likecoin/iscn-ipld/plugin/block/kernel"
	"github.com/likecoin/iscn-ipld/plugin/block/right"
	"github.com/likecoin/iscn-ipld/plugin/block/rights"
	"github.com/likecoin/iscn-ipld/plugin/block/stakeholder"
	"github.com/likecoin/iscn-ipld/plugin/block/stakeholders"
	timeperiod "github.com/likecoin/iscn-ipld/plugin/block/time_period"
	"github.com/multiformats/go-multibase"
)

var (
	IscnKernelCodecType   = uint64(iscnblock.CodecISCN)
	IscnContentCodecType  = uint64(iscnblock.CodecContent)
	EntityCodecType       = uint64(iscnblock.CodecEntity)
	RightTermsCodecType   = uint64(0x70) // MerkleDAG protobuf, for CIDv0 / normal file
	RightsCodecType       = uint64(iscnblock.CodecRights)
	StakeholdersCodecType = uint64(iscnblock.CodecStakeholders)
)

var CidMbaseEncoder multibase.Encoder

func init() {
	var err error
	CidMbaseEncoder, err = multibase.EncoderByName("base58btc")
	if err != nil {
		panic(err)
	}
	kernel.Register()
	rights.Register()
	stakeholders.Register()
	content.Register()
	entity.Register()

	right.Register()
	stakeholder.Register()
	timeperiod.Register()
}
