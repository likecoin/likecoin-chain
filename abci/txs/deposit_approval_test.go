package txs

import (
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	"github.com/tendermint/tendermint/crypto/tmhash"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDepositApprovalValidateFormat(t *testing.T) {
	Convey("For a deposit approval transaction", t, func() {
		txHash := make([]byte, tmhash.Size)
		depositApprovalTx := DepositApprovalTx(Alice.Address, txHash, 1, "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
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
		depositApprovalTx := DepositApprovalTx(Alice.Address, txHash, 1, "f05ba8c73518c9d65fa07d38124fa6008a8771fd35ae95f77ca0cb260bb8e84812e8fe6e4d0b27c01b45193437fa27214c3c9b9e96981471767458546c1cb1b41b")
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

		inputs1 := []deposit.Input{
			{
				FromAddr: *Alice.Address,
				Value:    types.NewBigInt(10),
			},
			{
				FromAddr: *Bob.Address,
				Value:    types.NewBigInt(20),
			},
		}
		inputs2 := []deposit.Input{
			{
				FromAddr: *Carol.Address,
				Value:    types.NewBigInt(10),
			},
			{
				FromAddr: *Mallory.Address,
				Value:    types.NewBigInt(20),
			},
		}

		account.NewAccountFromID(state, Alice.ID, Alice.Address)
		account.NewAccountFromID(state, Bob.ID, Bob.Address)
		account.NewAccountFromID(state, Carol.ID, Carol.Address)
		account.NewAccountFromID(state, Dave.ID, Dave.Address)
		account.NewAccountFromID(state, Erin.ID, Erin.Address)
		deposit.SetDepositApprovers(state, []deposit.Approver{
			// Specially set value, so the sums is larger than 2^32, and the sum of the first 3 will be just larger
			// than 2/3 of the sum
			{ID: Alice.ID, Weight: 2000000000},
			{ID: Bob.ID, Weight: 4000000001},
			{ID: Carol.ID, Weight: 2000000000},
			{ID: Dave.ID, Weight: 4000000000},
		})
		depositTx1 := DepositTx(Alice.Address, 1337, inputs1, 1, "af6109def6aebe496980b7cda873781fecf0c1d8f1b432b736a91f67b0372ee50d74c30d6718bdfcf0097aefeed5dbeebde3200030edbc7c064e29d76d5e60b11b")
		depositTx2 := DepositTx(Dave.ID, 1337, inputs2, 1, "2dce3d4d6613937900f40f70fb3e5d9f3f4e6f2c1c6bd407bc354a5966d252220395f1e6d8d13d77bb6df00c29fde0dd02ae818a32ffb74e969d19054716c6731c")
		depositTxHash1 := utils.HashRawTx(EncodeTx(depositTx1))
		depositTxHashDecoded, _ := utils.Hex2Bytes("ce03d2e7922d2c81842f3a937a06eef5e5616a2aff29c5fc45ada469060c885a")
		So(depositTxHash1, ShouldResemble, depositTxHashDecoded)
		depositTxHash2 := utils.HashRawTx(EncodeTx(depositTx2))
		depositTxHashDecoded, _ = utils.Hex2Bytes("af5885cdc7ea2c44073bab8064ea65113c372a6c90af0ad0e263fc60be1f3d6e")
		So(depositTxHash2, ShouldResemble, depositTxHashDecoded)
		r := depositTx1.DeliverTx(state, depositTxHash1)
		So(r.Code, ShouldEqual, 0)

		depositApprovalTx1 := DepositApprovalTx(Bob.Address, depositTxHash1, 1, "63b529352d681ae1c08c50e18702ab290ba2ebc0a77b0ec5ce86234c02185a5777b738199e167eac330636d7e899303f64f75b0609f4218faa3eef06fb23b1931c")
		Convey("If a deposit approval transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := depositApprovalTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(deposit.GetDepositApproval(state, Bob.ID, depositTx1.Proposal.BlockNumber), ShouldResemble, depositTxHash1)
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
				rawTx := RawDepositApprovalTx(Bob.Address, depositTxHash1, 1, "63b529352d681ae1c08c50e18702ab290ba2ebc0a77b0ec5ce86234c02185a5777b738199e167eac330636d7e899303f64f75b0609f4218faa3eef06fb23b1931c")
				So(rawTx, ShouldResemble, EncodeTx(depositApprovalTx1))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *DepositApprovalTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, depositApprovalTx1)
				})
			})
		})
		Convey("If a deposit approval transaction from LikeChainID is valid", func() {
			depositApprovalTx1.Approver = Bob.ID
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("06959fb8f14ac2a2e38194cbce81fb482129ab7f588e8bfaf1a769af0ec5fcc61afcb4f10796cea629ea2c49f2cb951b4b9ff86b5fd717008032f46c946a7df31b")}
			Convey("CheckTx should succeed", func() {
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
					r := depositApprovalTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(deposit.GetDepositApproval(state, Bob.ID, depositTx1.Proposal.BlockNumber), ShouldResemble, depositTxHash1)
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
			depositApprovalTx1.Approver = Mallory.Address
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("e74dfc17ca81d5430abaa7db4a714d0e9098dc69e1186ce7bac3a7c7ce3e2b8c24a24fed905c3bbf3b3d8c20d8cffc53309ed514b8b5191433e818f28a6d77651b")}
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
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("ef21709b824fc5ee1fe3f9af6b66a356c46875b237a5f20d6056008ac413325759daff0f92be27a3638de18e89a88b1c0e4ae6ff077f90df1e7e4a7e09d330ab1b")}
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
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("5714ec00a4fc7615faaee06d59d6796630e37120ff7d8f1c689aaab8f246b5fc0d2e0071118d3491d087d104836aae3caf1b1fa6b2741664f358d20f0a0b372c1b")}
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
			depositApprovalTx := DepositApprovalTx(Carol.ID, depositTxHash1, 1, "cff7853b921b49ccd4d4f9f3fd759f3785b992f7a37da37a8bb6eb109773b4c6075f002ec03153a66b2643ff8862d130d6224620cd7e6505ff111d89f389a1e91b")
			txHash := EncodeTx(depositApprovalTx)
			Convey("After DeliverTx, the deposit proposal should not be executed yet", func() {
				depositApprovalTx.DeliverTx(state, txHash)
				So(deposit.GetDepositExecution(state, depositTx1.Proposal.BlockNumber), ShouldBeNil)
				Convey("If after another deposit approval transaction, the proposal has total weight sum more than 2/3", func() {
					txHash := EncodeTx(depositApprovalTx1)
					Convey("After DeliverTx, the deposit proposal should be executed", func() {
						r := depositApprovalTx1.DeliverTx(state, txHash)
						So(r.Code, ShouldEqual, response.Success.Code)
						So(deposit.GetDepositExecution(state, depositTx1.Proposal.BlockNumber), ShouldResemble, depositTxHash1)
					})
				})
				Convey("If after another deposit approval transaction, the proposal has total weight still not yet sum more than 2/3", func() {
					depositApprovalTx := DepositApprovalTx(Dave.ID, depositTxHash1, 1, "ea1c130626f15f392877c29206c10095d488a7fc0ef91603dd1fe062038c23c414ed06829f431e8b50cc140f8a7f4403c26446d6695cca2b24b55620ccd24e301b")
					txHash := EncodeTx(depositApprovalTx)
					Convey("After DeliverTx, the deposit proposal should still not be executed yet", func() {
						r := depositApprovalTx.DeliverTx(state, txHash)
						So(r.Code, ShouldEqual, response.Success.Code)
						So(deposit.GetDepositExecution(state, depositTx1.Proposal.BlockNumber), ShouldBeNil)
					})
				})
			})
		})
		Convey("If after a deposit proposal execution, there is an deposit approval transaction for another proposal on the same block number", func() {
			r := depositTx2.DeliverTx(state, depositTxHash2)
			So(r.Code, ShouldEqual, response.Success.Code)
			deposit.ExecuteDepositProposal(state, depositTxHash1)
			depositApprovalTx := DepositApprovalTx(Bob.ID, depositTxHash2, 1, "d08a8cda9f44d2a3ff76b1723002fa2826b2c96ae53c67fb50c1eb676b22904f4f2b0012be5fadd41fa6a8aded020bf607222f72cc762a9747c1f08f08215b981b")
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
			r := depositTx2.DeliverTx(state, depositTxHash2)
			So(r.Code, ShouldEqual, response.Success.Code)
			depositApprovalTx := DepositApprovalTx(Alice.ID, depositTxHash2, 2, "d54b0599bc3261b588777a24fe8e25bcb8266bf44a5bb5879dd9e1d7870c5b592e7210196ac5413dd48c84c938724b9611ca29f00357f773b670e919c3418ea01c")
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
			r := depositTx2.DeliverTx(state, depositTxHash2)
			So(r.Code, ShouldEqual, response.Success.Code)
			txHash := utils.HashRawTx(EncodeTx(depositApprovalTx1))
			depositApprovalTx1.DeliverTx(state, txHash)
			depositApprovalTx := DepositApprovalTx(Bob.ID, depositTxHash2, 2, "46a5c70ac73632568f87dd029b57ab2d4e94404b399b77aac768b5098831882e131ed860d4e8901f7ff2b306aa02975d121af9d9b028c257a78d95c413615a711b")
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
			depositApprovalTx1.DepositTxHash = make([]byte, 32)
			depositApprovalTx1.Sig = &DepositApprovalJSONSignature{Sig("eba99d868db9489ab37d106d7a92c70b64b8c87f7e021f645ffb675dbfe0e57a23c6df5f8c77963c428895ce6e2cbbaac1c4a4452cda109f97707ee15c1a5f651b")}
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
