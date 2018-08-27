package register

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

var addr = &types.Address{Content: common.FromHex("0x064b663abf9d74277a07aa7563a8a64a54de8c0a")}
var sig = common.FromHex("0xb19ced763ac63a33476511ecce1df4ebd91bb9ae8b2c0d24b0a326d96c5717122ae0c9b5beacaf4560f3a2535a7673a3e567ff77f153e452907169d431c951091b")

func wrapRegisterTransaction(tx *types.RegisterTransaction) *types.Transaction {
	return &types.Transaction{
		Tx: &types.Transaction_RegisterTx{
			RegisterTx: tx,
		},
	}
}

func TestCheckAndDeliverRegister(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a Register Transaction", t, func() {
		Convey("If it is a valid Register transaction", func() {
			rawTx := wrapRegisterTransaction(&types.RegisterTransaction{
				Addr: addr,
				Sig: &types.Signature{
					Content: sig,
					Version: 1,
				},
			})

			Convey("CheckTx should return Code 0", func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, 0)
			})

			Convey("DeliverTx should return Code 0 and non-empty Data", func() {
				res := deliverRegister(state, rawTx)

				So(res.Code == 0 && len(res.Data) > 0, ShouldBeTrue)
			})

			state.Save()

			Convey("If it is given the same transaction", func() {
				code := response.RegisterCheckTxDuplicated.Code
				Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
					res := checkRegister(state, rawTx)

					So(res.Code, ShouldEqual, code)
				})

				code = response.RegisterDeliverTxDuplicated.Code
				Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
					res := deliverRegister(state, rawTx)

					So(res.Code, ShouldEqual, code)
				})
			})
		})

		Convey("If it is a Register transaction with invalid address format", func() {
			appCtx.Reset()

			rawTx := wrapRegisterTransaction(&types.RegisterTransaction{
				Addr: &types.Address{Content: []byte{}},
				Sig: &types.Signature{
					Content: sig,
					Version: 1,
				},
			})

			code := response.RegisterCheckTxInvalidFormat.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.RegisterDeliverTxInvalidFormat.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is a Register transaction with invalid signature version", func() {
			appCtx.Reset()

			rawTx := wrapRegisterTransaction(&types.RegisterTransaction{
				Addr: addr,
				Sig: &types.Signature{
					Content: sig,
					Version: 2,
				},
			})

			code := response.RegisterCheckTxInvalidFormat.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.RegisterDeliverTxInvalidFormat.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is a Register transaction with invalid signature format", func() {
			appCtx.Reset()

			rawTx := wrapRegisterTransaction(&types.RegisterTransaction{
				Addr: addr,
				Sig: &types.Signature{
					Content: common.FromHex("0xd880732022a41a404669ded27f41564df20e728280264860a968a2d3ae0e745f6a576539b36ac4a27e4e9bde1e74cdf58144dd130dc6d6328ab6440129c344f51c"),
					Version: 1,
				},
			})

			code := response.RegisterCheckTxInvalidSignature.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.RegisterDeliverTxInvalidSignature.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverRegister(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})
	})
}

func TestValidateRegisterSignature(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a Register transaction with valid signature", t, func() {
		tx := &types.RegisterTransaction{
			Addr: addr,
			Sig: &types.Signature{
				Content: sig,
				Version: 1,
			},
		}

		Convey("The signature should pass the validation", func() {
			So(validateRegisterSignature(state, tx), ShouldBeTrue)
		})
	})

	Convey("Given a Register transaction with invalid signature", t, func() {
		tx := &types.RegisterTransaction{
			Addr: addr,
			Sig: &types.Signature{
				Content: common.FromHex(""),
				Version: 1,
			},
		}

		Convey("The signature should not pass the validation", func() {
			So(validateRegisterSignature(state, tx), ShouldBeFalse)
		})
	})
}

func TestValidateRegisterTransaction(t *testing.T) {
	Convey("Given a Register transaction", t, func() {
		Convey("With valid format", func() {
			tx := &types.RegisterTransaction{
				Addr: addr,
				Sig: &types.Signature{
					Version: 1,
					Content: sig,
				},
			}

			Convey("It should pass the validation", func() {
				So(validateRegisterTransaction(tx), ShouldBeTrue)
			})
		})

		Convey("With invalid address format", func() {
			tx := &types.RegisterTransaction{
				Addr: &types.Address{Content: common.FromHex("")},
				Sig: &types.Signature{
					Version: 1,
					Content: sig,
				},
			}
			Convey("It should fail the validation", func() {
				So(validateRegisterTransaction(tx), ShouldBeFalse)
			})
		})

		Convey("With invalid signature version", func() {
			tx := &types.RegisterTransaction{
				Addr: addr,
				Sig: &types.Signature{
					Version: 0,
					Content: sig,
				},
			}
			Convey("It should fail the validation", func() {
				So(validateRegisterTransaction(tx), ShouldBeFalse)
			})
		})

		Convey("With invalid signature format", func() {
			tx := &types.RegisterTransaction{
				Addr: addr,
				Sig: &types.Signature{
					Version: 1,
					Content: common.FromHex(""),
				},
			}
			Convey("It should fail the validation", func() {
				So(validateRegisterTransaction(tx), ShouldBeFalse)
			})
		})
	})
}
