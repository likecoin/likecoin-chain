package types

import (
	gocid "github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

type CID = gocid.Cid

const CidCodecType = 0x0129 // DAG-JSON

func ComputeRecordCid(record []byte) CID {
	mhash, err := multihash.Sum(record, multihash.SHA2_256, -1)
	if err != nil {
		// should never happen
		panic(err)
	}
	return gocid.NewCidV1(CidCodecType, mhash)
}

func MustCidFromBytes(bz []byte) CID {
	_, cid, err := gocid.CidFromBytes(bz)
	if err != nil {
		panic(err)
	}
	return cid
}
