package htlc

import (
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashedTransferValidate(t *testing.T) {
	Convey("For a HashedTransfer", t, func() {
		ht := HashedTransfer{
			From:       fixture.Alice.Address,
			To:         fixture.Bob.ID,
			Value:      types.NewBigInt(10),
			HashCommit: [32]byte{},
			Expiry:     1,
		}
		Convey("If the HashedTransfer is valid", func() {
			Convey("Validate should return true", func() {
				So(ht.Validate(), ShouldBeTrue)
			})
		})
		Convey("If the HashedTransfer has nil From", func() {
			ht.From = nil
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the HashedTransfer has nil To", func() {
			ht.To = nil
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the HashedTransfer has nil value", func() {
			ht.Value.Int = nil
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the HashedTransfer has zero value", func() {
			ht.Value = types.NewBigInt(0)
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the HashedTransfer has negative value", func() {
			ht.Value = types.NewBigInt(-1)
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the HashedTransfer has value 2^256 - 1", func() {
			limit := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			limit.Sub(limit, big.NewInt(1))
			ht.Value = types.BigInt{Int: limit}
			Convey("Validate should return true", func() {
				So(ht.Validate(), ShouldBeTrue)
			})
		})
		Convey("If the HashedTransfer has value >= 2^256", func() {
			limit := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			ht.Value = types.BigInt{Int: limit}
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the HashedTransfer has negative expiry", func() {
			ht.Expiry = -1
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
		Convey("If the HashedTransfer has zero expiry", func() {
			ht.Expiry = 0
			Convey("Validate should return false", func() {
				So(ht.Validate(), ShouldBeFalse)
			})
		})
	})
}

func TestRawHashedTransfer(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		ht := HashedTransfer{
			From:       fixture.Alice.Address,
			To:         fixture.Bob.ID,
			Value:      types.NewBigInt(10),
			HashCommit: [32]byte{},
			Expiry:     1,
		}
		txHash := common.Hex2Bytes("0000000000000000000000000000000000000001")
		Convey("GetHashedTransfer should return nil", func() {
			So(GetHashedTransfer(state, txHash), ShouldBeNil)
			Convey("After CraeteHashedTransfer", func() {
				CreateHashedTransfer(state, &ht, txHash)
				Convey("GetHashedTransfer should return the HashedTransfer", func() {
					queriedHT := GetHashedTransfer(state, txHash)
					So(queriedHT, ShouldResemble, &ht)
					Convey("After RemoveHashedTransfer", func() {
						RemoveHashedTransfer(state, txHash)
						Convey("GetHashedTransfer should return nil", func() {
							So(GetHashedTransfer(state, txHash), ShouldBeNil)
						})
					})
				})
			})
		})
	})
}
