package txs

import (
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
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
				To:    types.Addr("0x0000000000000000000000000000000000000000"),
				Value: types.NewBigInt(0),
			},
			{
				To:     types.Addr("0x0000000000000000000000000000000000000000"),
				Value:  types.NewBigInt(0),
				Remark: make([]byte, 4096),
			},
		}
		transferTx := TransferTx(types.Addr("0x0000000000000000000000000000000000000000"), outputs, types.NewBigInt(0), 1, "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
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
				To:    types.Addr("0x833a907efe57af3040039c90f4a59946a0bb3d47"),
				Value: types.NewBigInt(1),
			},
			{
				To:     types.Addr("0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade"),
				Value:  types.NewBigInt(2),
				Remark: make([]byte, 1),
			},
		}
		transferTx := TransferTx(types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"), outputs, types.NewBigInt(0), 1, "187a9e180c7e950f3029a27e95f2cc0eabb6b73a28c706bafab75f0bd36c45ab2636092534cd5e8bc53d1d48b546f4aeb74c3e93e129855c3bb39020cfba2c7c1c")
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

		id1 := account.NewAccount(state, types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"))
		account.AddBalance(state, id1, big.NewInt(100))
		id2 := account.NewAccount(state, types.Addr("0x833a907efe57af3040039c90f4a59946a0bb3d47"))
		account.AddBalance(state, id2, big.NewInt(150))
		account.SaveBalance(state, types.Addr("0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade"), big.NewInt(200))

		outputs := []TransferOutput{
			{
				To:    types.Addr("0x833a907efe57af3040039c90f4a59946a0bb3d47"),
				Value: types.NewBigInt(50),
			},
			{
				To:     types.Addr("0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade"),
				Value:  types.NewBigInt(50),
				Remark: make([]byte, 1),
			},
		}
		transferTx := TransferTx(types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"), outputs, types.NewBigInt(0), 1, "86d924df60353825b12ae18a659a4de9df30555da9a1576d54fdfd9008d482397ce81cf6c83b9c1dcb8a733aa998f84f1c825bbd304f5c2774119d109e889cff1b")
		Convey("If a transfer transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(account.FetchBalance(state, id1).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchBalance(state, id2).Cmp(big.NewInt(200)), ShouldBeZeroValue)
					So(account.FetchBalance(state, types.Addr("0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade")).Cmp(big.NewInt(250)), ShouldBeZeroValue)
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
				rawTx := RawTransferTx(types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"), outputs, types.NewBigInt(0), 1, "86d924df60353825b12ae18a659a4de9df30555da9a1576d54fdfd9008d482397ce81cf6c83b9c1dcb8a733aa998f84f1c825bbd304f5c2774119d109e889cff1b")
				So(rawTx, ShouldResemble, EncodeTx(transferTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *TransferTransaction
					err := types.AminoCodec().UnmarshalBinary(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, transferTx)
				})
			})
		})
		Convey("If a transfer transaction from LikeChainID is valid", func() {
			outputs[0].To = id2
			transferTx := TransferTx(id1, outputs, types.NewBigInt(0), 1, "25ed8de68b496cbf6e0e012da6e3543c38f9ef941ed4a1ccb831874126fb422e0698d10c536c81d6adeccb5f77f3c3bfa1180650011ecdeae65a08dc1d87e81e1c")
			Convey("CheckTx should succeed", func() {
				r := transferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(transferTx))
					r := transferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(account.FetchBalance(state, id1).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchBalance(state, id2).Cmp(big.NewInt(200)), ShouldBeZeroValue)
					So(account.FetchBalance(state, types.Addr("0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade")).Cmp(big.NewInt(250)), ShouldBeZeroValue)
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
			transferTx.From = types.Addr("0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade")
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
			transferTx.Sig = &TransferJSONSignature{Sig("a176b40b9ffd6f4868ce5d608f3c6a560e8507abd57624449c0c595a3ea28a351e522c41e24beaccabc70e50999e9f653cecd7ce272eaf6ad051df1b758886fc1b")}
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
			transferTx.Sig = &TransferJSONSignature{Sig("45007b7c7c4ebb4ddb9580c1d11b1d8fd905c31bf4f709cdcdfe3bd12a3ebd860340c417662ec43b14cfe42391ea99c064eedf83ec396ab2244777c21dcb01691b")}
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
			outputs[1].To = types.IDStr("1MaeSeg6YEf0bkKy0FOh8MbnDqQ=")
			transferTx.Sig = &TransferJSONSignature{Sig("0f260b4879eaba6f3e4d467f18b18d1b3bfef45eb8eb7534706343b599867ba55312fb453eb794e227ffa2def17af2b456ed43f9d499e48de8bf42e60d60fd321c")}
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
			transferTx.Sig = &TransferJSONSignature{Sig("dee2c6eae68ca025c28169e1effea0b876ab692ce506069516e7d05dfd4ce3c9559a6dfa9bb02988fa953f8fcbcf9887d6d3038b970e6bd15444a6ef570061cf1b")}
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
			transferTx.Sig = &TransferJSONSignature{Sig("e7d3df2ee3c55fece54517a5f650c40eda930f7b442f7ea8a18905b0d3e024527ee6357f2d84ebaba55c39b91b34d62df16a3620de5c8ce417d129c70739937d1b")}
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
