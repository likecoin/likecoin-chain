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

		rawTx := (&types.TransferTransaction{
			From: fixture.Alice.ID.ToIdentifier(),
			ToList: []*types.TransferTransaction_TransferTarget{
				types.NewTransferTarget(fixture.Bob.ID.ToIdentifier(), "1000000000000000000", ""),
				types.NewTransferTarget(fixture.Carol.RawAddress.ToIdentifier(), "1000000000000000000", ""),
			},
			Nonce: 1,
			Fee:   types.NewBigInteger("1"),
			Sig:   types.NewSignatureFromHex("0x9fa395fef942fcade6200772872ce4a21f1a878120b459b500e9e27b848a8ffb4409545df564dfafed96828ce2bff77d73c43679b9e0980962bee8b5d2d93ed91c"),
		}).ToTransaction()

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
					So(res.Status, ShouldEqual, types.TxStatusSuccess)
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
					})
				})
			})
		})

		Convey("If it is a valid transaction with address as identifer", func() {
			tx := rawTx.GetTransferTx()
			tx.From = fixture.Alice.RawAddress.ToIdentifier()
			tx.Sig = types.NewSignatureFromHex("0xde3bb27979df8884ea69859dd841b91673814a474d6333e732979f047db72f8104a030aaca4667d31df8eff6dfd50775a5e93f5d44796351066394b70c6d50e11c")

			rawTxBytes, _ = proto.Marshal(rawTx)
			txHash = utils.HashRawTx(rawTxBytes)

			Convey("CheckTx should return code 0", func() {
				res := checkTransfer(state, rawTx)
				So(res.Code, ShouldEqual, 0)
			})

			Convey("DeliverTx should return Code 0", func() {
				res := deliverTransfer(state, rawTx, txHash)
				So(res.Code, ShouldEqual, 0)

				Convey("Balance should be updated correctly ", func() {
					aliceBalance := account.FetchBalance(state, fixture.Alice.ID.ToIdentifier())
					So(aliceBalance.String(), ShouldEqual, "6999999999999999999")
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
					So(res.Status, ShouldEqual, types.TxStatusFail)
				})
			})
		})

		Convey("If the sender is not registered", func() {
			tx := rawTx.GetTransferTx()
			tx.From = fixture.Mallory.ID.ToIdentifier()
			tx.Sig = types.NewSignatureFromHex("0x7e588f1532f1289cb52cddc57dc80d1ead7c9e84397faf7f3dc3f1500942470445ed0bc6c2adc05e54565da22b5820f611f9c4f2cacc5fbe885de887da68aaa21b")
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
			tx.Sig = types.NewSignatureFromHex("0x8621d00bffc0fddd29225dd25152958d4740a9263969951227ccae62b8d135370deb1917f490770cf0a3660e453dc390536cd5c2f1411179a2e80288bafdfdfe1b")

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
			tx.Sig = types.NewSignatureFromHex("0x6b69149b6fb251d14819f126d9e7d223db902875247ec6085b8b5f21a7a31ece40fbb9c657b242fe159a4c8b80b4ad2cb4a3a4e0778133938a7ad28f841ce7581b")

			Convey("If the LikeChain ID has been bound to the signing address", func() {
				account.NewAccountFromID(state, fixture.Alice.ID, fixture.Alice.Address)
				tx.From = fixture.Alice.ID.ToIdentifier()

				Convey("It should pass the validation", func() {
					So(validateTransferSignature(state, tx), ShouldBeTrue)
				})
			})

			Convey("If the LikeChain ID has not been bound to the signing address", func() {
				tx.From = fixture.Mallory.ID.ToIdentifier()

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

		Convey("If the remark of its receiver exceed the size limit", func() {
			tx.ToList[0].Remark = make([]byte, maxRemarkSize+1)

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
