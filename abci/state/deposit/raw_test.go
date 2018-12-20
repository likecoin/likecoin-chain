package deposit

import (
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateInput(t *testing.T) {
	Convey("Given a deposit input", t, func() {
		input := Input{
			FromAddr: *Alice.Address,
			Value:    types.NewBigInt(1),
		}
		Convey("If the input is valid", func() {
			Convey("Validate should return true", func() {
				So(input.Validate(), ShouldBeTrue)
			})
		})
		Convey("If the input has nil value", func() {
			input.Value.Int = nil
			Convey("Validate should return false", func() {
				So(input.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the input has value less than 0", func() {
			input.Value = types.NewBigInt(-1)
			Convey("Validate should return false", func() {
				So(input.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the input has value 0", func() {
			input.Value = types.NewBigInt(0)
			Convey("Validate should return true", func() {
				So(input.Validate(), ShouldBeTrue)
			})
		})
		Convey("If the input has value 2^256-1", func() {
			n := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			n.Sub(n, big.NewInt(1))
			input.Value.Int = n
			Convey("Validate should return true", func() {
				So(input.Validate(), ShouldBeTrue)
			})
		})
		Convey("If the input has value 2^256", func() {
			input.Value.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			Convey("Validate should return false", func() {
				So(input.Validate(), ShouldBeFalse)
			})
		})
	})
}

func TestValidateProposalInputs(t *testing.T) {
	Convey("Given a list of proposal inputs", t, func() {
		inputs := ProposalInputs{
			{
				FromAddr: *Alice.Address,
				Value:    types.NewBigInt(1),
			},
			{
				FromAddr: *Bob.Address,
				Value:    types.NewBigInt(2),
			},
		}
		Convey("If the list is valid", func() {
			Convey("Validate should return true", func() {
				So(inputs.Validate(), ShouldBeTrue)
			})
		})
		Convey("If there is invalid input in the list", func() {
			inputs[0].Value.Int = nil
			Convey("Validate should return false", func() {
				So(inputs.Validate(), ShouldBeFalse)
			})
		})
	})
}

func TestSortProposalInputs(t *testing.T) {
	Convey("Given a list of proposal inputs", t, func() {
		inputs := ProposalInputs{
			{
				FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
				Value:    types.NewBigInt(100),
			},
			{
				FromAddr: *types.Addr("0x0000000000000000000000000000000000000000"),
				Value:    types.NewBigInt(100),
			},
			{
				FromAddr: *types.Addr("0x0000000000000000000000000000000000000002"),
				Value:    types.NewBigInt(1),
			},
			{
				FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
				Value:    types.NewBigInt(200),
			},
			{
				FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
				Value:    types.NewBigInt(100),
			},
		}
		Convey("Len should return the correct length", func() {
			So(inputs.Len(), ShouldEqual, 5)
		})
		Convey("Less should return the correct comparison result", func() {
			So(inputs.Less(0, 0), ShouldBeFalse)
			So(inputs.Less(0, 1), ShouldBeFalse)
			So(inputs.Less(0, 2), ShouldBeTrue)
			So(inputs.Less(0, 3), ShouldBeTrue)
			So(inputs.Less(0, 4), ShouldBeFalse)
			So(inputs.Less(1, 0), ShouldBeTrue)
			So(inputs.Less(1, 1), ShouldBeFalse)
			So(inputs.Less(1, 2), ShouldBeTrue)
			So(inputs.Less(1, 3), ShouldBeTrue)
			So(inputs.Less(1, 4), ShouldBeTrue)
			So(inputs.Less(2, 0), ShouldBeFalse)
			So(inputs.Less(2, 1), ShouldBeFalse)
			So(inputs.Less(2, 2), ShouldBeFalse)
			So(inputs.Less(2, 3), ShouldBeFalse)
			So(inputs.Less(2, 4), ShouldBeFalse)
			So(inputs.Less(3, 0), ShouldBeFalse)
			So(inputs.Less(3, 1), ShouldBeFalse)
			So(inputs.Less(3, 2), ShouldBeTrue)
			So(inputs.Less(3, 3), ShouldBeFalse)
			So(inputs.Less(3, 4), ShouldBeFalse)
			So(inputs.Less(4, 0), ShouldBeFalse)
			So(inputs.Less(4, 1), ShouldBeFalse)
			So(inputs.Less(4, 2), ShouldBeTrue)
			So(inputs.Less(4, 3), ShouldBeTrue)
			So(inputs.Less(4, 4), ShouldBeFalse)
		})
		Convey("Swap should swap the corresponding entries", func() {
			length := len(inputs)
			for i := 0; i < length; i++ {
				for j := 0; j < length; j++ {
					originInputs := make(ProposalInputs, length)
					copy(originInputs, inputs)
					inputs.Swap(i, j)
					So(inputs[i], ShouldResemble, originInputs[j])
					So(inputs[j], ShouldResemble, originInputs[i])
					for k := 0; k < length; k++ {
						if k == i || k == j {
							continue
						}
						So(inputs[k], ShouldResemble, originInputs[k])
					}
				}
			}
		})
		Convey("Sort should sort the list", func() {
			originInputs := make(ProposalInputs, len(inputs))
			copy(originInputs, inputs)
			proposal := Proposal{
				BlockNumber: 1337,
				Inputs:      inputs,
			}
			proposal.Sort()
			So(proposal.Inputs[0], ShouldResemble, originInputs[1])
			So(proposal.Inputs[1], ShouldResemble, originInputs[0])
			So(proposal.Inputs[2], ShouldResemble, originInputs[4])
			So(proposal.Inputs[3], ShouldResemble, originInputs[3])
			So(proposal.Inputs[4], ShouldResemble, originInputs[2])
		})
	})
}

func TestHashProposal(t *testing.T) {
	Convey("Given 2 proposals with same input content but different orders", t, func() {
		proposal1 := Proposal{
			BlockNumber: 1337,
			Inputs: ProposalInputs{
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
					Value:    types.NewBigInt(100),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000000"),
					Value:    types.NewBigInt(100),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000002"),
					Value:    types.NewBigInt(1),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
					Value:    types.NewBigInt(200),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
					Value:    types.NewBigInt(100),
				},
			},
		}
		proposal2 := Proposal{
			BlockNumber: 1337,
			Inputs: ProposalInputs{
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000000"),
					Value:    types.NewBigInt(100),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
					Value:    types.NewBigInt(100),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
					Value:    types.NewBigInt(100),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000001"),
					Value:    types.NewBigInt(200),
				},
				{
					FromAddr: *types.Addr("0x0000000000000000000000000000000000000002"),
					Value:    types.NewBigInt(1),
				},
			},
		}
		Convey("Hash should return the same hash value", func() {
			So(proposal1.Hash(), ShouldResemble, proposal2.Hash())
		})
	})
}

func TestValidateProposal(t *testing.T) {
	Convey("Given a proposal", t, func() {
		proposal := Proposal{
			BlockNumber: 1337,
			Inputs: ProposalInputs{
				{
					FromAddr: *Alice.Address,
					Value:    types.NewBigInt(1),
				},
				{
					FromAddr: *Bob.Address,
					Value:    types.NewBigInt(2),
				},
			},
		}
		Convey("If the proposal is valid", func() {
			Convey("Validate should return true", func() {
				So(proposal.Validate(), ShouldBeTrue)
			})
		})
		Convey("If the input list of the proposal is empty", func() {
			proposal.Inputs = ProposalInputs{}
			Convey("Validate should return false", func() {
				So(proposal.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the input list of the proposal is nil", func() {
			proposal.Inputs = nil
			Convey("Validate should return false", func() {
				So(proposal.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the input list of the proposal is invalid", func() {
			proposal.Inputs[0].Value.Int = nil
			Convey("Validate should return false", func() {
				So(proposal.Validate(), ShouldBeFalse)
			})
		})
	})
}

func TestDepositApprovers(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		Convey("DepositApprovers should be nil", func() {
			approvers := GetDepositApprovers(state)
			So(approvers, ShouldBeNil)
			Convey("DepositApproversWeightSum should be 0", func() {
				approversWeightSum := GetDepositApproversWeightSum(state)
				So(approversWeightSum, ShouldBeZeroValue)
				Convey("After setting deposit approvers", func() {
					approvers := []Approver{
						{Alice.ID, 10},
						{Bob.ID, 20},
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
									{Bob.ID, 30},
									{Carol.ID, 40},
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

func TestDepositApproval(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id := Alice.ID
		proposalHash := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
		Convey("HasApprovedDeposit should return false", func() {
			So(HasApprovedDeposit(state, id, proposalHash), ShouldBeFalse)
			Convey("After setting DepositApproval", func() {
				setDepositApproval(state, id, proposalHash)
				Convey("HasApprovedDeposit should return true", func() {
					So(HasApprovedDeposit(state, id, proposalHash), ShouldBeTrue)
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
		proposalHash := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
		Convey("GetDepositExecution should return nil", func() {
			executedTxHash := GetDepositExecution(state, blockNumber)
			So(executedTxHash, ShouldBeNil)
			Convey("After setting DepositExecution", func() {
				setDepositExecution(state, blockNumber, proposalHash)
				Convey("GetDepositExecution should return the set TxHash", func() {
					executedTxHash := GetDepositExecution(state, blockNumber)
					So(executedTxHash, ShouldResemble, proposalHash)
				})
			})
		})
	})
}
