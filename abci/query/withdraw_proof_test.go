package query

import (
	"encoding/json"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/iavl"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

func TestQueryWithdrawProof(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		withdrawTree := state.MutableWithdrawTree()

		packedTx1 := make([]byte, 112)
		withdrawTree.Set(crypto.Sha256(packedTx1), []byte{1})
		_, version1, _ := withdrawTree.SaveVersion()

		height1 := int64(1)
		state.SetMetadataAtHeight(height1, context.TreeMetadata{
			WithdrawTreeVersion: version1,
		})

		query1 := abci.RequestQuery{
			Data:   packedTx1,
			Path:   "withdraw_proof",
			Height: height1,
		}

		packedTx2 := make([]byte, 112)
		packedTx2[0] = 1
		withdrawTree.Set(crypto.Sha256(packedTx2), []byte{1})
		_, version2, _ := withdrawTree.SaveVersion()

		height2 := int64(2)
		state.SetMetadataAtHeight(height2, context.TreeMetadata{
			WithdrawTreeVersion: version2,
		})

		query2 := abci.RequestQuery{
			Data:   packedTx2,
			Path:   "withdraw_proof",
			Height: height2,
		}

		height3 := int64(123)
		state.SetMetadataAtHeight(123, context.TreeMetadata{
			WithdrawTreeVersion: 123,
		})
		Convey("If the withdraw_proof query is valid with previous height", func() {
			Convey("queryWithdrawProof should succeed", func() {
				r := queryWithdrawProof(state, query1)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("The response should contain a valid Merkle proof of the withdraw tree", func() {
					var proof iavl.RangeProof
					err := json.Unmarshal(r.Data, &proof)
					So(err, ShouldBeNil)
					root := proof.ComputeRootHash()
					oldTree, err := withdrawTree.GetImmutable(version1)
					So(err, ShouldBeNil)
					So(root, ShouldResemble, oldTree.Hash())
				})
			})
		})
		Convey("If the withdraw_proof query is valid with present height", func() {
			Convey("queryWithdrawProof should succeed", func() {
				r := queryWithdrawProof(state, query2)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("The response should contain a valid Merkle proof of the withdraw tree", func() {
					var proof iavl.RangeProof
					err := json.Unmarshal(r.Data, &proof)
					So(err, ShouldBeNil)
					root := proof.ComputeRootHash()
					oldTree, err := withdrawTree.GetImmutable(version2)
					So(err, ShouldBeNil)
					So(root, ShouldResemble, oldTree.Hash())
				})
			})
		})
		Convey("If the withdraw_proof query has invalid height", func() {
			query1.Height = 0
			Convey("queryWithdrawProof should return InvalidHeight", func() {
				r := queryWithdrawProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryWithdrawProofInvalidHeight.Code)
			})
		})
		Convey("If the withdraw_proof query has non-existing height", func() {
			query1.Height = 3
			Convey("queryWithdrawProof should return InvalidHeight", func() {
				r := queryWithdrawProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryWithdrawProofInvalidHeight.Code)
			})
		})
		Convey("If the withdraw_proof query has height with GC-ed tree version", func() {
			query1.Height = height3
			Convey("queryWithdrawProof should return InvalidHeight", func() {
				r := queryWithdrawProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryWithdrawProofInvalidHeight.Code)
			})
		})
		Convey("If the withdraw_proof query has nil data", func() {
			query1.Data = nil
			Convey("queryWithdrawProof should return NotExist", func() {
				r := queryWithdrawProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryWithdrawProofNotExist.Code)
			})
		})
		Convey("If the withdraw_proof query has invalid data", func() {
			query1.Data = make([]byte, 32)
			query1.Data[0] = 2
			Convey("queryWithdrawProof should return NotExist", func() {
				r := queryWithdrawProof(state, query1)
				So(r.Code, ShouldEqual, response.QueryWithdrawProofNotExist.Code)
			})
		})
	})
}
