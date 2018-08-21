package account

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateLikeChainID(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given there is no LikeChain ID has generated before", t, func() {
		Convey("The seed of LikeChain ID should not exist in state tree", func() {
			_, seed := state.ImmutableStateTree().Get(likeChainIDSeedKey)
			So(seed, ShouldBeNil)
		})

		Convey("After generating the first LikeChain ID", func() {
			likeChainID1 := generateLikeChainID(state)

			Convey("The length of generated LikeChain ID should be 20", func() {
				So(len(likeChainID1.Content), ShouldEqual, 20)
			})

			Convey("The seed of LikeChain ID should exist in state tree", func() {
				_, seed1 := state.ImmutableStateTree().Get(likeChainIDSeedKey)
				So(seed1, ShouldNotBeNil)

				Convey("After generating the second LikeChain ID", func() {
					likeChainID2 := generateLikeChainID(state)

					Convey("The seed of LikeChain ID should be difference", func() {
						_, seed2 := state.ImmutableStateTree().Get(likeChainIDSeedKey)
						So(seed2, ShouldNotEqual, seed1)

						Convey("The generated LikeChain ID should be difference", func() {
							So(likeChainID2.Content, ShouldNotEqual, likeChainID1.Content)
						})
					})
				})
			})
		})
	})
}

func TestNewAccount(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a valid Ethereum address", t, func() {
		addr := common.HexToAddress("")

		Convey("An account is created with the address", func() {
			_, err := NewAccount(state, addr)

			So(err, ShouldBeNil)
		})
	})
}

func TestSaveAndFetchBalance(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()
	id := generateLikeChainID(state)

	Convey("Given a valid balance", t, func() {
		balance, _ := new(big.Int).SetString("1000000000000000000", 10)

		Convey("The balance can be sucessfully saved to DB", func() {
			So(SaveBalance(state, id, balance), ShouldBeNil)

			Convey("The balance can be sucessfully fetched from DB", func() {
				fetchedBalance := FetchBalance(state, id)
				So(fetchedBalance.Cmp(balance), ShouldEqual, 0)
			})
		})
	})
}

func TestIncrementAndFetchNextNounce(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("For a newly created account", t, func() {
		id, err := NewAccount(state, common.HexToAddress(""))
		if err != nil {
			panic("Unable to create new account")
		}

		Convey("The nonce should be 1", func() {
			So(FetchNextNonce(state, id), ShouldEqual, 1)

			Convey("After incrementing nonce", func() {
				IncrementNextNonce(state, id)

				Convey("It should become 2", func() {
					So(FetchNextNonce(state, id), ShouldEqual, 2)
				})
			})
		})
	})
}
