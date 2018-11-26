package txs

import (
	// "math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDepositApprovalValidateFormat(t *testing.T) {
	Convey("For a deposit approval transaction", t, func() {
		txHash := make([]byte, 20)
		depositApprovalTx := DepositApprovalTx(fixture.Alice.Address, txHash, 1, "f05ba8c73518c9d65fa07d38124fa6008a8771fd35ae95f77ca0cb260bb8e84812e8fe6e4d0b27c01b45193437fa27214c3c9b9e96981471767458546c1cb1b41b")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(depositApprovalTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has nil signature", func() {
			depositApprovalTx.Sig = nil
			Convey("The validation should fail", func() {
				So(depositApprovalTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil Approver field", func() {
			depositApprovalTx.Approver = nil
			Convey("The validation should fail", func() {
				So(depositApprovalTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil DepositTxHash field", func() {
			depositApprovalTx.DepositTxHash = nil
			Convey("The validation should fail", func() {
				So(depositApprovalTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has DepositTxHash field with wrong length", func() {
			depositApprovalTx.DepositTxHash = make([]byte, 21)
			Convey("The validation should fail", func() {
				So(depositApprovalTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestDepositApprovalSignature(t *testing.T) {
	Convey("If a deposit approval transaction is valid", t, func() {
		txHash := make([]byte, 20)
		depositApprovalTx := DepositApprovalTx(fixture.Alice.Address, txHash, 1, "f05ba8c73518c9d65fa07d38124fa6008a8771fd35ae95f77ca0cb260bb8e84812e8fe6e4d0b27c01b45193437fa27214c3c9b9e96981471767458546c1cb1b41b")
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := depositApprovalTx.Sig.RecoverAddress(depositApprovalTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the Approver address of the deposit transaction", func() {
				So(recoveredAddr, ShouldResemble, depositApprovalTx.Approver)
			})
		})
	})
}

func TestCheckAndDeliverDepositApproval(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID
		id3 := fixture.Dave.ID
		id4 := fixture.Erin.ID

		inputs1 := []deposit.Input{
			{
				FromAddr: *fixture.Alice.Address,
				Value:    types.NewBigInt(10),
			},
			{
				FromAddr: *fixture.Bob.Address,
				Value:    types.NewBigInt(20),
			},
		}
		inputs2 := []deposit.Input{
			{
				FromAddr: *fixture.Carol.Address,
				Value:    types.NewBigInt(10),
			},
			{
				FromAddr: *fixture.Mallory.Address,
				Value:    types.NewBigInt(20),
			},
		}

		account.NewAccountFromID(state, id0, fixture.Alice.Address)
		account.NewAccountFromID(state, id1, fixture.Bob.Address)
		account.NewAccountFromID(state, id2, fixture.Carol.Address)
		account.NewAccountFromID(state, id3, fixture.Dave.Address)
		account.NewAccountFromID(state, id4, fixture.Erin.Address)
		deposit.SetDepositApprovers(state, []deposit.Approver{
			// Specially set value, so the sums is larger than 2^32, and the sum of the first 3 will be just larger
			// than 2/3 of the sum
			{ID: id0, Weight: 2000000000},
			{ID: id1, Weight: 4000000001},
			{ID: id2, Weight: 2000000000},
			{ID: id3, Weight: 4000000000},
		})
		depositTx1 := DepositTx(fixture.Alice.Address, 1337, inputs1, 1, "af6109def6aebe496980b7cda873781fecf0c1d8f1b432b736a91f67b0372ee50d74c30d6718bdfcf0097aefeed5dbeebde3200030edbc7c064e29d76d5e60b11b")
		depositTx2 := DepositTx(id3, 1337, inputs2, 1, "2dce3d4d6613937900f40f70fb3e5d9f3f4e6f2c1c6bd407bc354a5966d252220395f1e6d8d13d77bb6df00c29fde0dd02ae818a32ffb74e969d19054716c6731c")
		depositTxHash1 := utils.HashRawTx(EncodeTx(depositTx1))
		depositTxHashDecoded, _ := utils.Hex2Bytes("ce03d2e7922d2c81842f3a937a06eef5e5616a2a")
		So(depositTxHash1, ShouldResemble, depositTxHashDecoded)
		depositTxHash2 := utils.HashRawTx(EncodeTx(depositTx2))
		depositTxHashDecoded, _ = utils.Hex2Bytes("af5885cdc7ea2c44073bab8064ea65113c372a6c")
		So(depositTxHash2, ShouldResemble, depositTxHashDecoded)
		r := depositTx1.DeliverTx(state, depositTxHash1)
		So(r.Code, ShouldEqual, 0)
		r = depositTx2.DeliverTx(state, depositTxHash2)
		So(r.Code, ShouldEqual, 0)

		depositApprovalTx1 := DepositApprovalTx(fixture.Bob.Address, depositTxHash1, 1, "1607f6a849f698540cd8e4d62a430eb04b560156c4a28ddb4aa29b302228fc6c56556678bd6ac4b269a0e86c617f0e9a55b1ffd81e3a2fa522a49d50daa2307c1c")
		Convey("If a deposit approval transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(deposit.GetDepositApproval(state, id1, depositTx1.Proposal.BlockNumber), ShouldResemble, depositTxHash1)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					Convey("CheckTx again should return Duplicated", func() {
						r := depositApprovalTx1.CheckTx(state)
						So(r.Code, ShouldEqual, response.DepositApprovalDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := depositApprovalTx1.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.DepositApprovalDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
						})
					})
				})
			})
			Convey("RawTransferTx should return the same encoded raw tx", func() {
				rawTx := RawDepositApprovalTx(fixture.Bob.Address, depositTxHash1, 1, "1607f6a849f698540cd8e4d62a430eb04b560156c4a28ddb4aa29b302228fc6c56556678bd6ac4b269a0e86c617f0e9a55b1ffd81e3a2fa522a49d50daa2307c1c")
				So(rawTx, ShouldResemble, EncodeTx(depositApprovalTx1))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *DepositApprovalTransaction
					err := types.AminoCodec().UnmarshalBinary(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, depositApprovalTx1)
				})
			})
		})
		Convey("If a deposit approval transaction from LikeChainID is valid", func() {
			depositApprovalTx1.Approver = id1
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("b19858ee7439436a703560cb5fe50e884f7d57760cfb7c97d806525c3e2de31303973898da36364237305724bbf11a81ae77efb8d393f3f96cbd12b31d31a0011b")}
			Convey("CheckTx should succeed", func() {
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(deposit.GetDepositApproval(state, id1, depositTx1.Proposal.BlockNumber), ShouldResemble, depositTxHash1)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
				})
			})
		})
		Convey("If a deposit approval transaction has invalid format", func() {
			depositApprovalTx1.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalInvalidFormat.Code)
				Convey("DeliverTx should return InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit approval transaction has unregistered sender", func() {
			depositApprovalTx1.Approver = fixture.Mallory.Address
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("773cdf99716432ff9c583af746b989bf5c1a3723ad0222c426f625dd4c2882e730bc5ae9881cd21ae11646a7b70cdd2a6ffb6bf448d38e375b7e76855f8292601b")}
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalSenderNotRegistered.Code)
				Convey("DeliverTx should return SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit approval transaction has invalid signature", func() {
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit approval transaction has signature not from sender", func() {
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("18ddf99f2c7167a60e2d17bc891a1f1d5943cac309b17ff23664892a7e235f8722bc5daea93ee2e661404a7bcca8a96fd112e4f9404daa62060df554c74108ff1b")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit approval transaction has invalid nonce", func() {
			depositApprovalTx1.Nonce = 2
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("2d91655d35e4cdd39b0797340a9e345e670324c5dc3845ff1e59b782139c16706404fe3e3f62e8ef1b1a6ab22ae56828e218cc51501434321dc62c1a462e823b1c")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalInvalidNonce.Code)
				Convey("DeliverTx should return InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If after a deposit approval transaction, the proposal has total weight sum not yet more than 2/3", func() {
			txHash := EncodeTx(depositApprovalTx1)
			Convey("After DeliverTx, the deposit proposal should not be executed yet", func() {
				depositApprovalTx1.DeliverTx(state, txHash)
				So(deposit.GetDepositExecution(state, depositTx1.Proposal.BlockNumber), ShouldBeNil)
				Convey("If after another deposit approval transaction, the proposal has total weight sum more than 2/3", func() {
					depositApprovalTx := DepositApprovalTx(id2, depositTxHash1, 1, "0fbfa1af5db7c43f9eeb47ad4302e74ecc929f0c3d46ca5a4f58174c5a2e83d53847826fcd288fdd2cf3037d2b631b883c78806cb54048e1543f6717ce0ad9431c")
					txHash := EncodeTx(depositApprovalTx)
					Convey("After DeliverTx, the deposit proposal should be executed", func() {
						depositApprovalTx.DeliverTx(state, txHash)
						So(deposit.GetDepositExecution(state, depositTx1.Proposal.BlockNumber), ShouldResemble, depositTxHash1)
					})
				})
				Convey("If after another deposit approval transaction, the proposal has total weight still not yet sum more than 2/3", func() {
					depositApprovalTx := DepositApprovalTx(id3, depositTxHash1, 1, "c74e72c4fe9f154bea1c6ddc9c294173f7e206606910c43c4508fbdac19153d83244a3eaba20579aa1fab7d865f1ca0f5d28c5621398ba68b8ab20d5be4192511c")
					txHash := EncodeTx(depositApprovalTx)
					Convey("After DeliverTx, the deposit proposal should still not be executed yet", func() {
						depositApprovalTx.DeliverTx(state, txHash)
						So(deposit.GetDepositExecution(state, depositTx1.Proposal.BlockNumber), ShouldBeNil)
					})
				})
			})
		})
		Convey("If after a deposit proposal execution, there is an deposit approval transaction for another proposal on the same block number", func() {
			deposit.ExecuteDepositProposal(state, depositTxHash1)
			depositApprovalTx := DepositApprovalTx(id1, depositTxHash2, 1, "c6f3080133ec27e67cd823924ad0838b273776348ea76842a3750a62830c1a1479ba3bb1701e048c7b84f1f86a10fd82511ae5c1eaaa54c9691049f408e997581c")
			Convey("CheckTx should return AlreadyExecuted", func() {
				r := depositApprovalTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalAlreadyExecuted.Code)
				Convey("DeliverTx should return AlreadyExecuted", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx))
					r := depositApprovalTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalAlreadyExecuted.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit approval transaction has approver who has already proposed another proposal", func() {
			depositApprovalTx := DepositApprovalTx(id0, depositTxHash2, 2, "28904ab129d715d0ea4548e79178324089666762b346a113461b18234ec20ee62dadbce87e65d746055844abbf0d28495f3b754b64a267bce8c0fff18f45f6c51b")
			Convey("CheckTx should return DoubleApproval", func() {
				r := depositApprovalTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalDoubleApproval.Code)
				Convey("DeliverTx should return DoubleApproval", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx))
					r := depositApprovalTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalDoubleApproval.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit approval transaction has approver who has already approved another proposal", func() {
			txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
			depositApprovalTx1.DeliverTx(state, txHash)
			depositApprovalTx := DepositApprovalTx(id1, depositTxHash2, 2, "962f71ca2f2bde4b69ae25e80e60f94416f45a83a44defd6519b578197ae553b7387c6585e49dc8e2c8daa3fe8c86136ee71ed78b8e3088e4ac2c5d770a837001b")
			Convey("CheckTx should return DoubleApproval", func() {
				r := depositApprovalTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalDoubleApproval.Code)
				Convey("DeliverTx should return DoubleApproval", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx))
					r := depositApprovalTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalDoubleApproval.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit approval transaction has non-existing proposal as DepositTxHash", func() {
			depositApprovalTx1.DepositTxHash = make([]byte, 20)
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("0c3416ea3dea94bd5187fed0948bb82803038687d70c34b8ee9d2ae4ef98aad44c8db86df67ddf69786983e9aacc71bbe02d16fe0ee015d4d6bbb7f9edce80461b")}
			Convey("CheckTx should return ProposalNotExist", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositApprovalProposalNotExist.Code)
				Convey("DeliverTx should return ProposalNotExist", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositApprovalProposalNotExist.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
	})
}
