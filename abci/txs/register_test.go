package txs

import (
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterValidateFormat(t *testing.T) {
	Convey("For a register transaction", t, func() {
		Convey("If the transaction has valid format", func() {
			regTx := RegisterTx("0x0000000000000000000000000000000000000000", "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
			Convey("The validation should succeed", func() {
				So(regTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has nil signature", func() {
			regTx := &RegisterTransaction{
				Addr: *types.Addr("0x0000000000000000000000000000000000000000"),
				Sig:  nil,
			}
			Convey("The validation should fail", func() {
				So(regTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestRegisterSignature(t *testing.T) {
	Convey("If a register transaction is valid", t, func() {
		regTx := RegisterTx("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", "65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141b")
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := regTx.Sig.RecoverAddress(regTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the address of the register transaction", func() {
				So(recoveredAddr, ShouldResemble, &regTx.Addr)
			})
		})
	})
}

func TestCheckAndDeliverRegister(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		Convey("If a register transaction is valid", func() {
			addrHex := "0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"
			sigHex := "65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141b"
			regTx := RegisterTx(addrHex, sigHex)
			Convey("CheckTx should succeed", func() {
				r := regTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(regTx))
					r := regTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					Convey("CheckTx again should return Duplicated", func() {
						r := regTx.CheckTx(state)
						So(r.Code, ShouldEqual, response.RegisterDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := regTx.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.RegisterDuplicated.Code)
						})
					})
				})
			})
			Convey("RawRegisterTx should return the same encoded raw tx", func() {
				rawTx := RawRegisterTx(addrHex, sigHex)
				So(rawTx, ShouldResemble, EncodeTx(regTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *RegisterTransaction
					err := types.AminoCodec().UnmarshalBinary(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, regTx)
				})
			})
		})
		Convey("If a register transaction has invalid format", func() {
			regTx := &RegisterTransaction{
				Addr: *types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"),
				Sig:  nil,
			}
			Convey("CheckTx should return InvalidFormat", func() {
				r := regTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.RegisterInvalidFormat.Code)
				Convey("DeliverTx should return InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(regTx))
					r := regTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.RegisterInvalidFormat.Code)
				})
			})
		})
		Convey("If a register transaction has invalid signature", func() {
			regTx := RegisterTx("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", "65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141c")
			Convey("CheckTx should return InvalidSignature", func() {
				r := regTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.RegisterInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(regTx))
					r := regTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.RegisterInvalidSignature.Code)
				})
			})
		})
		Convey("If a register transaction has signature from others", func() {
			regTx := RegisterTx("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", "b287bb3c420155326e0a7fe3a66fed6c397a4bdb5ddcd54960daa0f06c1fbf06300e862dbd3ae3daeae645630e66962b81cf6aa9ffb258aafde496e0310ab8551c")
			Convey("CheckTx should return InvalidSignature", func() {
				r := regTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.RegisterInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(regTx))
					r := regTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.RegisterInvalidSignature.Code)
				})
			})
		})
	})
}
