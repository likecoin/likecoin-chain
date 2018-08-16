package handlers

import (
	"fmt"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/errcode"
	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func wrapDepositTransaction(tx *types.DepositTransaction) *types.Transaction {
	return &types.Transaction{
		Tx: &types.Transaction_DepositTx{
			DepositTx: tx,
		},
	}
}

func TestCheckAndDeliverDeposit(t *testing.T) {
	ctx := context.NewMock()

	Convey("Given a Deposit Transaction", t, func() {

		Convey("If it is a valid transaction", func() {
			ctx.Reset()
			rawTx := wrapDepositTransaction(&types.DepositTransaction{
				// TODO
			})

			Convey("CheckTx should return Code 0", func() {
				res := checkDeposit(ctx, rawTx)

				So(res.Code, ShouldEqual, 0)
			})

			Convey("DeliverTx should return Code 0", func() {
				res := deliverDeposit(ctx, rawTx)

				So(res.Code, ShouldEqual, 0)

				Convey("Should be able to query the transaction info afterwards", func() {
					_ = res.Data // TODO: ID
					// TODO: query
				})
			})
		})

		Convey("If it is an invalid address format", func() {
			ctx.Reset()

			rawTx := wrapDepositTransaction(&types.DepositTransaction{
				// TODO
			})

			code, _ := errcode.DepositCheckTxInvalidFormat()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkDeposit(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.DepositDeliverTxInvalidFormat()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverDeposit(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is a replayed transaction", func() {
			ctx.Reset()

			rawTx := wrapDepositTransaction(&types.DepositTransaction{
				// TODO
			})

			code, _ := errcode.DepositCheckTxDuplicated()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkDeposit(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.DepositDeliverTxDuplicated()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverDeposit(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})
	})
}

func TestValidateDepositTransactionFormat(t *testing.T) {
	Convey("Given a Deposit transaction in valid format", t, func() {
		tx := &types.DepositTransaction{} // TODO

		Convey("The transaction should pass the validation", func() {
			So(validateDepositTransactionFormat(tx), ShouldBeTrue)
		})
	})

	Convey("Given a Deposit transaction in invalid format", t, func() {
		tx := &types.DepositTransaction{} // TODO

		Convey("The transaction should not pass the validation", func() {
			So(validateDepositTransactionFormat(tx), ShouldBeFalse)
		})
	})
}

func TestDeposit(t *testing.T) {
	ctx := context.NewMock()

	Convey("Given a valid Deposit transaction", t, func() {
		ctx.Reset()
		tx := &types.DepositTransaction{} // TODO

		Convey("The transaction should be pass", func() {
			deposit(ctx, tx)
			// TODO: checking
		})

		Convey("But the same Deposit transaction cannot be replayed", func() {
			deposit(ctx, tx)
			// TODO: checking
		})
	})

	Convey("Given an invalid Deposit transaction", t, func() {
		ctx.Reset()
		tx := &types.DepositTransaction{} // TODO

		Convey("The transaction should not be pass if receiver not exist", func() {
			deposit(ctx, tx)
			// TODO: checking
		})
	})
}
