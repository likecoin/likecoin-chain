package htlc

import (
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashedTransferIsExpired(t *testing.T) {
	Convey("Given an empty state and a valid HashedTransfer", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		ht := HashedTransfer{
			From:       fixture.Alice.Address,
			To:         fixture.Bob.ID,
			Value:      types.NewBigInt(10),
			HashCommit: [32]byte{},
			Expiry:     10,
		}
		Convey("If the expiry time of the HashedTransfer is less than the block time in context", func() {
			state.SetBlockTime(11)
			Convey("IsExpired should return true", func() {
				So(ht.IsExpired(state), ShouldBeTrue)
			})
		})
		Convey("If the expiry time of the HashedTransfer is equal to the block time in context", func() {
			state.SetBlockTime(10)
			Convey("IsExpired should return true", func() {
				So(ht.IsExpired(state), ShouldBeTrue)
			})
		})
		Convey("If the expiry time of the HashedTransfer is more than the block time in context", func() {
			state.SetBlockTime(9)
			Convey("IsExpired should return false", func() {
				So(ht.IsExpired(state), ShouldBeFalse)
			})
		})
	})
}

func TestCheckCreateHashedTransfer(t *testing.T) {
	Convey("Given an empty state and a valid HashedTransfer", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		ht := HashedTransfer{
			From:       fixture.Alice.Address,
			To:         fixture.Bob.ID,
			Value:      types.NewBigInt(10),
			HashCommit: [32]byte{},
			Expiry:     10,
		}
		Convey("If the expiry time of the HashedTransfer is less than the block time in context", func() {
			state.SetBlockTime(11)
			Convey("CheckCreateHashedTransfer should return InvalidExpiry", func() {
				r := CheckCreateHashedTransfer(state, &ht)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidExpiry.Code)
			})
		})
		Convey("If the expiry time of the HashedTransfer is equal to the block time in context", func() {
			state.SetBlockTime(10)
			Convey("CheckCreateHashedTransfer should return InvalidExpiry", func() {
				r := CheckCreateHashedTransfer(state, &ht)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidExpiry.Code)
			})
		})
		Convey("If the expiry time of the HashedTransfer is more than the block time in context", func() {
			state.SetBlockTime(9)
			Convey("CheckCreateHashedTransfer should return Success", func() {
				r := CheckCreateHashedTransfer(state, &ht)
				So(r.Code, ShouldEqual, response.Success.Code)
			})
		})
	})
}

func TestCheckClaimHashedTransfer(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		Convey("After creating a HashedTransfer", func() {
			secret := make([]byte, 32)
			hash := crypto.Sha256(secret)
			hash32 := [32]byte{}
			copy(hash32[:], hash)
			ht := HashedTransfer{
				From:       fixture.Alice.Address,
				To:         fixture.Bob.ID,
				Value:      types.NewBigInt(10),
				HashCommit: hash32,
				Expiry:     10,
			}
			state.SetBlockTime(9)
			Convey("For a valid secret", func() {
				Convey("If the expiry time of the HashedTransfer is less than the block time in context", func() {
					state.SetBlockTime(11)
					Convey("CheckClaimHashedTransfer should return Expired", func() {
						r := CheckClaimHashedTransfer(state, &ht, secret)
						So(r.Code, ShouldEqual, response.ClaimHashedTransferExpired.Code)
					})
				})
				Convey("If the expiry time of the HashedTransfer is equal to the block time in context", func() {
					state.SetBlockTime(10)
					Convey("CheckClaimHashedTransfer should return Expired", func() {
						r := CheckClaimHashedTransfer(state, &ht, secret)
						So(r.Code, ShouldEqual, response.ClaimHashedTransferExpired.Code)
					})
				})
				Convey("If the expiry time of the HashedTransfer is more than the block time in context", func() {
					Convey("CheckClaimHashedTransfer should return Success", func() {
						r := CheckClaimHashedTransfer(state, &ht, secret)
						So(r.Code, ShouldEqual, response.Success.Code)
					})
				})
			})
			Convey("For an invalid secret", func() {
				secret[0] = 1
				Convey("CheckClaimHashedTransfer should return Success", func() {
					r := CheckClaimHashedTransfer(state, &ht, secret)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidSecret.Code)
				})
			})
		})
	})
}

func TestClaimHashedTransfer(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		account.NewAccountFromID(state, fixture.Alice.ID, fixture.Alice.Address)
		account.NewAccountFromID(state, fixture.Bob.ID, fixture.Bob.Address)
		Convey("After creating a HashedTransfer", func() {
			ht := HashedTransfer{
				From:       fixture.Alice.Address,
				To:         fixture.Bob.ID,
				Value:      types.NewBigInt(10),
				HashCommit: [32]byte{},
				Expiry:     10,
			}
			txHash := make([]byte, tmhash.Size)
			CreateHashedTransfer(state, &ht, txHash)
			Convey("After ClaimHashedTransfer", func() {
				ClaimHashedTransfer(state, &ht, txHash)
				Convey("The balance should change accordingly", func() {
					aliceBalance := account.FetchBalance(state, fixture.Alice.ID)
					So(aliceBalance.Cmp(big.NewInt(0)), ShouldBeZeroValue)
					bobBalance := account.FetchBalance(state, fixture.Bob.ID)
					So(bobBalance.Cmp(ht.Value.Int), ShouldBeZeroValue)
					Convey("The HashedTransfer should have been removed", func() {
						So(GetHashedTransfer(state, txHash), ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestCheckRevokeHashedTransfer(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		Convey("After creating a HashedTransfer", func() {
			ht := HashedTransfer{
				From:       fixture.Alice.Address,
				To:         fixture.Bob.ID,
				Value:      types.NewBigInt(10),
				HashCommit: [32]byte{},
				Expiry:     10,
			}
			Convey("If the expiry time of the HashedTransfer is less than the block time in context", func() {
				state.SetBlockTime(11)
				Convey("CheckRevokeHashedTransfer should return Success", func() {
					r := CheckRevokeHashedTransfer(state, &ht)
					So(r.Code, ShouldEqual, response.Success.Code)
				})
			})
			Convey("If the expiry time of the HashedTransfer is equal to the block time in context", func() {
				state.SetBlockTime(10)
				Convey("CheckRevokeHashedTransfer should return Success", func() {
					r := CheckRevokeHashedTransfer(state, &ht)
					So(r.Code, ShouldEqual, response.Success.Code)
				})
			})
			Convey("If the expiry time of the HashedTransfer is more than the block time in context", func() {
				state.SetBlockTime(9)
				Convey("CheckRevokeHashedTransfer should return NotYetExpired", func() {
					r := CheckRevokeHashedTransfer(state, &ht)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferNotYetExpired.Code)
				})
			})
		})
	})
}

func TestRevokeHashedTransfer(t *testing.T) {
	Convey("Given an empty state", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		account.NewAccountFromID(state, fixture.Alice.ID, fixture.Alice.Address)
		account.NewAccountFromID(state, fixture.Bob.ID, fixture.Bob.Address)
		Convey("After creating a HashedTransfer", func() {
			ht := HashedTransfer{
				From:       fixture.Alice.Address,
				To:         fixture.Bob.ID,
				Value:      types.NewBigInt(10),
				HashCommit: [32]byte{},
				Expiry:     10,
			}
			txHash := make([]byte, tmhash.Size)
			CreateHashedTransfer(state, &ht, txHash)
			Convey("After RevokeHashedTransfer", func() {
				RevokeHashedTransfer(state, &ht, txHash)
				Convey("The balance should change accordingly", func() {
					aliceBalance := account.FetchBalance(state, fixture.Alice.ID)
					So(aliceBalance.Cmp(ht.Value.Int), ShouldBeZeroValue)
					bobBalance := account.FetchBalance(state, fixture.Bob.ID)
					So(bobBalance.Cmp(big.NewInt(0)), ShouldBeZeroValue)
					Convey("The HashedTransfer should have been removed", func() {
						So(GetHashedTransfer(state, txHash), ShouldBeNil)
					})
				})
			})
		})
	})
}
