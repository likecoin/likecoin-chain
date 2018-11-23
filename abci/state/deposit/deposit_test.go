package deposit

import (
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDepositApproverInfo(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		approvers := []Approver{
			{Alice.ID, 10},
			{Bob.ID, 20},
		}
		Convey("GetDepositApproverInfo should be nil", func() {
			approverInfo := GetDepositApproverInfo(state, Alice.ID)
			So(approverInfo, ShouldBeNil)
			Convey("After setting deposit approvers", func() {
				SetDepositApprovers(state, approvers)
				Convey("GetDepositApproverInfo should return the approver's info", func() {
					approverInfo := GetDepositApproverInfo(state, Alice.ID)
					So(approverInfo, ShouldResemble, &approvers[0])
					approverInfo = GetDepositApproverInfo(state, Bob.ID)
					So(approverInfo, ShouldResemble, &approvers[1])
					Convey("GetDepositApproverInfo for non-existing approver should return nil", func() {
						approverInfo := GetDepositApproverInfo(state, Carol.ID)
						So(approverInfo, ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestCheckDeposit(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		approvers := []Approver{
			{Alice.ID, 10},
			{Bob.ID, 20},
			{Carol.ID, 30},
		}
		SetDepositApprovers(state, approvers)
		proposal1 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*Alice.Address, types.NewBigInt(100)},
				{*Bob.Address, types.NewBigInt(200)},
			},
		}
		proposal2 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*Carol.Address, types.NewBigInt(100)},
				{*Mallory.Address, types.NewBigInt(200)},
			},
		}
		Convey("For a valid DepositProposal", func() {
			Convey("CheckDeposit should return Success", func() {
				r := CheckDeposit(state, proposal1, Alice.ID)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a DepositProposal with proposer who is not a DepositApprover", func() {
			Convey("CheckDeposit should return DepositNotApprover", func() {
				r := CheckDeposit(state, proposal1, Mallory.ID)
				So(r, ShouldResemble, response.DepositNotApprover)
			})
		})
		Convey("For a DepositProposal with proposer who had already proposed another proposal on the same block number", func() {
			ProcessDeposit(state, proposal1, Alice.ID)
			Convey("CheckDeposit should succeed", func() {
				r := CheckDeposit(state, proposal2, Alice.ID)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a DepositProposal which is already proposed by the same proposer", func() {
			ProcessDeposit(state, proposal1, Alice.ID)
			Convey("CheckDeposit should return DepositDoubleApproval", func() {
				r := CheckDeposit(state, proposal1, Alice.ID)
				So(r, ShouldResemble, response.DepositDoubleApproval)
			})
		})
		Convey("For a DepositProposal which the block number has another ongoing proposal", func() {
			ProcessDeposit(state, proposal1, Alice.ID)
			Convey("CheckDeposit should return Success", func() {
				r := CheckDeposit(state, proposal2, Bob.ID)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a DepositProposal which the block number has another executed proposal", func() {
			ProcessDeposit(state, proposal1, Bob.ID)
			ProcessDeposit(state, proposal1, Carol.ID)
			Convey("CheckDeposit should return DepositAlreadyExecuted", func() {
				r := CheckDeposit(state, proposal2, Alice.ID)
				So(r, ShouldResemble, response.DepositAlreadyExecuted)
			})
		})
	})
}

func TestProcessDeposit(t *testing.T) {
	Convey("Given a valid DepositProposal", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		proposal := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*Alice.Address, types.NewBigInt(100)},
				{*Bob.Address, types.NewBigInt(200)},
				{*Carol.Address, types.NewBigInt(300)},
			},
		}
		Convey("GetDepositProposalWeight should return 0 before proposing the proposal", func() {
			proposalHash := proposal.Hash()
			queriedWeight := GetDepositProposalWeight(state, proposalHash)
			So(queriedWeight, ShouldEqual, 0)
			Convey("ProcessDeposit should return false when proposer's weight is not enough to execute the proposal", func() {
				approvers := []Approver{
					{Alice.ID, 33},
					{Bob.ID, 34},
					{Carol.ID, 33},
				}
				SetDepositApprovers(state, approvers)
				account.NewAccountFromID(state, Carol.ID, Carol.Address)
				executed := ProcessDeposit(state, proposal, Alice.ID)
				So(executed, ShouldBeFalse)
				Convey("HasApprovedDeposit should return true", func() {
					So(HasApprovedDeposit(state, Alice.ID, proposalHash), ShouldBeTrue)
					Convey("GetDepositProposalWeight should return proposer's weight", func() {
						queriedWeight := GetDepositProposalWeight(state, proposalHash)
						So(queriedWeight, ShouldEqual, approvers[0].Weight)
						Convey("When someone further propose the same proposal, making the total weight >2/3", func() {
							Convey("ProcessDeposit should return true", func() {
								executed := ProcessDeposit(state, proposal, Bob.ID)
								So(executed, ShouldBeTrue)
								Convey("Account balance should change accordingly", func() {
									balance := account.FetchBalance(state, Alice.Address)
									So(balance.String(), ShouldResemble, "100")
									balance = account.FetchBalance(state, Bob.Address)
									So(balance.String(), ShouldResemble, "200")
									balance = account.FetchBalance(state, Carol.ID)
									So(balance.String(), ShouldResemble, "300")
								})
							})
						})
					})
				})
			})
		})
		Convey("ProcessDeposit should return true when proposer's weight is enough to execute the proposal", func() {
			approvers := []Approver{
				{Alice.ID, 67},
				{Bob.ID, 33},
			}
			SetDepositApprovers(state, approvers)
			account.NewAccountFromID(state, Carol.ID, Carol.Address)
			Convey("ProcessDeposit should return true", func() {
				executed := ProcessDeposit(state, proposal, Alice.ID)
				So(executed, ShouldBeTrue)
				Convey("Account balance should change accordingly", func() {
					balance := account.FetchBalance(state, Alice.Address)
					So(balance.String(), ShouldResemble, "100")
					balance = account.FetchBalance(state, Bob.Address)
					So(balance.String(), ShouldResemble, "200")
					balance = account.FetchBalance(state, Carol.ID)
					So(balance.String(), ShouldResemble, "300")
				})
			})
		})
	})
}
