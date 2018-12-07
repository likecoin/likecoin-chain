package contract

import (
	"testing"

	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProposalBytes(t *testing.T) {
	Convey("Proposals should have unique bytes", t, func() {
		proposal := Proposal{
			ContractIndex:   0xDEADBEEF1337C0DE,
			ContractAddress: *types.Addr("0x1234567890123456789012345678901234567890"),
		}
		So(proposal.Bytes(), ShouldResemble, common.Hex2Bytes("deadbeef1337c0de1234567890123456789012345678901234567890"))
	})
}

func TestContractUpdaters(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		Convey("ContractUpdaters should be nil", func() {
			updaters := GetContractUpdaters(state)
			So(updaters, ShouldBeNil)
			Convey("ContractUpdatersWeightSum should be 0", func() {
				updatersWeightSum := GetContractUpdatersWeightSum(state)
				So(updatersWeightSum, ShouldBeZeroValue)
				Convey("After setting contract updaters", func() {
					updaters := []Updater{
						{Alice.ID, 10},
						{Bob.ID, 20},
					}
					SetContractUpdaters(state, updaters)
					Convey("GetContractUpdaters should return the set contract updaters", func() {
						queriedUpdaters := GetContractUpdaters(state)
						So(queriedUpdaters, ShouldResemble, updaters)
						Convey("GetContractUpdatersWeightSum should return the sum of the weights of the set contract updaters", func() {
							updatersWeightSum := uint64(0)
							for _, updater := range updaters {
								updatersWeightSum += uint64(updater.Weight)
							}
							queriedUpdatersWeightSum := GetContractUpdatersWeightSum(state)
							So(queriedUpdatersWeightSum, ShouldEqual, updatersWeightSum)
							Convey("After setting another list of contract updaters", func() {
								updaters := []Updater{
									{Bob.ID, 30},
									{Carol.ID, 40},
								}
								SetContractUpdaters(state, updaters)
								Convey("GetContractUpdaters should return the newly set contract updaters", func() {
									queriedUpdaters := GetContractUpdaters(state)
									So(queriedUpdaters, ShouldResemble, updaters)
									Convey("GetContractUpdatersWeightSum should return the sum of the weights of the newly set contract updaters", func() {
										updatersWeightSum := uint64(0)
										for _, updater := range updaters {
											updatersWeightSum += uint64(updater.Weight)
										}
										queriedUpdatersWeightSum := GetContractUpdatersWeightSum(state)
										So(queriedUpdatersWeightSum, ShouldEqual, updatersWeightSum)
										Convey("After setting an empty list as contract updaters", func() {
											SetContractUpdaters(state, []Updater{})
											Convey("GetContractUpdaters should return nil", func() {
												queriedUpdaters := GetContractUpdaters(state)
												So(queriedUpdaters, ShouldBeNil)
												Convey("GetContractUpdatersWeightSum should return 0", func() {
													queriedUpdatersWeightSum := GetContractUpdatersWeightSum(state)
													So(queriedUpdatersWeightSum, ShouldBeZeroValue)
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

func TestUpdateProposalWeight(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		proposalBytes := common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000")
		Convey("GetUpdateProposalWeight should return 0", func() {
			weight := GetUpdateProposalWeight(state, proposalBytes)
			So(weight, ShouldBeZeroValue)
			Convey("After calling IncreaseUpdateProposalWeight", func() {
				IncreaseUpdateProposalWeight(state, proposalBytes, 123)
				Convey("GetUpdateProposalWeight should return the set weight", func() {
					weight := GetUpdateProposalWeight(state, proposalBytes)
					So(weight, ShouldEqual, 123)
					Convey("After further calling IncreaseUpdateProposalWeight", func() {
						IncreaseUpdateProposalWeight(state, proposalBytes, 234)
						Convey("GetUpdateProposalWeight should return the sum of the weights", func() {
							weight := GetUpdateProposalWeight(state, proposalBytes)
							So(weight, ShouldEqual, 123+234)
						})
					})
				})
			})
		})
	})
}

func TestUpdateApproval(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		id := Alice.ID
		proposalBytes := common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000")
		Convey("HasApprovedUpdate should return false", func() {
			HasApprovedUpdate(state, id, proposalBytes)
			So(HasApprovedUpdate(state, id, proposalBytes), ShouldBeFalse)
			Convey("After setting DepositApproval", func() {
				setUpdateApproval(state, id, proposalBytes)
				Convey("HasApprovedUpdate should return true", func() {
					So(HasApprovedUpdate(state, id, proposalBytes), ShouldBeTrue)
				})
			})
		})
	})
}

func TestUpdateExecution(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		proposal := Proposal{
			ContractIndex:   1337,
			ContractAddress: *types.Addr("0x1111111111111111111111111111111111111111"),
		}
		Convey("GetUpdateExecution should return nil", func() {
			executedContractAddr := GetUpdateExecution(state, proposal.ContractIndex)
			So(executedContractAddr, ShouldBeNil)
			Convey("GetUpdateExecutionWithProof should return nil", func() {
				_, version, _ := appCtx.GetMutableState().MutableWithdrawTree().SaveVersion()
				executedContractAddr, proof := GetUpdateExecutionWithProof(state, proposal.ContractIndex, version)
				So(executedContractAddr, ShouldBeNil)
				So(proof, ShouldBeNil)
				Convey("After setting UpdateExecution", func() {
					setUpdateExecution(state, &proposal)
					Convey("GetUpdateExecution should return the proposal address", func() {
						executedContractAddr := GetUpdateExecution(state, proposal.ContractIndex)
						So(executedContractAddr, ShouldResemble, &proposal.ContractAddress)
						Convey("GetUpdateExecutionWithProof should return valid proof", func() {
							_, version, _ := appCtx.GetMutableState().MutableWithdrawTree().SaveVersion()
							executedContractAddr, proof := GetUpdateExecutionWithProof(state, proposal.ContractIndex, version)
							So(executedContractAddr, ShouldResemble, &proposal.ContractAddress)
							So(proof, ShouldNotBeNil)
							So(proof.ComputeRootHash(), ShouldResemble, appCtx.GetMutableState().MutableWithdrawTree().Hash())
							Convey("GetUpdateExecutionWithProof on previous version should return nil", func() {
								executedContractAddr, proof := GetUpdateExecutionWithProof(state, proposal.ContractIndex, version-1)
								So(executedContractAddr, ShouldBeNil)
								So(proof, ShouldBeNil)
							})
						})
					})
				})
			})
		})
	})
}

func TestContractIndex(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		Convey("GetContractIndex should return 0", func() {
			So(GetContractIndex(state), ShouldBeZeroValue)
			Convey("After calling increaseContractIndex", func() {
				for i := 0; i < 1234; i++ {
					increaseContractIndex(state)
				}
				Convey("GetContractIndex should increase contract index correspondingly", func() {
					So(GetContractIndex(state), ShouldEqual, 1234)
				})
			})
		})
	})
}
