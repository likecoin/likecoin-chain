package txs

import (
	"math/big"
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

func TestDepositValidateFormat(t *testing.T) {
	Convey("For a deposit transaction", t, func() {
		inputs := []deposit.Input{
			{
				FromAddr: *fixture.Alice.Address,
				Value:    types.NewBigInt(100),
			},
			{
				FromAddr: *fixture.Bob.Address,
				Value:    types.NewBigInt(200),
			},
		}
		depositTx := DepositTx(fixture.Alice.Address, 1, inputs, 1, "60409476e6dc6449434c33e3de4de8f62cbdc40ccbc909a7d538e3a2b4d656c952dc7b204115816d4b10b92e79fb4a50ff0309c3cd97a447f3c8198f542e8ae11c")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(depositTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has nil signature", func() {
			depositTx.Sig = nil
			Convey("The validation should fail", func() {
				So(depositTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil Proposer field", func() {
			depositTx.Proposer = nil
			Convey("The validation should fail", func() {
				So(depositTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has no inputs", func() {
			depositTx.Proposal.Inputs = nil
			Convey("The validation should fail", func() {
				So(depositTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil value input", func() {
			depositTx.Proposal.Inputs[0].Value.Int = nil
			Convey("The validation should fail", func() {
				So(depositTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has negative value input", func() {
			depositTx.Proposal.Inputs[0].Value.Int = big.NewInt(-1)
			Convey("The validation should fail", func() {
				So(depositTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestDepositSignature(t *testing.T) {
	Convey("If a deposit transaction is valid", t, func() {
		inputs := []deposit.Input{
			{
				FromAddr: *fixture.Alice.Address,
				Value:    types.NewBigInt(100),
			},
			{
				FromAddr: *fixture.Bob.Address,
				Value:    types.NewBigInt(200),
			},
		}
		depositTx := DepositTx(fixture.Alice.Address, 1, inputs, 1, "60409476e6dc6449434c33e3de4de8f62cbdc40ccbc909a7d538e3a2b4d656c952dc7b204115816d4b10b92e79fb4a50ff0309c3cd97a447f3c8198f542e8ae11c")
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := depositTx.Sig.RecoverAddress(depositTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the Proposer address of the deposit transaction", func() {
				So(recoveredAddr, ShouldResemble, depositTx.Proposer)
			})
		})
	})
}

func TestCheckAndDeliverDeposit(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID

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
		deposit.SetDepositApprovers(state, []deposit.Approver{
			{ID: id0, Weight: 33},
			{ID: id1, Weight: 67},
		})
		depositTx1 := DepositTx(fixture.Alice.Address, 1337, inputs1, 1, "af6109def6aebe496980b7cda873781fecf0c1d8f1b432b736a91f67b0372ee50d74c30d6718bdfcf0097aefeed5dbeebde3200030edbc7c064e29d76d5e60b11b")
		depositTx2 := DepositTx(fixture.Bob.ID, 1337, inputs2, 1, "7cb046891ca457d0dc7082b22d427fbf62616d4c23ffdddc1c10b8cd4b65260235db90afaf4d22e452e1591e66835a3b11cc4129fc8c8ed732779e9b1b68c3661c")
		Convey("If a deposit transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(deposit.GetDepositProposal(state, txHash), ShouldResemble, &depositTx1.Proposal)
					So(deposit.GetDepositApproval(state, id0, depositTx1.Proposal.BlockNumber), ShouldResemble, txHash)
					So(r.Status, ShouldEqual, txstatus.TxStatusPending)
					Convey("CheckTx again should return Duplicated", func() {
						r := depositTx1.CheckTx(state)
						So(r.Code, ShouldEqual, response.DepositDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := depositTx1.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.DepositDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
						})
					})
				})
			})
			Convey("RawTransferTx should return the same encoded raw tx", func() {
				rawTx := RawDepositTx(fixture.Alice.Address, 1337, inputs1, 1, "af6109def6aebe496980b7cda873781fecf0c1d8f1b432b736a91f67b0372ee50d74c30d6718bdfcf0097aefeed5dbeebde3200030edbc7c064e29d76d5e60b11b")
				So(rawTx, ShouldResemble, EncodeTx(depositTx1))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *DepositTransaction
					err := types.AminoCodec().UnmarshalBinary(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, depositTx1)
				})
			})
		})
		Convey("If a deposit transaction from LikeChainID is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := depositTx2.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx2))
					r := depositTx2.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(deposit.GetDepositProposal(state, txHash), ShouldResemble, &depositTx2.Proposal)
					So(deposit.GetDepositApproval(state, id1, depositTx2.Proposal.BlockNumber), ShouldResemble, txHash)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
				})
			})
		})
		Convey("If a deposit transaction has invalid format", func() {
			depositTx1.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositInvalidFormat.Code)
				Convey("DeliverTx should return InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit transaction has unregistered sender", func() {
			depositTx1.Proposer = fixture.Mallory.Address
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositSenderNotRegistered.Code)
				Convey("DeliverTx should return SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit transaction has invalid signature", func() {
			depositTx1.Sig = &DepositJSONSignature{Sig("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit transaction has signature not from sender", func() {
			depositTx1.Sig = &DepositJSONSignature{Sig("ae79fa211a4cec72255d87fd58436a9049f3343717e3cfb611a85c08e71444086be967022ab92075efb72c35b49d0ac37eac7e58fe1a832ba68be8c9305970ef1b")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit transaction has invalid nonce", func() {
			depositTx1.Nonce = 2
			depositTx1.Sig = &DepositJSONSignature{Sig("580e89ccc13628789a5854dafb7e419149cdd2a1c2ea7335448c98309cf334932dbc717cf0d5eac1cbac48c5f6c6fd466d93473370693864df7a05883f75e5dc1c")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositInvalidNonce.Code)
				Convey("DeliverTx should return InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a valid deposit transaction has proposer who has over 2/3 weight", func() {
			Convey("The deposit proposal should be executed", func() {
				txHash := utils.HashRawTx(EncodeTx(depositTx2))
				r := depositTx2.DeliverTx(state, txHash)
				So(r.Code, ShouldEqual, response.Success.Code)
				So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
				So(deposit.GetDepositProposal(state, txHash), ShouldResemble, &depositTx2.Proposal)
				So(deposit.GetDepositApproval(state, id1, depositTx2.Proposal.BlockNumber), ShouldResemble, txHash)
				So(deposit.GetDepositExecution(state, depositTx2.Proposal.BlockNumber), ShouldResemble, txHash)
				So(account.FetchBalance(state, id2).Cmp(depositTx2.Proposal.Inputs[0].Value.Int), ShouldBeZeroValue)
				So(account.FetchBalance(state, fixture.Mallory.Address).Cmp(depositTx2.Proposal.Inputs[1].Value.Int), ShouldBeZeroValue)
			})
		})
		Convey("If a deposit transaction is on some block number which already has another proposal executed", func() {
			txHash := utils.HashRawTx(EncodeTx(depositTx2))
			depositTx2.DeliverTx(state, txHash)
			Convey("CheckTx should return AlreadyExecuted", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositAlreadyExecuted.Code)
				Convey("DeliverTx should return AlreadyExecuted", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositAlreadyExecuted.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit transaction has proposer who has already proposed another proposal", func() {
			txHash := utils.HashRawTx(EncodeTx(depositTx1))
			depositTx1.DeliverTx(state, txHash)
			depositTx2.Proposer = fixture.Alice.ID
			depositTx2.Nonce = 2
			depositTx2.Sig = &DepositJSONSignature{Sig("4748562cf12388776d769cd0bfc93ff3cf64fd77c52d00f14bfb4c295f3b3ad76a45058242054844a7f4fce57ed3c5505a9022cd9e720801d9ee051f7be59f691b")}
			Convey("CheckTx should return DoubleApproval", func() {
				r := depositTx2.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositDoubleApproval.Code)
				Convey("DeliverTx should return DoubleApproval", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx2))
					r := depositTx2.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositDoubleApproval.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a deposit transaction has proposer who is not a deposit approver", func() {
			depositTx1.Proposer = fixture.Carol.ID
			depositTx1.Sig = &DepositJSONSignature{Sig("410962be0eedb26cdb299265c1c2b08b9565246a3dc51601534799ae956072e3619178d676ab90f8e6d327236d6a72897abb1aaabae60c901f2850f04c066d191c")}
			Convey("CheckTx should return NotApprover", func() {
				r := depositTx1.CheckTx(state)
				So(r.Code, ShouldEqual, response.DepositNotApprover.Code)
				Convey("DeliverTx should return NotApprover", func() {
					txHash := utils.HashRawTx(EncodeTx(depositTx1))
					r := depositTx1.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.DepositNotApprover.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
	})
}
