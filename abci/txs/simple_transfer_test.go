package txs

import (
	"math/big"
	"strings"
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

func TestSimpleTransferValidateFormat(t *testing.T) {
	Convey("For a SimpleTransfer transaction", t, func() {
		tx := SimpleTransferTx(Alice.Address, Bob.ID, types.NewBigInt(50), "this is spartaaaaaaaa", types.NewBigInt(1), 1, "2af037daf098a5019f28a83196e28818faa74d5ec788953f8332036688b431d720a523246dc32c40a9f0c2da882a9cc68b44d090c26477827213ded82240e0101b")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(tx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has nil From field", func() {
			tx.From = nil
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil fee", func() {
			tx.Fee.Int = nil
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has negative fee", func() {
			tx.Fee.Int = big.NewInt(-1)
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has 0 fee", func() {
			tx.Fee.Int = big.NewInt(0)
			Convey("The validation should succedd", func() {
				So(tx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has fee 2^256-1", func() {
			tx.Fee.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			tx.Fee.Int.Sub(tx.Fee.Int, big.NewInt(1))
			Convey("The validation should succeed", func() {
				So(tx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has fee 2^256", func() {
			tx.Fee.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil value", func() {
			tx.Value.Int = nil
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has negative value", func() {
			tx.Value.Int = big.NewInt(-1)
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has 0 value", func() {
			tx.Value.Int = big.NewInt(0)
			Convey("The validation should succedd", func() {
				So(tx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has value 2^256-1", func() {
			tx.Value.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			tx.Value.Int.Sub(tx.Fee.Int, big.NewInt(1))
			Convey("The validation should succeed", func() {
				So(tx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has value 2^256", func() {
			tx.Value.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil signature", func() {
			tx.Sig = nil
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has output with remark with length 0", func() {
			tx.Remark = ""
			Convey("The validation should succeed", func() {
				So(tx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has output with remark with length 4096", func() {
			tx.Remark = strings.Repeat("w", 4096)
			Convey("The validation should succeed", func() {
				So(tx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has output with remark with length 4097", func() {
			tx.Remark = strings.Repeat("w", 4097)
			Convey("The validation should fail", func() {
				So(tx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestSimpleTransferJSONSignature(t *testing.T) {
	Convey("For a SimpleTransfer transaction with JSON signature", t, func() {
		Convey("If the transaction is valid with some remark", func() {
			tx := SimpleTransferTx(Alice.Address, Bob.ID, types.NewBigInt(50), "this is spartaaaaaaaa", types.NewBigInt(1), 1, "2af037daf098a5019f28a83196e28818faa74d5ec788953f8332036688b431d720a523246dc32c40a9f0c2da882a9cc68b44d090c26477827213ded82240e0101b")
			Convey("Address recovery should succeed", func() {
				recoveredAddr, err := tx.Sig.RecoverAddress(tx)
				So(err, ShouldBeNil)
				Convey("The recovered address should be the From address of the transfer transaction", func() {
					So(recoveredAddr, ShouldResemble, tx.From)
				})
			})
		})
		Convey("If the transaction is valid with no remark", func() {
			tx := SimpleTransferTx(Alice.Address, Bob.ID, types.NewBigInt(50), "", types.NewBigInt(1), 1, "f4d07f91ab07941284d6a5385d361ddf6a8112529b036b5d55ed066615a231293a86add2311a5305e816f1f85bedfeebe0d0e5ec62e4fdb60addb99d81e7c4a61b")
			Convey("Address recovery should succeed", func() {
				recoveredAddr, err := tx.Sig.RecoverAddress(tx)
				So(err, ShouldBeNil)
				Convey("The recovered address should be the From address of the transfer transaction", func() {
					So(recoveredAddr, ShouldResemble, tx.From)
				})
			})
		})
	})
}

func TestSimpleTransferEIP712Signature(t *testing.T) {
	Convey("For a SimpleTransfer transaction with EIP-712 signature", t, func() {
		Convey("If the transaction is valid with some remark", func() {
			tx := &SimpleTransferTransaction{
				From:   Alice.Address,
				To:     Bob.ID,
				Value:  types.NewBigInt(50),
				Remark: "this is spartaaaaaaaa",
				Fee:    types.NewBigInt(1),
				Nonce:  1,
				Sig:    &SimpleTransferEIP712Signature{SigEIP712("813156e45f00f8e885daba596da70854f6d572dd46a37460a104b383d8cd76734c59eda0e177752371e187c66249def595e3dc191433eb0c35e8bbc2af6c0da31b")},
			}
			Convey("Address recovery should succeed", func() {
				recoveredAddr, err := tx.Sig.RecoverAddress(tx)
				So(err, ShouldBeNil)
				Convey("The recovered address should be the From address of the transfer transaction", func() {
					So(recoveredAddr, ShouldResemble, tx.From)
				})
			})
		})
		Convey("If the transaction is valid with no remark", func() {
			tx := &SimpleTransferTransaction{
				From:   Alice.Address,
				To:     Bob.ID,
				Value:  types.NewBigInt(50),
				Remark: "",
				Fee:    types.NewBigInt(1),
				Nonce:  1,
				Sig:    &SimpleTransferEIP712Signature{SigEIP712("60881422ecacad9fa37350e7f25e17caa79dd17be547c2ff0776bd3d942a22326454c37c8f25a9977115f851e2b056ed76185322d9cd574f00d41f3244f667641b")},
			}
			Convey("Address recovery should succeed", func() {
				recoveredAddr, err := tx.Sig.RecoverAddress(tx)
				So(err, ShouldBeNil)
				Convey("The recovered address should be the From address of the transfer transaction", func() {
					So(recoveredAddr, ShouldResemble, tx.From)
				})
			})
		})
	})
}

func TestCheckAndDeliverSimpleTransfer(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		account.NewAccountFromID(state, Alice.ID, Alice.Address)
		account.NewAccountFromID(state, Bob.ID, Bob.Address)
		account.AddBalance(state, Alice.ID, big.NewInt(200))
		tx := SimpleTransferTx(Alice.Address, Bob.ID, types.NewBigInt(50), "this is spartaaaaaaaa", types.NewBigInt(1), 1, "2af037daf098a5019f28a83196e28818faa74d5ec788953f8332036688b431d720a523246dc32c40a9f0c2da882a9cc68b44d090c26477827213ded82240e0101b")
		Convey("If a SimpleTransfer transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(150)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(50)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
					Convey("CheckTx again should return Duplicated", func() {
						r := tx.CheckTx(state)
						So(r.Code, ShouldEqual, response.SimpleTransferDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := tx.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.SimpleTransferDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
						})
					})
				})
			})
			Convey("RawSimpleTransferTx should return the same encoded raw tx", func() {
				rawTx := RawSimpleTransferTx(Alice.Address, Bob.ID, types.NewBigInt(50), "this is spartaaaaaaaa", types.NewBigInt(1), 1, "2af037daf098a5019f28a83196e28818faa74d5ec788953f8332036688b431d720a523246dc32c40a9f0c2da882a9cc68b44d090c26477827213ded82240e0101b")
				So(rawTx, ShouldResemble, EncodeTx(tx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *SimpleTransferTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, tx)
				})
			})
		})
		Convey("If a transfer transaction from LikeChainID is valid", func() {
			tx.From = Alice.ID
			tx.Sig = &SimpleTransferJSONSignature{Sig("582815c336ec854fdd2306be5ad4166ae38a8c290daa20726522c191ea7535603bd262ba53f5ccee9edd8ebaefd0b0550d45515bfe240b00981c94a9b723da671c")}
			Convey("CheckTx should succeed", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(150)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(50)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a SimpleTransfer transaction has invalid format", func() {
			tx.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.SimpleTransferInvalidFormat.Code)
				Convey("DeliverTx should InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.SimpleTransferInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a SimpleTransfer transaction has unregistered sender", func() {
			tx.From = Carol.Address
			tx.Sig = &SimpleTransferJSONSignature{Sig("cc926aadee3c732c317453f49e735153908d31ca8d8d48067fe3567834ea45820b2410bf8bfc12fc42f03894595095bae13537392fbc1a55b3224ff5f4c2dd541c")}
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.SimpleTransferSenderNotRegistered.Code)
				Convey("DeliverTx should SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.SimpleTransferSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a SimpleTransfer transaction has signature not from the sender", func() {
			tx.Sig = &SimpleTransferJSONSignature{Sig("0a741be73c2a58b4f3d832d44c0dcea34fb247b4f0a5ac3c3ee36a9c86810d0b6618842e6151a5bcda4e76898946ad5aeddd97860400ac8f32415da2e835d8e31b")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.SimpleTransferInvalidSignature.Code)
				Convey("DeliverTx should InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.SimpleTransferInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a SimpleTransfer transaction has nonce exceeding next nonce", func() {
			tx.Nonce = 2
			tx.Sig = &SimpleTransferJSONSignature{Sig("af04d983bd17905f3c9f83ed1c82f67b42684a92d013cac4ff9cfd98953fe4673f8b9543d3dd4401946262054286c67c40a129371d192ec4a9b2ec75f5052e451c")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.SimpleTransferInvalidNonce.Code)
				Convey("DeliverTx should InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.SimpleTransferInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a SimpleTransfer transaction has unregistered LikeChainID as receiver", func() {
			tx.To = Carol.ID
			tx.Sig = &SimpleTransferJSONSignature{Sig("2a87e86c6652bc39870c0666b0d31c7c416177804e914573015aa7662e7b039b3a0037a781b9fc99bf82bb8209d468b2af24bf1c55bec88d4a676ccb0f72c2461c")}
			Convey("CheckTx should return InvalidReceiver", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.SimpleTransferInvalidReceiver.Code)
				Convey("DeliverTx should InvalidReceiver", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.SimpleTransferInvalidReceiver.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a SimpleTransfer transaction has value and fee sum exactly equals to sender's balance", func() {
			tx.Fee = types.NewBigInt(150)
			tx.Sig = &SimpleTransferJSONSignature{Sig("e6029333f6320b492e334f7ba5d8d1a53a772432f6b814f765577db5bc5a585f6e4634148ab30d00716748e1bee45615d8ca15b92f4c4e309786bb1ed86c7d7f1c")}
			Convey("CheckTx should succeed", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(150)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(50)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a SimpleTransfer transaction has value and fee sum just more than sender's balance", func() {
			tx.Fee = types.NewBigInt(151)
			tx.Sig = &SimpleTransferJSONSignature{Sig("2e1dad51e73c9064a1172d1be460577370d276cee608f59e0af9ba8ef0ba6eba5a8eab7a5574a8fcd27a1b9ed70ca2af16450eecef2aa8ede4f3aedcc5210b2a1c")}
			Convey("CheckTx should return NotEnoughBalance", func() {
				r := tx.CheckTx(state)
				So(r.Code, ShouldEqual, response.SimpleTransferNotEnoughBalance.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(tx))
					r := tx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.SimpleTransferNotEnoughBalance.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
	})
}
