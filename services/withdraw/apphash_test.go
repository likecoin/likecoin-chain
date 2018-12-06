package withdraw

import (
	"bytes"
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAppHashProof(t *testing.T) {
	Convey("Given a Tendermint block header", t, func() {
		header := types.Header{
			Version: version.Consensus{
				Block: 123,
				App:   456,
			},
			ChainID:  "some chain ID",
			Height:   789,
			Time:     time.Now(),
			NumTxs:   10,
			TotalTxs: 1337,
			LastBlockID: types.BlockID{
				Hash: cmn.RandBytes(32),
				PartsHeader: types.PartSetHeader{
					Total: 1,
					Hash:  cmn.RandBytes(32),
				},
			},
			LastCommitHash:     cmn.RandBytes(32),
			DataHash:           cmn.RandBytes(32),
			ValidatorsHash:     cmn.RandBytes(32),
			NextValidatorsHash: cmn.RandBytes(32),
			ConsensusHash:      cmn.RandBytes(32),
			AppHash:            cmn.RandBytes(64),
			LastResultsHash:    cmn.RandBytes(32),
			EvidenceHash:       cmn.RandBytes(32),
			ProposerAddress:    cmn.RandBytes(20),
		}
		root, proof := Proof(&header)
		So(root, ShouldResemble, []byte(header.Hash()))
		// Below mimics the logic in Relay contract
		buf := new(bytes.Buffer)
		buf.WriteByte(byte(len(header.AppHash)))
		buf.Write(header.AppHash)
		hash := tmhash.Sum(buf.Bytes())
		buf = new(bytes.Buffer)
		buf.WriteByte(32)
		buf.Write(hash)
		buf.WriteByte(32)
		buf.Write(proof[0])
		hash = tmhash.Sum(buf.Bytes())
		buf = new(bytes.Buffer)
		buf.WriteByte(32)
		buf.Write(hash)
		buf.WriteByte(32)
		buf.Write(proof[1])
		hash = tmhash.Sum(buf.Bytes())
		buf = new(bytes.Buffer)
		buf.WriteByte(32)
		buf.Write(proof[2])
		buf.WriteByte(32)
		buf.Write(hash)
		hash = tmhash.Sum(buf.Bytes())
		buf = new(bytes.Buffer)
		buf.WriteByte(32)
		buf.Write(proof[3])
		buf.WriteByte(32)
		buf.Write(hash)
		hash = tmhash.Sum(buf.Bytes())
		So(hash, ShouldResemble, root)
	})
}
