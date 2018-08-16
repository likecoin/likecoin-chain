package handlers

import (
	"fmt"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/errcode"
	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func wrapTransferTransaction(tx *types.TransferTransaction) *types.Transaction {
	return &types.Transaction{
		Tx: &types.Transaction_TransferTx{
			TransferTx: tx,
		},
	}
}

func TestCheckAndDeliverTransfer(t *testing.T) {
	ctx := context.NewMock()

	Convey("Given a Transfer Transaction", t, func() {

		Convey("If it is a valid transaction", func() {
			ctx.Reset()
			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			Convey("CheckTx should return Code 0", func() {
				res := checkTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, 0)
			})

			Convey("DeliverTx should return Code 0", func() {
				res := deliverTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, 0)

				Convey("Should be able to query the transaction info afterwards", func() {
					_ = res.Data // TODO: ID
					// TODO: query
				})
			})
		})

		Convey("If it is an invalid address format", func() {
			ctx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code, _ := errcode.TransferCheckTxInvalidFormat()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.TransferDeliverTxInvalidFormat()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is an invalid signature version", func() {
			ctx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code, _ := errcode.TransferCheckTxInvalidSignature()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.TransferDeliverTxInvalidSignature()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is an invalid signature format", func() {
			ctx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code, _ := errcode.TransferCheckTxInvalidSignature()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.TransferDeliverTxInvalidSignature()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is a replayed transaction", func() {
			ctx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code, _ := errcode.TransferCheckTxDuplicated()
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code, _ = errcode.TransferDeliverTxDuplicated()
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(ctx, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})
	})
}

func TestValidateTransferSignature(t *testing.T) {
	Convey("Given a valid Transfer signature", t, func() {
		sig := &types.Signature{} // TODO

		Convey("The signature should pass the validation", func() {
			So(validateTransferSignature(sig), ShouldBeTrue)
		})
	})

	Convey("Given an invalid Transfer signature", t, func() {
		sig := &types.Signature{} // TODO

		Convey("The signature should not pass the validation", func() {
			So(validateTransferSignature(sig), ShouldBeFalse)
		})
	})
}

func TestvalidateTransferTransactionFormat(t *testing.T) {
	Convey("Given a Transfer transaction in valid format", t, func() {
		tx := &types.TransferTransaction{} // TODO

		Convey("The transaction should pass the validation", func() {
			So(validateTransferTransactionFormat(tx), ShouldBeTrue)
		})
	})

	Convey("Given a Transfer transaction in invalid format ", t, func() {
		tx := &types.TransferTransaction{} // TODO

		Convey("The transaction should not pass the validation", func() {
			So(validateTransferTransactionFormat(tx), ShouldBeFalse)
		})
	})

	Convey("Given a Transfer transaction with invalid nouce", t, func() {
		tx := &types.TransferTransaction{} // TODO

		Convey("The transaction should not pass the validation", func() {
			So(validateTransferTransactionFormat(tx), ShouldBeFalse)
		})
	})
}

func TestTransfer(t *testing.T) {
	ctx := context.NewMock()

	Convey("Given a valid Transfer transaction", t, func() {
		ctx.Reset()
		tx := &types.TransferTransaction{} // TODO

		Convey("The transaction should be pass", func() {
			transfer(ctx, tx)
			// TODO check
		})

		Convey("But the same transaction cannot be replayed", func() {
			transfer(ctx, tx)
			// TODO check
		})
	})

	Convey("Given an invalid Transfer transaction", t, func() {
		ctx.Reset()
		tx := &types.TransferTransaction{} // TODO

		Convey("The transaction should not be pass if sender not exist ", func() {
			transfer(ctx, tx)
			// TODO check
		})

		tx = &types.TransferTransaction{} // TODO

		Convey("The transaction should not be pass if receiver not exist", func() {
			transfer(ctx, tx)
			// TODO check
		})

		tx = &types.TransferTransaction{} // TODO

		Convey("The transaction should not be pass if there is not enough balance", func() {
			transfer(ctx, tx)
			// TODO check
		})
	})
}
