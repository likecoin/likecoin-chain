package txs

import (
	"testing"

	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONSignature(t *testing.T) {
	Convey("In the beginning", t, func() {
		Convey("If a JSON signature is valid", func() {
			sigHex := "65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141b"
			sig := JSONSignature(Sig(sigHex))
			Convey("Address recovery should succeed", func() {
				recoveredAddr, err := sig.RecoverAddress(map[string]interface{}{"addr": "0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"})
				So(err, ShouldBeNil)
				Convey("The recovered address should match the signing address", func() {
					So(recoveredAddr, ShouldResemble, types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"))
				})
			})
			Convey("The string representation of the signature should resemble the bytes in the signature", func() {
				s := sig.String()
				So(s, ShouldEqual, "0x"+sigHex)
			})
		})
		Convey("If a JSON structure contains invalid JSON values", func() {
			sigHex := "65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141b"
			sig := JSONSignature(Sig(sigHex))
			Convey("Address recovery should fail", func() {
				_, err := sig.RecoverAddress(map[string]interface{}{"c": 1i})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "unsupported type")
			})
		})
		Convey("If s JSON signature is invalid", func() {
			sigHex := "65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f541400"
			sig := JSONSignature(Sig(sigHex))
			Convey("Address recovery should fail", func() {
				_, err := sig.RecoverAddress(map[string]interface{}{"addr": "0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid signature recovery id")
			})
		})
		Convey("If s signature is not a valid hex string", func() {
			sigHex := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
			Convey("Constructing JSONSignature with Sig should panic", func() {
				So(func() { Sig(sigHex) }, ShouldPanic)
			})
		})
		Convey("If s signature has invalid hex length", func() {
			sigHex := "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
			Convey("Constructing JSONSignature with Sig should panic", func() {
				So(func() { Sig(sigHex) }, ShouldPanic)
			})
		})
	})
}
