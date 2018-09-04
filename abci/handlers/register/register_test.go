package register

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	. "github.com/smartystreets/goconvey/convey"
)

var sigHex = types.NewSignatureFromHex("0xb19ced763ac63a33476511ecce1df4ebd91bb9ae8b2c0d24b0a326d96c5717122ae0c9b5beacaf4560f3a2535a7673a3e567ff77f153e452907169d431c951091b")

func TestCheckAndDeliverRegister(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given an empty Transaction, CheckTx and DeliverTx should panic", t, func() {
		rawTx := &types.Transaction{}

		So(func() { checkRegister(state, rawTx) }, ShouldPanic)
		So(func() { deliverRegister(state, rawTx, nil) }, ShouldPanic)
	})

	Convey("Given a Register Transaction", t, func() {
		appCtx.Reset()

		rawTx := types.RegisterTransaction{
			Addr: fixture.Alice.RawAddress,
			Sig:  sigHex,
		}.ToTransaction()

		rawTxBytes, _ := proto.Marshal(rawTx)
		txHash := utils.HashRawTx(rawTxBytes)

		Convey("If it is a valid transaction", func() {
			Convey("CheckTx should return Code 0", func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, 0)
			})

			Convey("For DeliverTx", func() {
				res := deliverRegister(state, rawTx, txHash)

				Convey("It should return Code 0 and non-empty Data", func() {
					So(res.Code == 0 && len(res.Data) > 0, ShouldBeTrue)
				})

				state.Save()

				Convey("If replay transaction", func() {
					Convey("For CheckTx", func() {
						res := checkRegister(state, rawTx)

						code := response.RegisterCheckTxDuplicated.Code
						Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
							So(res.Code, ShouldEqual, code)
						})
					})

					Convey("For DeliverTx", func() {
						res := deliverRegister(state, rawTx, txHash)

						code := response.RegisterDeliverTxDuplicated.Code
						Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
							So(res.Code, ShouldEqual, code)
						})
					})
				})
			})
		})

		Convey("If its address format is invalid", func() {
			rawTx.GetRegisterTx().Addr = &types.Address{Content: []byte{}}

			code := response.RegisterCheckTxInvalidFormat.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.RegisterDeliverTxInvalidFormat.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverRegister(state, rawTx, txHash)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If its signature version is invalid", func() {
			rawTx.GetRegisterTx().Sig.Version = 2

			code := response.RegisterCheckTxInvalidFormat.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.RegisterDeliverTxInvalidFormat.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverRegister(state, rawTx, txHash)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If its signature is invalid", func() {
			rawTx.GetRegisterTx().Sig = types.NewZeroSignature()

			code := response.RegisterCheckTxInvalidSignature.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.RegisterDeliverTxInvalidSignature.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverRegister(state, rawTx, txHash)

				So(res.Code, ShouldEqual, code)
			})
		})
	})
}

func TestValidateRegisterSignature(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a Register transaction", t, func() {
		appCtx.Reset()
		tx := &types.RegisterTransaction{
			Addr: fixture.Alice.RawAddress,
			Sig:  sigHex,
		}

		Convey("If its signature is valid", func() {
			Convey("It should pass the validation", func() {
				So(validateRegisterSignature(state, tx), ShouldBeTrue)
			})
		})

		Convey("If its signature is invalid", func() {
			tx.Sig = types.NewZeroSignature()

			Convey("It should fail the validation", func() {
				So(validateRegisterSignature(state, tx), ShouldBeFalse)
			})
		})

		Convey("If its signing address is not match", func() {
			tx.Sig = types.NewSignatureFromHex("0x87671a38598a89a850c0883f7e6882fad980d1cd82498dfe6e5cdf512f3988fe599b9fd1f1b5ec6da32d721f466fa956c2248a58ee801dc144ac4b80cc8b35941b")

			Convey("It should fail the validation", func() {
				So(validateRegisterSignature(state, tx), ShouldBeFalse)
			})
		})
	})
}

func TestValidateRegisterTransactionFormat(t *testing.T) {
	Convey("Given a Register transaction", t, func() {
		tx := &types.RegisterTransaction{
			Addr: types.NewZeroAddress(),
			Sig:  types.NewZeroSignature(),
		}

		Convey("If its format is valid", func() {
			Convey("It should pass the validation", func() {
				So(validateRegisterTransactionFormat(tx), ShouldBeTrue)
			})
		})

		Convey("If its address format is invalid", func() {
			tx.Addr = types.NewAddressFromHex("")

			Convey("It should fail the validation", func() {
				So(validateRegisterTransactionFormat(tx), ShouldBeFalse)
			})
		})

		Convey("If its signature version is invalid", func() {
			tx.Sig.Version = 2

			Convey("It should fail the validation", func() {
				So(validateRegisterTransactionFormat(tx), ShouldBeFalse)
			})
		})

		Convey("With invalid signature format", func() {
			tx.Sig = types.NewSignatureFromHex("")

			Convey("It should fail the validation", func() {
				So(validateRegisterTransactionFormat(tx), ShouldBeFalse)
			})
		})
	})
}
