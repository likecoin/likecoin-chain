package query

import (
	"fmt"
	"testing"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuery(t *testing.T) {
	appCtx := context.NewMock()
	state := appCtx.GetMutableState()

	Convey("Given a query", t, func() {
		reqQuery := abci.RequestQuery{
			Path: "tx_state",
			Data: []byte(""),
		}

		Convey("If its path is invalid", func() {
			reqQuery.Path = ""
			res := Query(state, reqQuery)

			code := response.QueryPathNotExist.Code
			Convey(fmt.Sprintf("Should return code %d", code), func() {
				So(res.Code, ShouldEqual, response.QueryPathNotExist.Code)
			})
		})

		Convey("If its data is invalid", func() {
			reqQuery.Data = nil
			res := Query(state, reqQuery)

			code := response.QueryParsingRequestError.Code
			Convey(fmt.Sprintf("Should return code %d", code), func() {
				So(res.Code, ShouldEqual, response.QueryParsingRequestError.Code)
			})
		})
	})
}

func TestJsonMapToResponse(t *testing.T) {
	Convey("Given a JSON map", t, func() {
		m := jsonMap{"foo": "bar"}

		Convey("If it is valid", func() {
			Convey("It should return code 0", func() {
				So(m.ToResponse().Code, ShouldEqual, 0)
			})
		})

		Convey("If it is invalid", func() {
			m["foo"] = func() {}
			code := response.QueryParsingResponseError.Code
			Convey(fmt.Sprintf("It should return code %d", code), func() {
				So(m.ToResponse().Code, ShouldEqual, code)
			})
		})
	})
}
