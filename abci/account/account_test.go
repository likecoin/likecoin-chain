package account

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateLikeChainID(t *testing.T) {
	ctx := context.NewMock()

	Convey("Generated LikeChain ID should match the format and fulfil some criteria", t, func() {
		_ = generateLikeChainID(ctx)

		// TODO: Checking
	})
}

func TestNewAccount(t *testing.T) {
	ctx := context.NewMock()
	Convey("Given a valid Ethereum address", t, func() {
		ctx.Reset()
		addr := common.HexToAddress("0x0123456789012345678901234567890123456789")

		Convey("An account is created with the address", func() {
			_, err := NewAccount(ctx, addr)

			So(err, ShouldBeFalse)
			// TODO: Check the identifier
		})
	})

	Convey("Given an invalid Ethereum address", t, func() {
		ctx.Reset()
		addr := common.HexToAddress("")

		Convey("Error is returned when creating account with the address", func() {
			_, err := NewAccount(ctx, addr)

			So(err, ShouldBeTrue)
		})
	})
}

func TestSaveAndFetchBalance(t *testing.T) {
	ctx := context.NewMock()
	// TODO: setup ctx
	id := types.LikeChainID{Content: nil} // TODO

	Convey("Given a valid balance", t, func() {
		balance, _ := new(big.Int).SetString("1000000000000000000", 10)

		Convey("The balance can be sucessfully saved to DB", func() {
			So(SaveBalance(ctx, id, balance), ShouldBeTrue)

			Convey("The balance can be sucessfully fetched from DB", func() {
				fetchedBalance := FetchBalance(ctx, id)
				So(fetchedBalance, ShouldEqual, balance)
			})
		})
	})
}

func TestSaveAndFetchNextNounce(t *testing.T) {
	ctx := context.NewMock()
	id := types.LikeChainID{Content: nil} // TODO

	Convey("For a newly created account", t, func() {
		ctx.Reset()
		// TODO: new account
		Convey("The nonce should be 0", func() {
			So(FetchNextNonce(ctx, id), ShouldEqual, 0)
			Convey("After incrementing nonce", func() {
				IncrementNextNonce(ctx, id)
				Convey("It should become 1", func() {
					So(FetchNextNonce(ctx, id), ShouldEqual, 1)
				})
			})
		})
	})

	Convey("For any current nonce value", t, func() {
		ctx.Reset()
		// TODO: setup account, randomize nonce value
		nonce := FetchNextNonce(ctx, id)
		Convey("After incrementing nonce", func() {
			IncrementNextNonce(ctx, id)
			Convey("It should increase by 1", func() {
				So(FetchNextNonce(ctx, id), ShouldEqual, nonce+1)
			})
		})
	})
}
