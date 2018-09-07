package transfer

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/utils"

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
		So(func() { deliverTransfer(state, rawTx, nil) }, ShouldPanic)
	})

	Convey("Given a Transfer Transaction", t, func() {
		appCtx.Reset()
		account.NewAccountFromID(state, fixture.Alice.ID, fixture.Alice.Address)
		account.AddBalance(state, fixture.Alice.ID.ToIdentifier(), big.NewInt(9000000000000000000))
		account.NewAccountFromID(state, fixture.Bob.ID, fixture.Bob.Address)

		rawTx := wrapTransferTransaction(&types.TransferTransaction{
			From: fixture.Alice.ID.ToIdentifier(),
			ToList: []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(fixture.Bob.ID.ToIdentifier(), "1000000000000000000", ""),
				types.NewTransferTarget(fixture.Carol.RawAddress.ToIdentifier(), "1000000000000000000", ""),
			},
			Nonce: 1,
			Fee:   types.NewBigInteger("1"),
			Sig:   types.NewSignatureFromHex("0xf194fd5457c6a25bda697821283b9e2cc81279362215e448cc80d9c36c17cc2a3dc29ecf46f11f4263af85339cb47bc0c576ec32da184d396e8312b0fac0bb201b"),
		})

		rawTxBytes, _ := proto.Marshal(rawTx)
		txHash := utils.HashRawTx(rawTxBytes)

		Convey("If it is a valid transaction", func() {
			Convey("CheckTx should return code 0", func() {
				res := checkTransfer(state, rawTx)
				So(res.Code, ShouldEqual, 0)
			})

			Convey("For DeliverTx", func() {
				res := deliverTransfer(state, rawTx, txHash)

				Convey("Should return Code 0", func() {
					So(res.Code, ShouldEqual, 0)
				})

				Convey("Tx Status should be success", func() {
					So(GetStatus(state, txHash), ShouldEqual, types.TxStatusSuccess)
				})

				Convey("Balance of those accounts in the state should be updated correctly ", func() {
					aliceBalance := account.FetchBalance(state, fixture.Alice.ID.ToIdentifier())
					So(aliceBalance.String(), ShouldEqual, "6999999999999999999")

					bobBalance := account.FetchBalance(state, fixture.Bob.ID.ToIdentifier())
					So(bobBalance.String(), ShouldEqual, "1000000000000000000")

					carolBalance := account.FetchBalance(state, fixture.Carol.RawAddress.ToIdentifier())
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
						res := deliverTransfer(state, rawTx, txHash)

						code := response.TransferDeliverTxDuplicated.Code
						Convey(fmt.Sprintf("Should return Code %d", code), func() {
							So(res.Code, ShouldEqual, code)
						})

						Convey("Tx Status should be success", func() {
							So(GetStatus(state, txHash), ShouldEqual, types.TxStatusSuccess)
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
				res := deliverTransfer(state, rawTx, txHash)

				So(res.Code, ShouldEqual, code)

				Convey("Tx Status should be failed", func() {
					So(GetStatus(state, txHash), ShouldEqual, types.TxStatusFailed)
				})
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
				res := deliverTransfer(state, rawTx, txHash)

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
				res := deliverTransfer(state, rawTx, txHash)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If its nonce is invalid", func() {
			tx := rawTx.GetTransferTx()
			tx.Nonce = 2
			tx.Sig = types.NewSignatureFromHex("0xe3af48868498f69b792905b653ee11533af5dac613c1fb4a358a2776c277e7f800149e1f8af68ec5f4c01a1a24909bc339ce9f70cf0addc8635f54b3981a66561b0x2dddc89b34896c0ab3622a3d2c4b1192b85606b6ed72a9f7e372ce0795498f830c4823cd461d80169ff666a6a55f23153dc820fc42a322ce1cad63b48cfc42a41b")

			code := response.TransferCheckTxInvalidNonce.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxInvalidNonce.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx, txHash)

				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If the sender balance is not enough", func() {
			account.SaveBalance(state, fixture.Alice.ID.ToIdentifier(), big.NewInt(0))

			code := response.TransferCheckTxNotEnoughBalance.Code
			Convey(fmt.Sprintf("CheckTx should return Code %d", code), func() {
				res := checkTransfer(state, rawTx)

				So(res.Code, ShouldEqual, code)
			})

			code = response.TransferDeliverTxNotEnoughBalance.Code
			Convey(fmt.Sprintf("DeliverTx should return Code %d", code), func() {
				res := deliverTransfer(state, rawTx, txHash)

				So(res.Code, ShouldEqual, code)
			})
		})
	})
}

func TestValidateTransferSignature(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a Transfer transaction", t, func() {
		tx := &types.TransferTransaction{
			From: fixture.Alice.RawAddress.ToIdentifier(),
			ToList: []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(
					fixture.Bob.RawAddress.ToIdentifier(),
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
			tx.From = types.NewZeroAddress().ToIdentifier()

			Convey("It should not pass the validation", func() {
				So(validateTransferSignature(state, tx), ShouldBeFalse)
			})
		})

		Convey("If its sender identifier is a LikeChain ID", func() {
			tx.Sig = types.NewSignatureFromHex("0x61235afd564c0b96c44342edd4456ef26ed142da9528bd88a7c321aa8595c96044a7a9c7747e049f88e89df6db0419c0aaee6c4d795c5540424379efd7e7a6731c")

			Convey("If the LikeChain ID has been bound to the signing address", func() {
				account.NewAccountFromID(state, fixture.Alice.ID, fixture.Alice.Address)
				tx.From = fixture.Alice.ID.ToIdentifier()

				Convey("It should pass the validation", func() {
					So(validateTransferSignature(state, tx), ShouldBeTrue)
				})
			})

			Convey("If the LikeChain ID has not been bound to the signing address", func() {
				mallory := fixture.NewUser("mallory", "")
				tx.From = mallory.ID.ToIdentifier()

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

	account.SaveBalance(state, fixture.Alice.ID.ToIdentifier(), big.NewInt(1000000000000000000))
	account.IncrementNextNonce(state, fixture.Alice.ID)

	account.SaveBalance(state, fixture.Bob.ID.ToIdentifier(), big.NewInt(0))

	Convey("Given a Transfer transaction", t, func() {
		tx := &types.TransferTransaction{
			From: fixture.Alice.ID.ToIdentifier(),
			ToList: []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(fixture.Bob.ID.ToIdentifier(), "1000000000000000000", ""),
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
