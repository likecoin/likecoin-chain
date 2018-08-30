package transfer

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
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
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a Transfer Transaction", t, func() {

		Convey("If it is a valid transaction", func() {
			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			Convey("CheckTx should return Code 0", func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, 0)
			})

			Convey("DeliverTx should return Code 0", func() {
				res := deliverTransfer(state, rawTx)

				So(res.Code, ShouldEqual, 0)

				Convey("Should be able to query the transaction info afterwards", func() {
					_ = res.Data // TODO: ID
					// TODO: query
				})
			})
		})

		Convey("If it is an invalid address format", func() {
			appCtx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code := response.TransferCheckTxInvalidFormat.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxInvalidFormat.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is an invalid signature version", func() {
			appCtx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code := response.TransferCheckTxInvalidSignature.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxInvalidSignature.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is an invalid signature format", func() {
			appCtx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code := response.TransferCheckTxInvalidSignature.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxInvalidSignature.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If it is a replayed transaction", func() {
			appCtx.Reset()

			rawTx := wrapTransferTransaction(&types.TransferTransaction{
				// TODO
			})

			code := response.TransferCheckTxDuplicated.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxDuplicated.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})
	})
}

func TestValidateTransferSignature(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	senderAddr := types.NewAddressFromHex("0x064b663abf9d74277a07aa7563a8a64a54de8c0a")

	Convey("Given a Transfer transaction", t, func() {
		tx := &types.TransferTransaction{
			From: senderAddr.ToIdentifier(),
			ToList: []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(
					types.NewAddressFromHex("0xbef509a0ab4a60111a8957322fee016cdf713ad2").ToIdentifier(),
					"1000000000000000000",
					"",
				),
			},
			Nonce: 1,
			Fee:   types.NewBigInteger("10000000000"),
			Sig:   types.NewSignatureFromHex("0x6f44fdcfc7c4516af854404bbad55a229dbc3898621146f6c737e38ca22117d81f564f425449cd167d38df480ea8e26441b9f5342a935973bd72a903927cc4641b"),
		}

		Convey("If its sender address is match with the signing address", func() {
			Convey("It should pass the validation", func() {
				So(validateTransferSignature(state, tx), ShouldBeTrue)
			})
		})

		Convey("If its sender address is not match with the signing address", func() {
			tx.From = types.NewAddressFromHex("0xbef509a0ab4a60111a8957322fee016cdf713ad2").ToIdentifier()

			Convey("It should not pass the validation", func() {
				So(validateTransferSignature(state, tx), ShouldBeFalse)
			})
		})

		Convey("If its sender identifier is a LikeChain ID", func() {
			tx.Sig = types.NewSignatureFromHex("0x61235afd564c0b96c44342edd4456ef26ed142da9528bd88a7c321aa8595c96044a7a9c7747e049f88e89df6db0419c0aaee6c4d795c5540424379efd7e7a6731c")

			Convey("If the LikeChain ID has been bound to the signing address", func() {
				aliceID := types.NewLikeChainID([]byte("alice"))
				account.NewAccountFromID(state, aliceID, senderAddr.ToEthereum())
				tx.From = aliceID.ToIdentifier()

				Convey("It should pass the validation", func() {
					So(validateTransferSignature(state, tx), ShouldBeTrue)
				})
			})

			Convey("If the LikeChain ID has not been bound to the signing address", func() {
				malloryID := types.NewLikeChainID([]byte("mallory"))
				tx.From = malloryID.ToIdentifier()

				Convey("It should not pass the validation", func() {
					So(validateTransferSignature(state, tx), ShouldBeFalse)
				})
			})
		})
	})
}

func TestValidateTransferTransactionFormat(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	aliceID := types.NewLikeChainID([]byte("alice"))
	account.SaveBalance(state, *aliceID, big.NewInt(1000000000000000000))
	account.IncrementNextNonce(state, *aliceID)

	bobID := types.NewLikeChainID([]byte("bob"))
	account.SaveBalance(state, *bobID, big.NewInt(0))

	Convey("Given a Transfer transaction", t, func() {
		tx := &types.TransferTransaction{
			From: aliceID.ToIdentifier(),
			ToList: []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(bobID.ToIdentifier(), "1000000000000000000", ""),
			},
			Nonce: 1,
			Fee:   types.NewBigInteger("10000000000"),
			Sig:   types.NewSignatureFromHex("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		}

		Convey("If its format is valid", func() {
			Convey("It should pass the validation", func() {
				So(validateTransferTransactionFormat(state, tx), ShouldBeTrue)
			})
		})

		Convey("If its sender is invalid", func() {
			tx.From = &types.Identifier{
				Id: &types.Identifier_LikeChainID{
					LikeChainID: &types.LikeChainID{Content: []byte{}},
				},
			}

			Convey("It should not pass the validation", func() {
				So(validateTransferTransactionFormat(state, tx), ShouldBeFalse)
			})
		})

		Convey("If its receivers are invalid", func() {
			tx.ToList = []*types.TransferTransaction_TransferTarget{}

			Convey("It should not pass the validation", func() {
				So(validateTransferTransactionFormat(state, tx), ShouldBeFalse)
			})
		})

		Convey("If its signature format is invalid", func() {
			tx.Sig = types.NewSignatureFromHex("")

			Convey("It should not pass the validation", func() {
				So(validateTransferTransactionFormat(state, tx), ShouldBeFalse)
			})
		})
	})
}

func TestTransfer(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a valid Transfer transaction", t, func() {
		tx := &types.TransferTransaction{} // TODO

		Convey("The transaction should be pass", func() {
			transfer(state, tx)
			// TODO check
		})

		Convey("But the same transaction cannot be replayed", func() {
			transfer(state, tx)
			// TODO check
		})
	})

	Convey("Given an invalid Transfer transaction", t, func() {
		appCtx.Reset()
		tx := &types.TransferTransaction{} // TODO

		Convey("The transaction should not be pass if sender not exist ", func() {
			transfer(state, tx)
			// TODO check
		})

		tx = &types.TransferTransaction{} // TODO

		Convey("The transaction should not be pass if receiver not exist", func() {
			transfer(state, tx)
			// TODO check
		})

		tx = &types.TransferTransaction{} // TODO

		Convey("The transaction should not be pass if there is not enough balance", func() {
			transfer(state, tx)
			// TODO check
		})
	})
}
