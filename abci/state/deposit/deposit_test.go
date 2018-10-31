package deposit

import (
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDepositApproverInfo(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		approvers := []Approver{
			{id0, 10},
			{id1, 20},
		}
		Convey("GetDepositApproverInfo should be nil", func() {
			approverInfo := GetDepositApproverInfo(state, id0)
			So(approverInfo, ShouldBeNil)
			Convey("After setting deposit approvers", func() {
				SetDepositApprovers(state, approvers)
				Convey("GetDepositApproverInfo should return the approver's info", func() {
					approverInfo := GetDepositApproverInfo(state, id0)
					So(approverInfo, ShouldResemble, &approvers[0])
					approverInfo = GetDepositApproverInfo(state, id1)
					So(approverInfo, ShouldResemble, &approvers[1])
					Convey("GetDepositApproverInfo for non-existing approver should return nil", func() {
						approverInfo := GetDepositApproverInfo(state, fixture.Carol.ID)
						So(approverInfo, ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestCheckDepositProposal(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID
		approvers := []Approver{
			{id0, 10},
			{id1, 20},
			{id2, 30},
		}
		SetDepositApprovers(state, approvers)
		proposal1 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Alice.Address, types.NewBigInt(100)},
				{*fixture.Bob.Address, types.NewBigInt(200)},
			},
		}
		txHash1 := common.Hex2Bytes("0000000000000000000000000000000000000001")
		proposal2 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Carol.Address, types.NewBigInt(100)},
				{*fixture.Mallory.Address, types.NewBigInt(200)},
			},
		}
		Convey("For a valid DepositProposal", func() {
			Convey("CheckDepositProposal should return Success", func() {
				r := CheckDepositProposal(state, proposal1, id0)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a DepositProposal with proposer who is not a DepositApprover", func() {
			Convey("CheckDepositProposal should return DepositNotApprover", func() {
				r := CheckDepositProposal(state, proposal1, fixture.Mallory.ID)
				So(r, ShouldResemble, response.DepositNotApprover)
			})
		})
		Convey("For a DepositProposal with proposer who had already proposed another proposal on the same block number", func() {
			CreateDepositProposal(state, txHash1, proposal1, id0)
			Convey("CheckDepositProposal should return DepositDoubleApproval", func() {
				r := CheckDepositProposal(state, proposal2, id0)
				So(r, ShouldResemble, response.DepositDoubleApproval)
			})
		})
		Convey("For a DepositProposal with proposer who had already approved another proposal", func() {
			CreateDepositProposal(state, txHash1, proposal1, id0)
			CreateDepositApproval(state, txHash1, id1)
			Convey("CheckDepositProposal should return DepositDoubleApproval", func() {
				r := CheckDepositProposal(state, proposal2, id1)
				So(r, ShouldResemble, response.DepositDoubleApproval)
			})
		})
		Convey("For a DepositProposal which the block number has another ongoing proposal", func() {
			CreateDepositProposal(state, txHash1, proposal1, id0)
			Convey("CheckDepositProposal should return Success", func() {
				r := CheckDepositProposal(state, proposal2, id1)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a DepositProposal which the block number has another executed proposal", func() {
			CreateDepositProposal(state, txHash1, proposal1, id0)
			ExecuteDepositProposal(state, txHash1)
			Convey("CheckDepositProposal should return DepositAlreadyExecuted", func() {
				r := CheckDepositProposal(state, proposal2, id1)
				So(r, ShouldResemble, response.DepositAlreadyExecuted)
			})
		})
	})
}

func TestCreateDepositProposal(t *testing.T) {
	Convey("Given a valid DepositProposal", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID
		approvers := []Approver{
			{id0, 10},
			{id1, 20},
			{id2, 30},
		}
		SetDepositApprovers(state, approvers)
		proposal1 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Alice.Address, types.NewBigInt(100)},
				{*fixture.Bob.Address, types.NewBigInt(200)},
			},
		}
		txHash1 := common.Hex2Bytes("0000000000000000000000000000000000000001")
		Convey("CreateDepositProposal should return proposer's weight", func() {
			weight := CreateDepositProposal(state, txHash1, proposal1, id0)
			So(weight, ShouldEqual, approvers[0].Weight)
			Convey("Should be able to get proposal by GetDepositProposal", func() {
				queriedProposal := GetDepositProposal(state, txHash1)
				So(queriedProposal, ShouldResemble, &proposal1)
				Convey("GetDepositProposalWeight should return proposer's weight", func() {
					queriedWeight := GetDepositProposalWeight(state, txHash1)
					So(queriedWeight, ShouldEqual, approvers[0].Weight)
					Convey("Should be able to get proposer's approval on the proposal's block number", func() {
						approvalTxHash := GetDepositApproval(state, id0, proposal1.BlockNumber)
						So(approvalTxHash, ShouldResemble, txHash1)
					})
				})
			})
		})
	})
}

func TestCheckDepositApproval(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID
		approvers := []Approver{
			{id0, 10},
			{id1, 20},
			{id2, 30},
		}
		SetDepositApprovers(state, approvers)
		proposal1 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Alice.Address, types.NewBigInt(100)},
				{*fixture.Bob.Address, types.NewBigInt(200)},
			},
		}
		txHash1 := common.Hex2Bytes("0000000000000000000000000000000000000000")
		proposal2 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Carol.Address, types.NewBigInt(100)},
				{*fixture.Mallory.Address, types.NewBigInt(200)},
			},
		}
		txHash2 := common.Hex2Bytes("0000000000000000000000000000000000000002")
		Convey("For a valid DepositApproval", func() {
			Convey("CheckDepositApproval should return DepositApprovalProposalNotExist", func() {
				r := CheckDepositApproval(state, txHash1, id1)
				So(r, ShouldResemble, response.DepositApprovalProposalNotExist)
				Convey("After proposing the corresponding DepositProposal", func() {
					CreateDepositProposal(state, txHash1, proposal1, id0)
					Convey("CheckDepositApproval should return Success", func() {
						r := CheckDepositApproval(state, txHash1, id1)
						So(r, ShouldResemble, response.Success)
						Convey("If the approver then try to approve another proposal with the same block number", func() {
							CreateDepositApproval(state, txHash1, id1)
							CreateDepositProposal(state, txHash2, proposal2, id2)
							Convey("CheckDepositApproval should return DepositApprovalDoubleApproval", func() {
								r := CheckDepositApproval(state, txHash2, id1)
								So(r, ShouldResemble, response.DepositApprovalDoubleApproval)
							})
						})
					})
				})
			})
		})
		Convey("For a DepositApproval with approver who is not a DepositApprover", func() {
			CreateDepositProposal(state, txHash1, proposal1, id0)
			Convey("CheckDepositApproval should return DepositApprovalNotApprover", func() {
				r := CheckDepositApproval(state, txHash1, fixture.Mallory.ID)
				So(r, ShouldResemble, response.DepositApprovalNotApprover)
			})
		})
		Convey("For a DepositApproval which the block number has another ongoing proposal", func() {
			CreateDepositProposal(state, txHash1, proposal1, id0)
			CreateDepositProposal(state, txHash2, proposal2, id1)
			Convey("CheckDepositApproval should return Success", func() {
				r := CheckDepositApproval(state, txHash2, id2)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a DepositApproval which the block number has another executed proposal", func() {
			CreateDepositProposal(state, txHash1, proposal1, id0)
			CreateDepositProposal(state, txHash2, proposal2, id1)
			ExecuteDepositProposal(state, txHash1)
			Convey("CheckDepositApproval should return DepositAlreadyExecuted", func() {
				r := CheckDepositApproval(state, txHash2, id2)
				So(r, ShouldResemble, response.DepositApprovalAlreadyExecuted)
			})
		})
	})
}

func TestCreateDepositApproval(t *testing.T) {
	Convey("Given a valid DepositApproval", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID
		approvers := []Approver{
			{id0, 10},
			{id1, 20},
			{id2, 30},
		}
		SetDepositApprovers(state, approvers)
		proposal1 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Alice.Address, types.NewBigInt(100)},
				{*fixture.Bob.Address, types.NewBigInt(200)},
			},
		}
		txHash1 := common.Hex2Bytes("0000000000000000000000000000000000000000")
		CreateDepositProposal(state, txHash1, proposal1, id0)
		Convey("CreateDepositApproval should return the proposal's new weight", func() {
			weight := CreateDepositApproval(state, txHash1, id2)
			So(weight, ShouldEqual, 40)
			Convey("Should be able to get proposer's approval on the proposal's block number", func() {
				queriedTxHash := GetDepositApproval(state, id2, proposal1.BlockNumber)
				So(queriedTxHash, ShouldResemble, txHash1)
			})
		})
	})
}

func TestExecuteDepositProposal(t *testing.T) {
	Convey("Given a state with ongoing DepositProposal", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID
		account.NewAccountFromID(state, id0, fixture.Alice.Address)
		approvers := []Approver{
			{id0, 10},
			{id1, 20},
			{id2, 30},
		}
		SetDepositApprovers(state, approvers)
		proposal1 := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Alice.Address, types.NewBigInt(100)},
				{*fixture.Bob.Address, types.NewBigInt(200)},
			},
		}
		txHash1 := common.Hex2Bytes("0000000000000000000000000000000000000000")
		CreateDepositProposal(state, txHash1, proposal1, id0)
		Convey("After ExecuteDepositProposal", func() {
			ExecuteDepositProposal(state, txHash1)
			Convey("Account balance should increase according the proposal", func() {
				balance := account.FetchBalance(state, id0)
				So(balance.Cmp(proposal1.Inputs[0].Value.Int), ShouldBeZeroValue)
				balance = account.FetchBalance(state, fixture.Bob.Address)
				So(balance.Cmp(proposal1.Inputs[1].Value.Int), ShouldBeZeroValue)
				Convey("GetDepositExecution for the block number should return the txHash", func() {
					queriedTxHash := GetDepositExecution(state, proposal1.BlockNumber)
					So(queriedTxHash, ShouldResemble, txHash1)
				})
			})
		})
	})
}
