package transfer

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

	Convey("Given an empty Transaction, CheckTx and DeliverTx should panic", t, func() {
		rawTx := &types.Transaction{}

		So(func() { checkTransfer(state, rawTx) }, ShouldPanic)
		So(func() { deliverTransfer(state, rawTx) }, ShouldPanic)
	})

	aliceID := types.NewLikeChainID([]byte("alice"))
	bobID := types.NewLikeChainID([]byte("bob"))
	carolID := types.NewLikeChainID([]byte("carol"))

	Convey("Given a Transfer Transaction", t, func() {
		appCtx.Reset()
		account.NewAccountFromID(state, aliceID, common.HexToAddress("0x064b663abf9d74277a07aa7563a8a64a54de8c0a"))
		account.AddBalance(state, aliceID, big.NewInt(9000000000000000000))
		account.NewAccountFromID(state, bobID, common.HexToAddress("0xbef509a0ab4a60111a8957322fee016cdf713ad2"))
		account.NewAccountFromID(state, carolID, common.HexToAddress("0xba0ad74ab6cfea30e0cfa4998392873ad1a11388"))

		rawTx := wrapTransferTransaction(&types.TransferTransaction{
			From: aliceID.ToIdentifier(),
			ToList: []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(bobID.ToIdentifier(), "1000000000000000000", ""),
				types.NewTransferTarget(carolID.ToIdentifier(), "1000000000000000000", ""),
			},
			Nonce: 1,
			Fee:   types.NewBigInteger("1"),
			Sig:   types.NewSignatureFromHex("0xcfa49131425aba2a3c089a69db5b9d4c1f793d9f3c448084564304485df2b5da46bb2274730cb8bf81b794e77096dbee9156b52c9592e95b0bd9b3dc054e7fb91b"),
		})

		Convey("If it is a valid transaction", func() {
			Convey("CheckTx should return code 0", func() {
				res := checkTransfer(state, rawTx)
				So(res.Code, ShouldEqual, 0)
			})

			Convey("For DeliverTx", func() {
				res := deliverTransfer(state, rawTx)

				Convey("Should return Code 0", func() {
					So(res.Code, ShouldEqual, 0)
				})

				Convey("Balance of those accounts in the state should be updated correctly ", func() {
					aliceBalance := account.FetchBalance(state, *aliceID)
					So(aliceBalance.String(), ShouldEqual, "6999999999999999999")

					bobBalance := account.FetchBalance(state, *bobID)
					So(bobBalance.String(), ShouldEqual, "1000000000000000000")

					carolBalance := account.FetchBalance(state, *carolID)
					So(carolBalance.String(), ShouldEqual, "1000000000000000000")
				})

				Convey("If replay this transaction", func() {
					Convey("For CheckTx", func() {
						res := checkTransfer(state, rawTx)

						code := response.TransferCheckTxDuplicated.Code
						Convey(fmt.Sprintf("Should return Code %d", code), func() {
							So(res.Code, ShouldEqual, code)
						})
					})

					Convey("For DeliverTx", func() {
						res := deliverTransfer(state, rawTx)

						code := response.TransferDeliverTxDuplicated.Code
						Convey(fmt.Sprintf("Should return Code %d", code), func() {
							So(res.Code, ShouldEqual, code)
						})
					})
				})
			})
		})

		Convey("If its format is invalid", func() {
			rawTx.GetTransferTx().ToList = []*types.TransferTransaction_TransferTarget{}

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

		Convey("If the sender is not registered", func() {
			rawTx.GetTransferTx().From = types.NewLikeChainID([]byte("mallory")).ToIdentifier()

			code := response.TransferCheckTxSenderNotRegistered.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxSenderNotRegistered.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If its signature is invalid", func() {
			rawTx.GetTransferTx().Sig = types.NewZeroSignature()

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

		Convey("If its nonce is invalid", func() {
			tx := rawTx.GetTransferTx()
			tx.Nonce = 2
			tx.Sig = types.NewSignatureFromHex("0x2dddc89b34896c0ab3622a3d2c4b1192b85606b6ed72a9f7e372ce0795498f830c4823cd461d80169ff666a6a55f23153dc820fc42a322ce1cad63b48cfc42a41b")

			code := response.TransferCheckTxInvalidNonce.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxInvalidNonce.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If the sender balance is not enough", func() {
			account.SaveBalance(state, *aliceID, big.NewInt(0))

			code := response.TransferCheckTxNotEnoughBalance.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxNotEnoughBalance.Code
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
			Sig:   types.NewZeroSignature(),
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
			tx.ToList = []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(nil, "0", ""),
			}

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
