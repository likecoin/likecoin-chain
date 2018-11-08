package txs

import (
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransferValidateFormat(t *testing.T) {
	Convey("For a transfer transaction", t, func() {
		outputs := []TransferOutput{
			{
				To:    Bob.ID,
				Value: types.NewBigInt(0),
			},
			{
				To:     Carol.Address,
				Value:  types.NewBigInt(0),
				Remark: make([]byte, 4096),
			},
		}
		transferTx := TransferTx(Alice.Address, outputs, types.NewBigInt(0), 1, "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(transferTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has nil From field", func() {
			transferTx.From = nil
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil fee", func() {
			transferTx.Fee.Int = nil
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has negative fee", func() {
			transferTx.Fee.Int = big.NewInt(-1)
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil signature", func() {
			transferTx.Sig = nil
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has no outputs", func() {
			transferTx.Outputs = nil
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil receiver output", func() {
			transferTx.Outputs[0].To = nil
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil value output", func() {
			transferTx.Outputs[0].Value.Int = nil
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has negative value output", func() {
			transferTx.Outputs[0].Value.Int = big.NewInt(-1)
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has output with remark with length 4097", func() {
			transferTx.Outputs[0].Remark = make([]byte, 4097)
			Convey("The validation should fail", func() {
				So(transferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestTransferSignature(t *testing.T) {
	Convey("If a transfer transaction is valid", t, func() {
		outputs := []TransferOutput{
			{
				To:    Bob.ID,
				Value: types.NewBigInt(1),
			},
			{
				To:     Carol.Address,
				Value:  types.NewBigInt(2),
				Remark: make([]byte, 1),
			},
		}
		transferTx := TransferTx(Alice.Address, outputs, types.NewBigInt(0), 1, "52ade067038b6fa07aca930c4261ae683439f105ff7722015e10620c05bfe7c13253ba01573fef145f5f2a605f8db4cb68003116a07546bd3b52431c30aacffa1b")
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := transferTx.Sig.RecoverAddress(transferTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the From address of the transfer transaction", func() {
				So(recoveredAddr, ShouldResemble, transferTx.From)
			})
		})
	})
}

func TestCheckAndDeliverTransfer(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		account.NewAccountFromID(state, Alice.ID, Alice.Address)
		account.NewAccountFromID(state, Bob.ID, Bob.Address)
		account.AddBalance(state, Alice.ID, big.NewInt(100))
		account.AddBalance(state, Bob.ID, big.NewInt(150))
		account.SaveBalance(state, Carol.Address, big.NewInt(200))

		outputs := []TransferOutput{
			{
				To:    Bob.Address,
				Value: types.NewBigInt(50),
			},
			{
				To:     Carol.Address,
				Value:  types.NewBigInt(50),
				Remark: make([]byte, 1),
			},
		}
		transferTx := TransferTx(Alice.Address, outputs, types.NewBigInt(0), 1, "6c2e543b2bc323563aef9ae0ae8143dbfc1c2485a9d70fbf68ffc08ebac5ca840aa29974a1673806abb5a6d2a6d22ea5cc7e24b42bd391b6923b6126efb1652b1c")
		Convey("If a transfer transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(200)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Carol.Address).Cmp(big.NewInt(250)), ShouldBeZeroValue)
					Convey("CheckTx again should return Duplicated", func() {
						r := transferTx.CheckTx(state)
						So(r.Code, ShouldEqual, response.TransferDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := transferTx.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.TransferDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
						})
					})
				})
			})
			Convey("RawTransferTx should return the same encoded raw tx", func() {
				rawTx := RawTransferTx(Alice.Address, outputs, types.NewBigInt(0), 1, "6c2e543b2bc323563aef9ae0ae8143dbfc1c2485a9d70fbf68ffc08ebac5ca840aa29974a1673806abb5a6d2a6d22ea5cc7e24b42bd391b6923b6126efb1652b1c")
				So(rawTx, ShouldResemble, EncodeTx(transferTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *TransferTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx.From, ShouldResemble, transferTx.From)
					So(decodedTx.Fee.Cmp(transferTx.Fee.Int), ShouldBeZeroValue)
					So(decodedTx.Outputs, ShouldResemble, transferTx.Outputs)
					So(decodedTx.Nonce, ShouldResemble, transferTx.Nonce)
					So(decodedTx.Sig, ShouldResemble, transferTx.Sig)
				})
			})
		})
		Convey("If a transfer transaction from LikeChainID is valid", func() {
			outputs[0].To = Bob.ID
			transferTx := TransferTx(Alice.ID, outputs, types.NewBigInt(0), 1, "1cb0cb8ad4c21c93682d089819aacf1ccaf04b297912013f1578452dbdb9f4e47f769b99ae758d0cb68bac271a55575da33c672940a7b6e375ffa97f37aa25491c")
			Convey("CheckTx should succeed", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(200)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Carol.Address).Cmp(big.NewInt(250)), ShouldBeZeroValue)
				})
			})
		})
		Convey("If a transfer transaction has invalid format", func() {
			transferTx.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferInvalidFormat.Code)
				Convey("DeliverTx should InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a transfer transaction has unregistered sender", func() {
			transferTx.From = Carol.Address
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferSenderNotRegistered.Code)
				Convey("DeliverTx should SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a transfer transaction has invalid signature", func() {
			transferTx.Sig = &TransferJSONSignature{Sig("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferInvalidSignature.Code)
				Convey("DeliverTx should InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a transfer transaction has signature not from the sender", func() {
			transferTx.Sig = &TransferJSONSignature{Sig("74bbe427f9bba311d773f041a84c1e556c0e68e66694bf9bb86e46414536d1055632a0d0a1522a1ae3e81eed7082d214c2d9b1dcd44b955c4ee8d4dc1e01ba2c1c")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferInvalidSignature.Code)
				Convey("DeliverTx should InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a transfer transaction has nonce exceeding next nonce", func() {
			transferTx.Nonce = 2
			transferTx.Sig = &TransferJSONSignature{Sig("f05955c06ace300cde0fcb999ac488e427e8853a72d1998c828edd1c1718ca2b6168bed18c69ed3b8ad91f9f025464879d4cf92b96a0e6838bec096a17927ee81b")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferInvalidNonce.Code)
				Convey("DeliverTx should InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a transfer transaction has unregistered LikeChainID as output receiver", func() {
			outputs[1].To = Carol.ID
			transferTx.Sig = &TransferJSONSignature{Sig("08ac90f92a3987009501ba58c16289fbeb58f7c264ce3615107d0933e62c415c53bc5183819b9e88e4cf3034b5de9ec782bd4d6727f958efc3fb57ff391d0e951c")}
			Convey("CheckTx should return InvalidReceiver", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferInvalidReceiver.Code)
				Convey("DeliverTx should InvalidReceiver", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferInvalidReceiver.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a transfer transaction is transferring more balance than the sender has", func() {
			outputs[0].Value.Int = big.NewInt(51)
			transferTx.Sig = &TransferJSONSignature{Sig("3b37f8f60d79e428bc3f16715c54a6b249b79c29349f1bf2661c8e0ba2acbc9302356ddfe12c2842eed2cac961cd075053d5c6e33a5d19505a70fb8b8af2ce791b")}
			Convey("CheckTx should return NotEnoughBalance", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferNotEnoughBalance.Code)
				Convey("DeliverTx should NotEnoughBalance", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferNotEnoughBalance.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a transfer transaction has not enough balance for fee", func() {
			transferTx.Fee.Int = big.NewInt(1)
			transferTx.Sig = &TransferJSONSignature{Sig("c5663a3e9d9a32feb656980d066557e67133a79a66c3c246cf6ae546adbc58df77d837b7003c14680c1b8c66f523873c8c5468a96aa879ff8945dabe600b4e471b")}
			Convey("CheckTx should return NotEnoughBalance", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.TransferNotEnoughBalance.Code)
				Convey("DeliverTx should NotEnoughBalance", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.TransferNotEnoughBalance.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
	})
}
