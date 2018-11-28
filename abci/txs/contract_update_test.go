package txs

import (
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/contract"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContractUpdateValidateFormat(t *testing.T) {
	Convey("For a contract update transaction", t, func() {
		contractUpdateTx := ContractUpdateTx(Alice.Address, 1, types.Addr("0x1111111111111111111111111111111111111111"), 1, "81229384bdbe2ecc3572cb6deda168f2f0723161bc872bc53a9ef641ca2d087f1cd4aad56ddc7690550bc64aee0de76820d476b374699d4bd71dad9e7b5be8d81b")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(contractUpdateTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has nil signature", func() {
			contractUpdateTx.Sig = nil
			Convey("The validation should fail", func() {
				So(contractUpdateTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil Proposer field", func() {
			contractUpdateTx.Proposer = nil
			Convey("The validation should fail", func() {
				So(contractUpdateTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestContractUpdateSignature(t *testing.T) {
	Convey("If a contract update transaction is valid", t, func() {
		contractUpdateTx := ContractUpdateTx(Alice.Address, 1, types.Addr("0x1111111111111111111111111111111111111111"), 1, "81229384bdbe2ecc3572cb6deda168f2f0723161bc872bc53a9ef641ca2d087f1cd4aad56ddc7690550bc64aee0de76820d476b374699d4bd71dad9e7b5be8d81b")
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := contractUpdateTx.Sig.RecoverAddress(contractUpdateTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the Proposer address of the contract update transaction", func() {
				So(recoveredAddr, ShouldResemble, contractUpdateTx.Proposer)
			})
		})
	})
}

func TestCheckAndDeliverContractUpdate(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		account.NewAccountFromID(state, Alice.ID, Alice.Address)
		account.NewAccountFromID(state, Bob.ID, Bob.Address)
		account.NewAccountFromID(state, Carol.ID, Carol.Address)
		contract.SetContractUpdaters(state, []contract.Updater{
			{ID: Alice.ID, Weight: 33},
			{ID: Bob.ID, Weight: 67},
		})
		contractUpdateTx1 := ContractUpdateTx(Alice.Address, 1, types.Addr("0x1111111111111111111111111111111111111111"), 1, "81229384bdbe2ecc3572cb6deda168f2f0723161bc872bc53a9ef641ca2d087f1cd4aad56ddc7690550bc64aee0de76820d476b374699d4bd71dad9e7b5be8d81b")
		proposalBytes1 := contractUpdateTx1.Proposal.Bytes()
		contractUpdateTx2 := ContractUpdateTx(Bob.ID, 1, types.Addr("0x2222222222222222222222222222222222222222"), 1, "f0a8f180402795880004ab46ef99e5551782987e70f6a6899bc60c0b1c4a7bb56f0ef1db6a1c2f4b6f553332ce9cbf9856db250699822d8f9dbd03c915a90cc71c")
		proposalBytes2 := contractUpdateTx2.Proposal.Bytes()
		Convey("If a contract update transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(contract.HasApprovedUpdate(state, Alice.ID, proposalBytes1), ShouldBeTrue)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					Convey("CheckTx again should return Duplicated", func() {
						r := contractUpdateTx1.CheckTx(state)
						So(r.Code, ShouldEqual, response.ContractUpdateDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := contractUpdateTx1.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.ContractUpdateDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
						})
					})
				})
			})
			Convey("RawTransferTx should return the same encoded raw tx", func() {
				rawTx := RawContractUpdateTx(Alice.Address, 1, types.Addr("0x1111111111111111111111111111111111111111"), 1, "81229384bdbe2ecc3572cb6deda168f2f0723161bc872bc53a9ef641ca2d087f1cd4aad56ddc7690550bc64aee0de76820d476b374699d4bd71dad9e7b5be8d81b")
				So(rawTx, ShouldResemble, EncodeTx(contractUpdateTx1))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *ContractUpdateTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, contractUpdateTx1)
				})
			})
		})
		Convey("If a contract update transaction from LikeChainID is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := contractUpdateTx2.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx2))
					r := contractUpdateTx2.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(contract.HasApprovedUpdate(state, Bob.ID, proposalBytes2), ShouldBeTrue)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
				})
			})
		})
		Convey("If a contract update transaction has invalid format", func() {
			contractUpdateTx1.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateInvalidFormat.Code)
				Convey("DeliverTx should return InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a contract update transaction has unregistered sender", func() {
			contractUpdateTx1.Proposer = Mallory.Address
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateSenderNotRegistered.Code)
				Convey("DeliverTx should return SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a contract update transaction has invalid signature", func() {
			contractUpdateTx1.Sig = &ContractUpdateJSONSignature{Sig("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a contract update transaction has signature not from sender", func() {
			contractUpdateTx1.Sig = &ContractUpdateJSONSignature{Sig("5e7f8caefa0aa08236bcb9ff9fa7682b23ca76374c352fad2d666001b4e70647260c9e7dd601dd6ebec7d2097f16c88fcb55b194c5ad5b127f37d069b5840f901b")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a contract update transaction has invalid nonce", func() {
			contractUpdateTx1.Nonce = 2
			contractUpdateTx1.Sig = &ContractUpdateJSONSignature{Sig("968877403c283e7b6f551c998c6a3e65a0491dc86234fb7ec80523c6491f5570208d32b457f0fa00fb1b30e3eb35a85b3d870189ebe885125cec55227eed97091c")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateInvalidNonce.Code)
				Convey("DeliverTx should return InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a valid contract update transaction has proposer who has over 2/3 weight", func() {
			Convey("The contract update proposal should be executed", func() {
				txHash := utils.HashRawTx(EncodeTx(contractUpdateTx2))
				r := contractUpdateTx2.DeliverTx(state, txHash)
				So(r.Code, ShouldEqual, response.Success.Code)
				So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
				So(contract.HasApprovedUpdate(state, Bob.ID, proposalBytes2), ShouldBeTrue)
				So(contract.GetUpdateExecution(state, contractUpdateTx2.Proposal.ContractIndex), ShouldResemble, &contractUpdateTx2.Proposal.ContractAddress)
			})
		})
		Convey("If a contract update transaction is on a contract index which already has another proposal executed", func() {
			txHash := utils.HashRawTx(EncodeTx(contractUpdateTx2))
			contractUpdateTx2.DeliverTx(state, txHash)
			Convey("CheckTx should return InvalidIndex", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateInvalidIndex.Code)
				Convey("DeliverTx should return InvalidIndex", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateInvalidIndex.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a contract update transaction is on a contract index which is not the current index", func() {
			contractUpdateTx1.Proposal.ContractIndex = 2
			contractUpdateTx1.Sig = &ContractUpdateJSONSignature{Sig("985c46374d37270fda653bb20bc782d9559a503a5cd5df4cdac0fac3e76fab386028c3bb274f30d78f3d026f94b4638f87e542fe1ed53a1858cce6c6b39d6a2c1c")}
			Convey("CheckTx should return InvalidIndex", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateInvalidIndex.Code)
				Convey("DeliverTx should return InvalidIndex", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateInvalidIndex.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a contract update transaction has proposer who has already proposed another proposal", func() {
			txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
			contractUpdateTx1.DeliverTx(state, txHash)
			contractUpdateTx2.Proposer = Alice.ID
			contractUpdateTx2.Nonce = 2
			contractUpdateTx2.Sig = &ContractUpdateJSONSignature{Sig("098c7c4bd40d04c13367e250f7435f51cd7a2f6c71552561bebc259bb8c18c15643c679620d23d45dd187e12cbd6162c8bd4bc7ab67857036eccad3e6d0c08291c")}
			Convey("CheckTx should succeed", func() {
				r := contractUpdateTx2.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx2))
					r := contractUpdateTx2.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
				})
			})
		})
		Convey("If a contract update transaction has proposer who has already proposed the same proposal", func() {
			txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
			contractUpdateTx1.DeliverTx(state, txHash)
			contractUpdateTx2.Proposer = Alice.ID
			contractUpdateTx2.Proposal = contractUpdateTx1.Proposal
			contractUpdateTx2.Nonce = 2
			contractUpdateTx2.Sig = &ContractUpdateJSONSignature{Sig("1d1bf94ad7d95088cd39d55893abbfe3a405d627786bcc3417fde4a579174fda205f6ea25754baf65c06e10960d3af2365c5d21397a331dc6bcd42c370b441431c")}
			Convey("CheckTx should return ContractUpdateDoubleApproval", func() {
				r := contractUpdateTx2.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateDoubleApproval.Code)
				Convey("DeliverTx should return ContractUpdateDoubleApproval", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx2))
					r := contractUpdateTx2.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateDoubleApproval.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a contract update transaction has proposer who is not a contract updater", func() {
			contractUpdateTx1.Proposer = Carol.ID
			contractUpdateTx1.Sig = &ContractUpdateJSONSignature{Sig("78c818991872f28407560caf7bac6eee277f48f54105769480f97eb5865f5ea776e4225a85503ee90fe7a27b17e0dee07fac949ac01533863d09603bf60db2ee1b")}
			Convey("CheckTx should return NotUpdater", func() {
				r := contractUpdateTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.ContractUpdateNotUpdater.Code)
				Convey("DeliverTx should return NotUpdater", func() {
					txHash := utils.HashRawTx(EncodeTx(contractUpdateTx1))
					r := contractUpdateTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ContractUpdateNotUpdater.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
	})
}
