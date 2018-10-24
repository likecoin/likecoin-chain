package query

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"

	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryAddressInfo(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()
	account.NewAccountFromID(state, fixture.Alice.ID, fixture.Alice.Address)

	account.SaveBalance(state, fixture.Bob.Address, big.NewInt(1))

	Convey("Given an address info query", t, func() {
		reqQuery := abci.RequestQuery{
			Data: []byte(fixture.Alice.Address.String()),
			Path: "address_info",
		}

		Convey("If it is a valid query for an Ethereum address with LikeChain ID", func() {
			res := Query(state, reqQuery)

			Convey("Should successfully return code 0", func() {
				So(res.Code, ShouldEqual, 0)
			})
		})

		Convey("If it is a valid query for an Ethereum address without LikeChain ID but with some balance", func() {
			reqQuery.Data = []byte(fixture.Bob.Address.String())
			res := Query(state, reqQuery)

			Convey("Should return code 0", func() {
				So(res.Code, ShouldEqual, 0)
			})
		})

		Convey("If it is a valid query for an Ethereum address without LikeChain ID and no balance", func() {
			reqQuery.Data = []byte(fixture.Carol.Address.String())
			res := Query(state, reqQuery)

			Convey("Should return code 0", func() {
				So(res.Code, ShouldEqual, 0)
			})
		})

		Convey("If the query contains invalid character", func() {
			reqQuery.Data = []byte("0x000000000000000000000000000000000000000g")
			res := Query(state, reqQuery)

			code := response.QueryInvalidIdentifier.Code
			Convey(fmt.Sprintf("Should return code %d", code), func() {
				So(res.Code, ShouldEqual, code)
			})
		})

		Convey("If the query length is wrong", func() {
			reqQuery.Data = []byte("0x000000000000000000000000000000000000000")
			res := Query(state, reqQuery)

			code := response.QueryInvalidIdentifier.Code
			Convey(fmt.Sprintf("Should return code %d", code), func() {
				So(res.Code, ShouldEqual, code)
			})
		})
	})
}
