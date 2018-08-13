package account

import (
	"testing"

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
	Convey("Given a valid Ethereum address", t, func() {
		addr := []byte("0x0123456789012345678901234567890123456789")

		Convey("An account is created with the address", func() {
			_, err := NewAccount(addr)

			So(err, ShouldBeFalse)
			// TODO: Check the identifier
		})
	})

	Convey("Given an invalid Ethereum address", t, func() {
		addr := []byte("")

		Convey("Error is returned when creating account with the address", func() {
			_, err := NewAccount(addr)

			So(err, ShouldBeTrue)
		})
	})
}

func TestSaveAndFetchBalance(t *testing.T) {
	ctx := context.NewMock()
	identifier := &types.Identifier{} // TODO

	Convey("Given a valid balance", t, func() {
		balance := types.BigInteger{Content: []byte("1000000000000000000")}

		Convey("The balance can be sucessfully saved to DB", func() {
			So(SaveBalance(ctx, identifier, balance), ShouldBeTrue)

			Convey("The balance can be sucessfully fetched from DB", func() {
				fetchedBalance := FetchBalance(ctx, identifier)
				So(fetchedBalance, ShouldEqual, balance)
			})
		})
	})
}

func TestSaveAndFetchNextNounce(t *testing.T) {
	ctx := context.NewMock()
	identifier := &types.Identifier{} // TODO

	Convey("Given a valid nonce", t, func() {
		nextNonce := uint64(2)

		Convey("The nonce can be sucessfully saved to DB", func() {
			So(SaveNextNonce(ctx, identifier, nextNonce), ShouldBeTrue)

			Convey("The nonce can be sucessfully fetched from DB", func() {
				fetchedNonce := FetchNextNonce(ctx, identifier)
				So(fetchedNonce, ShouldEqual, nextNonce)
			})
		})
	})

	Convey("Given a smaller nonce", t, func() {
		nextNonce := uint64(1)

		Convey("The nonce can not be saved to DB", func() {
			So(SaveNextNonce(ctx, identifier, nextNonce), ShouldBeFalse)
		})
	})
}
