package handlers

import (
	"fmt"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/errcode"
	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func wrapWithdrawTransaction(tx *types.WithdrawTransaction) *types.Transaction {
	return &types.Transaction{
		Tx: &types.Transaction_WithdrawTx{
			WithdrawTx: tx,
		},
	}
}

func TestCheckAndDeliverWithdraw(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a Withdraw Transaction", t, func() {

		Convey("If it is a valid transaction", func() {
			appCtx.Reset()
			rawTx := wrapWithdrawTransaction(&types.WithdrawTransaction{
				// TODO
			})

			Convey("CheckTx should return Code 0", func() {
				res := checkWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, 0)
			})

			Convey("DeliverTx should return Code 0", func() {
				res := deliverWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, 0)

				Convey("Should be able to query the transaction info afterwards", func() {
					_ = res.Data // TODO: ID
					// TODO: query
				})
			})
		})

		Convey("If it is an invalid address format", func() {
			appCtx.Reset()

			rawTx := wrapWithdrawTransaction(&types.WithdrawTransaction{
				// TODO
			})

			code, _ := errcode.WithdrawCheckTxInvalidFormat()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.WithdrawDeliverTxInvalidFormat()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is an invalid signature version", func() {
			appCtx.Reset()

			rawTx := wrapWithdrawTransaction(&types.WithdrawTransaction{
				// TODO
			})

			code, _ := errcode.WithdrawCheckTxInvalidSignature()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.WithdrawDeliverTxInvalidSignature()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is an invalid signature format", func() {
			appCtx.Reset()

			rawTx := wrapWithdrawTransaction(&types.WithdrawTransaction{
				// TODO
			})

			code, _ := errcode.WithdrawCheckTxInvalidSignature()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.WithdrawDeliverTxInvalidSignature()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is a replayed transaction", func() {
			appCtx.Reset()

			rawTx := wrapWithdrawTransaction(&types.WithdrawTransaction{
				// TODO
			})

			code, _ := errcode.WithdrawCheckTxDuplicated()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.WithdrawDeliverTxDuplicated()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverWithdraw(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})
	})
}

func TestvalidateWithdrawTransactionFormat(t *testing.T) {
	Convey("Given a Withdraw transaction in valid format", t, func() {
		tx := &types.WithdrawTransaction{} // TODO

		Convey("The transaction should pass the validation", func() {
			So(validateWithdrawTransactionFormat(tx), ShouldBeTrue)
		})
	})

	Convey("Given a Withdraw transaction in invalid format", t, func() {
		tx := &types.WithdrawTransaction{} // TODO

		Convey("The transaction should not pass the validation", func() {
			So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
		})
	})

	Convey("Given a Withdraw transaction with invalid nouce", t, func() {
		tx := &types.WithdrawTransaction{} // TODO

		Convey("The transaction should not pass the validation", func() {
			So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
		})
	})
}

func TestWithdraw(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a valid Withdraw transaction", t, func() {
		tx := &types.WithdrawTransaction{} // TODO

		Convey("The transaction should be pass", func() {
			withdraw(state, tx)
			// TODO: checking
		})

		Convey("But the same Withdraw transaction cannot be replayed", func() {
			withdraw(state, tx)
			// TODO: checking
		})
	})

	Convey("Given an invalid Withdraw transaction", t, func() {
		appCtx.Reset()
		tx := &types.WithdrawTransaction{} // TODO

		Convey("The transaction should not be pass if sender not exist ", func() {
			withdraw(state, tx)
			// TODO: checking
		})

		tx = &types.WithdrawTransaction{} // TODO

		Convey("The transaction should not be pass if receiver not exist", func() {
			withdraw(state, tx)
			// TODO: checking
		})

		tx = &types.WithdrawTransaction{} // TODO

		Convey("The transaction should not be pass if there is not enough balance", func() {
			withdraw(state, tx)
			// TODO: checking
		})
	})
}
