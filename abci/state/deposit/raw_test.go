package deposit

import (
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDepositApprovers(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id0 := fixture.Alice.ID
		id1 := fixture.Bob.ID
		id2 := fixture.Carol.ID
		Convey("DepositApprovers should be nil", func() {
			approvers := GetDepositApprovers(state)
			So(approvers, ShouldBeNil)
			Convey("DepositApproversWeightSum should be 0", func() {
				approversWeightSum := GetDepositApproversWeightSum(state)
				So(approversWeightSum, ShouldBeZeroValue)
				Convey("After setting deposit approvers", func() {
					approvers := []Approver{
						{id0, 10},
						{id1, 20},
					}
					SetDepositApprovers(state, approvers)
					Convey("GetDepositApprovers should return the set DepositApprovers", func() {
						queriedApprovers := GetDepositApprovers(state)
						So(queriedApprovers, ShouldResemble, approvers)
						Convey("GetDepositApproversWeightSum should return the sum of the weights of the set DepositApprovers", func() {
							approversWeightSum := uint64(0)
							for _, approver := range approvers {
								approversWeightSum += uint64(approver.Weight)
							}
							queriedApproversWeightSum := GetDepositApproversWeightSum(state)
							So(queriedApproversWeightSum, ShouldEqual, approversWeightSum)
							Convey("After setting another list of DepositApprovers", func() {
								approvers := []Approver{
									{id1, 30},
									{id2, 40},
								}
								SetDepositApprovers(state, approvers)
								Convey("GetDepositApprovers should return the newly set DepositApprovers", func() {
									queriedApprovers := GetDepositApprovers(state)
									So(queriedApprovers, ShouldResemble, approvers)
									Convey("GetDepositApproversWeightSum should return the sum of the weights of the newly set DepositApprovers", func() {
										approversWeightSum := uint64(0)
										for _, approver := range approvers {
											approversWeightSum += uint64(approver.Weight)
										}
										queriedApproversWeightSum := GetDepositApproversWeightSum(state)
										So(queriedApproversWeightSum, ShouldEqual, approversWeightSum)
										Convey("After setting an empty list as DepositApprovers", func() {
											SetDepositApprovers(state, []Approver{})
											Convey("GetDepositApprovers should return nil", func() {
												queriedApprovers := GetDepositApprovers(state)
												So(queriedApprovers, ShouldBeNil)
												Convey("GetDepositApproversWeightSum should return 0", func() {
													queriedApproversWeightSum := GetDepositApproversWeightSum(state)
													So(queriedApproversWeightSum, ShouldBeZeroValue)
												})
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})
}

func TestDepositProposal(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		proposal := Proposal{
			BlockNumber: 1337,
			Inputs: []Input{
				{*fixture.Alice.Address, types.NewBigInt(100)},
				{*fixture.Bob.Address, types.NewBigInt(200)},
			},
		}
		txHash := common.Hex2Bytes("0000000000000000000000000000000000000001")
		Convey("Getting DepositProposal from non-existing txHash should return nil", func() {
			queriedProposal := GetDepositProposal(state, txHash)
			So(queriedProposal, ShouldBeNil)
			Convey("Weight of non-existing DepositProposal should be 0", func() {
				weight := GetDepositProposalWeight(state, txHash)
				So(weight, ShouldBeZeroValue)
				Convey("After setting DepositProposal", func() {
					SetDepositProposal(state, txHash, proposal)
					Convey("GetDepositProposal should return the set DepositProposal", func() {
						queriedProposal := GetDepositProposal(state, txHash)
						So(queriedProposal, ShouldResemble, &proposal)
						Convey("After increasing DepositProposal weight", func() {
							IncreaseDepositProposalWeight(state, txHash, 10)
							IncreaseDepositProposalWeight(state, txHash, 20)
							Convey("GetDepositProposalWeight should return the increased weight", func() {
								weight := GetDepositProposalWeight(state, txHash)
								So(weight, ShouldEqual, 30)
							})
						})
					})
				})
			})
		})
	})
}

func TestDepositApproval(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id := fixture.Alice.ID
		blockNumber := uint64(1337)
		txHash := common.Hex2Bytes("0000000000000000000000000000000000000000")
		Convey("GetDepositApproval should return nil", func() {
			approvalTxHash := GetDepositApproval(state, id, blockNumber)
			So(approvalTxHash, ShouldBeNil)
			Convey("After setting DepositApproval", func() {
				SetDepositApproval(state, id, blockNumber, txHash)
				Convey("GetDepositApproval should return the set TxHash", func() {
					approvalTxHash := GetDepositApproval(state, id, blockNumber)
					So(approvalTxHash, ShouldResemble, txHash)
				})
			})
		})
	})
}

func TestDepositExecution(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		blockNumber := uint64(1337)
		txHash := common.Hex2Bytes("0000000000000000000000000000000000000000")
		Convey("GetDepositExecution should return nil", func() {
			executedTxHash := GetDepositExecution(state, blockNumber)
			So(executedTxHash, ShouldBeNil)
			Convey("After setting DepositExecution", func() {
				SetDepositExecution(state, blockNumber, txHash)
				Convey("GetDepositExecution should return the set TxHash", func() {
					executedTxHash := GetDepositExecution(state, blockNumber)
					So(executedTxHash, ShouldResemble, txHash)
				})
			})
		})
	})
}
