package withdraw

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/db"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProof(t *testing.T) {
	Convey("Given an IAVL tree", t, func() {
		tree := iavl.NewMutableTree(db.NewMemDB(), 0)
		Convey("Given an IAVL RangeProof", func() {
			key := [32]byte{}
			value := []byte{1}
			for i := 0; i < 100000; i++ {
				// treat key as a 32-byte little-endian integer and increment it
				for j := 0; j < 32; j++ {
					if j != 255 {
						for k := 0; k < j-1; k++ {
							key[k] = 0
						}
						key[j]++
					}
				}
				tree.Set(key[:], value)
			}
			tree.SaveVersion()
			_, proof, err := tree.GetWithProof(key[:])
			So(err, ShouldBeNil)
			Convey("The JSON representation should be parsed into a RangeProof in withdraw service", func() {
				jsonBytes, err := json.Marshal(proof)
				So(err, ShouldBeNil)
				myProof := ParseRangeProof(jsonBytes)
				So(myProof, ShouldNotBeNil)
				So(myProof.ComputeRootHash(), ShouldResemble, proof.ComputeRootHash())
				Convey("The contract proof should be valid", func() {
					contractProof := myProof.ContractProof()
					i := 0
					So(contractProof[i], ShouldEqual, 0)
					i++
					versionLen := int(contractProof[i])
					i++
					leafVersionBytes := contractProof[i : i+versionLen]
					leafVersion, n, err := amino.DecodeVarint(leafVersionBytes)
					So(err, ShouldBeNil)
					So(n, ShouldEqual, versionLen)
					So(leafVersion, ShouldEqual, myProof.Leaves[0].Version)
					i += versionLen

					// Below mimics the logic in Relay contract
					buf := new(bytes.Buffer)
					buf.Write([]byte{0, 2})
					buf.Write(leafVersionBytes)
					buf.WriteByte(32)
					buf.Write(key[:])
					buf.WriteByte(32)
					buf.Write(tmhash.Sum(value))
					hash := tmhash.Sum(buf.Bytes())

					pathLength := int(contractProof[i])
					So(pathLength, ShouldEqual, len(myProof.LeftPath))
					i++
					for j := 0; j < pathLength; j++ {
						prefixLengthAndOrder := int(contractProof[i])
						order := prefixLengthAndOrder & 0x80
						prefixLen := prefixLengthAndOrder & 0x7F
						i++
						buf = new(bytes.Buffer)
						buf.Write(contractProof[i : i+prefixLen])
						i += prefixLen
						if order == 0 {
							buf.WriteByte(32)
							buf.Write(hash)
							buf.WriteByte(32)
							buf.Write(contractProof[i : i+32])
						} else {
							buf.WriteByte(32)
							buf.Write(contractProof[i : i+32])
							buf.WriteByte(32)
							buf.Write(hash)
						}
						i += 32
						hash = tmhash.Sum(buf.Bytes())
					}
					So(hash, ShouldResemble, myProof.ComputeRootHash())
				})
			})
		})
	})
}
