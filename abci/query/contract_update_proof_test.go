package query

import (
	"encoding/json"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/contract"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	"github.com/tendermint/iavl"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryContractUpdateProof(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		withdrawTree := state.MutableWithdrawTree()

		proposal1 := contract.Proposal{
			ContractIndex:   1,
			ContractAddress: *types.Addr("0x1111111111111111111111111111111111111111"),
		}
		contractIndexBytes1 := utils.EncodeUint64(proposal1.ContractIndex)
		key1 := []byte("exec")
		key1 = append(key1, contractIndexBytes1...)
		key1 = crypto.Sha256(key1)
		withdrawTree.Set(key1, proposal1.ContractAddress[:])
		_, version1, _ := withdrawTree.SaveVersion()

		height1 := int64(1)
		state.SetMetadataAtHeight(height1, context.TreeMetadata{
			WithdrawTreeVersion: version1,
		})

		query1 := abci.RequestQuery{
			Data:   contractIndexBytes1,
			Path:   "contract_update_proof",
			Height: height1,
		}

		proposal2 := contract.Proposal{
			ContractIndex:   2,
			ContractAddress: *types.Addr("0x2222222222222222222222222222222222222222"),
		}
		contractIndexBytes2 := utils.EncodeUint64(proposal2.ContractIndex)
		key2 := []byte("exec")
		key2 = append(key2, contractIndexBytes2...)
		key2 = crypto.Sha256(key2)
		withdrawTree.Set(key2, proposal1.ContractAddress[:])
		_, version2, _ := withdrawTree.SaveVersion()

		height2 := int64(2)
		state.SetMetadataAtHeight(height2, context.TreeMetadata{
			WithdrawTreeVersion: version2,
		})

		query2 := abci.RequestQuery{
			Data:   contractIndexBytes2,
			Path:   "contract_update_proof",
			Height: height2,
		}

		height3 := int64(123)
		state.SetMetadataAtHeight(123, context.TreeMetadata{
			WithdrawTreeVersion: 123,
		})

		Convey("If the contract_update_proof query is valid with previous height", func() {
			Convey("queryContractUpdateProof should succeed", func() {
				r := queryContractUpdateProof(state, query1)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("The response should contain a valid Merkle proof of the withdraw tree", func() {
					var res struct {
						ContractAddress types.Address   `json:"contract_address"`
						Proof           iavl.RangeProof `json:"proof"`
					}
					err := json.Unmarshal(r.Data, &res)
					So(err, ShouldBeNil)
					root := res.Proof.ComputeRootHash()
					oldTree, err := withdrawTree.GetImmutable(version1)
					So(err, ShouldBeNil)
					So(root, ShouldResemble, oldTree.Hash())
				})
			})
		})
		Convey("If the contract_update_proof query is valid with present height", func() {
			Convey("queryContractUpdateProof should succeed", func() {
				r := queryContractUpdateProof(state, query2)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("The response should contain a valid Merkle proof of the withdraw tree", func() {
					var res struct {
						ContractAddress types.Address   `json:"contract_address"`
						Proof           iavl.RangeProof `json:"proof"`
					}
					err := json.Unmarshal(r.Data, &res)
					So(err, ShouldBeNil)
					root := res.Proof.ComputeRootHash()
					oldTree, err := withdrawTree.GetImmutable(version2)
					So(err, ShouldBeNil)
					So(root, ShouldResemble, oldTree.Hash())
				})
			})
		})
		Convey("If the contract_update_proof query has invalid height", func() {
			query1.Height = 0
			Convey("queryContractUpdateProof should return InvalidHeight", func() {
				r := queryContractUpdateProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryContractUpdateProofInvalidHeight.Code)
			})
		})
		Convey("If the contract_update_proof query has non-existing height", func() {
			query1.Height = 3
			Convey("queryContractUpdateProof should return InvalidHeight", func() {
				r := queryContractUpdateProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryContractUpdateProofInvalidHeight.Code)
			})
		})
		Convey("If the contract_update_proof query has height with GC-ed tree version", func() {
			query1.Height = height3
			Convey("queryContractUpdateProof should return InvalidHeight", func() {
				r := queryContractUpdateProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryContractUpdateProofInvalidHeight.Code)
			})
		})
		Convey("If the contract_update_proof query has nil data", func() {
			query1.Data = nil
			Convey("queryContractUpdateProof should return NotExist", func() {
				r := queryContractUpdateProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryContractUpdateProofNotExist.Code)
			})
		})
		Convey("If the contract_update_proof query has invalid data", func() {
			query1.Data = make([]byte, 32)
			query1.Data[0] = 2
			Convey("queryContractUpdateProof should return NotExist", func() {
				r := queryContractUpdateProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryContractUpdateProofNotExist.Code)
			})
		})
	})
}
