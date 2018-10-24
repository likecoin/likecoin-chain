package query

import (
	"fmt"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryAccountInfo(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()
	account.NewAccountFromID(state, fixture.Alice.ID, fixture.Alice.Address)

	Convey("Given an account info query", t, func() {
		reqQuery := abci.RequestQuery{
			Data: []byte(fixture.Alice.ID.String()),
			Path: "account_info",
		}

		Convey("If it is a valid query using LikeChain ID as identity", func() {
			res := Query(state, reqQuery)

			Convey("Should return code 0", func() {
				So(res.Code, ShouldEqual, 0)
			})
		})

		Convey("If it is a valid query using address as identity", func() {
			reqQuery.Data = []byte(fixture.Alice.Address.String())
			res := Query(state, reqQuery)

			Convey("Should return code 0", func() {
				So(res.Code, ShouldEqual, 0)
			})
		})

		Convey("If its identity is not registered", func() {
			reqQuery.Data = []byte(fixture.Bob.ID.String())
			res := Query(state, reqQuery)

			code := response.QueryInvalidIdentifier.Code
			Convey(fmt.Sprintf("Should return code %d", code), func() {
				So(res.Code, ShouldEqual, response.QueryInvalidIdentifier.Code)
			})
		})
	})
}
