package contract

import (
	"testing"

	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDepositUpdaterInfo(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		updaters := []Updater{
			{Alice.ID, 10},
			{Bob.ID, 20},
		}
		Convey("GetContractUpdaterInfo should be nil", func() {
			updaterInfo := GetContractUpdaterInfo(state, Alice.ID)
			So(updaterInfo, ShouldBeNil)
			Convey("After setting contract updaters", func() {
				SetContractUpdaters(state, updaters)
				Convey("GetContractUpdaterInfo should return the updater's info", func() {
					updaterInfo := GetContractUpdaterInfo(state, Alice.ID)
					So(updaterInfo, ShouldResemble, &updaters[0])
					updaterInfo = GetContractUpdaterInfo(state, Bob.ID)
					So(updaterInfo, ShouldResemble, &updaters[1])
					Convey("GetContractUpdaterInfo for non-existing updater should return nil", func() {
						updaterInfo := GetContractUpdaterInfo(state, Carol.ID)
						So(updaterInfo, ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestCheckUpdate(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		updaters := []Updater{
			{Alice.ID, 10},
			{Bob.ID, 20},
			{Carol.ID, 30},
		}
		SetContractUpdaters(state, updaters)
		proposal1 := &Proposal{
			ContractIndex:   1,
			ContractAddress: *types.Addr("0x1111111111111111111111111111111111111111"),
		}
		proposal2 := &Proposal{
			ContractIndex:   1,
			ContractAddress: *types.Addr("0x2222222222222222222222222222222222222222"),
		}
		Convey("For a valid proposal", func() {
			Convey("CheckUpdate should return Success", func() {
				r := CheckUpdate(state, proposal1, Alice.ID)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a proposal with proposer who is not a contract updater", func() {
			Convey("CheckUpdate should return ContractUpdateNotUpdater", func() {
				r := CheckUpdate(state, proposal1, Mallory.ID)
				So(r, ShouldResemble, response.ContractUpdateNotUpdater)
			})
		})
		Convey("For a proposal with proposer who had already proposed another proposal on the same contract index", func() {
			ProcessUpdate(state, proposal1, Alice.ID)
			Convey("CheckUpdate should succeed", func() {
				r := CheckUpdate(state, proposal2, Alice.ID)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a proposal which is already proposed by the same proposer", func() {
			ProcessUpdate(state, proposal1, Alice.ID)
			Convey("CheckUpdate should return DepositDoubleApproval", func() {
				r := CheckUpdate(state, proposal1, Alice.ID)
				So(r, ShouldResemble, response.ContractUpdateDoubleApproval)
			})
		})
		Convey("For a proposal which the contract index has another ongoing proposal", func() {
			ProcessUpdate(state, proposal1, Alice.ID)
			Convey("CheckUpdate should return Success", func() {
				r := CheckUpdate(state, proposal2, Bob.ID)
				So(r, ShouldResemble, response.Success)
			})
		})
		Convey("For a proposal which the contract index is not the current index", func() {
			proposal1.ContractIndex = 2
			Convey("CheckUpdate should return ContractUpdateInvalidIndex", func() {
				r := CheckUpdate(state, proposal1, Alice.ID)
				So(r, ShouldResemble, response.ContractUpdateInvalidIndex)
			})
		})
		Convey("For a proposal which the contract index has another executed proposal", func() {
			ProcessUpdate(state, proposal1, Bob.ID)
			ProcessUpdate(state, proposal1, Carol.ID)
			Convey("CheckUpdate should return ContractUpdateInvalidIndex", func() {
				r := CheckUpdate(state, proposal2, Alice.ID)
				So(r, ShouldResemble, response.ContractUpdateInvalidIndex)
			})
		})
	})
}

func TestProcessDeposit(t *testing.T) {
	Convey("Given a valid proposal", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		proposal := &Proposal{
			ContractIndex:   1,
			ContractAddress: *types.Addr("0x1111111111111111111111111111111111111111"),
		}
		Convey("GetUpdateProposalWeight should return 0 before proposing the proposal", func() {
			proposalBytes := proposal.Bytes()
			queriedWeight := GetUpdateProposalWeight(state, proposalBytes)
			So(queriedWeight, ShouldEqual, 0)
			Convey("ProcessUpdate should return false when proposer's weight is not enough to execute the proposal", func() {
				updaters := []Updater{
					{Alice.ID, 33},
					{Bob.ID, 34},
					{Carol.ID, 33},
				}
				SetContractUpdaters(state, updaters)
				executed := ProcessUpdate(state, proposal, Alice.ID)
				So(executed, ShouldBeFalse)
				Convey("HasApprovedUpdate should return true", func() {
					So(HasApprovedUpdate(state, Alice.ID, proposalBytes), ShouldBeTrue)
					Convey("GetUpdateProposalWeight should return proposer's weight", func() {
						queriedWeight := GetUpdateProposalWeight(state, proposalBytes)
						So(queriedWeight, ShouldEqual, updaters[0].Weight)
						Convey("When someone further propose the same proposal, making the total weight >2/3", func() {
							Convey("ProcessUpdate should return true", func() {
								executed := ProcessUpdate(state, proposal, Bob.ID)
								So(executed, ShouldBeTrue)
								Convey("GetUpdateExecution should return the contract address", func() {
									So(GetUpdateExecution(state, proposal.ContractIndex), ShouldResemble, &proposal.ContractAddress)
									Convey("GetContractIndex should return increased value", func() {
										So(GetContractIndex(state), ShouldEqual, 1)
									})
								})
							})
						})
					})
				})
			})
		})
		Convey("ProcessUpdate should return true when proposer's weight is enough to execute the proposal", func() {
			updaters := []Updater{
				{Alice.ID, 67},
				{Bob.ID, 33},
			}
			SetContractUpdaters(state, updaters)
			Convey("ProcessUpdate should return true", func() {
				executed := ProcessUpdate(state, proposal, Alice.ID)
				So(executed, ShouldBeTrue)
				Convey("GetUpdateExecution should return the contract address", func() {
					So(GetUpdateExecution(state, proposal.ContractIndex), ShouldResemble, &proposal.ContractAddress)
					Convey("GetContractIndex should return increased value", func() {
						So(GetContractIndex(state), ShouldEqual, 1)
					})
				})
			})
		})
	})
}
