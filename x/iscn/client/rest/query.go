package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	gocid "github.com/ipfs/go-cid"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/likecoin/likechain/x/iscn/types"
	"github.com/multiformats/go-multibase"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/iscn/params",
		paramsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/iscn/kernels/{iscnID}",
		kernelHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/iscn/cids/{cid}",
		cidHandlerFn(cliCtx),
	).Methods("GET")
}

func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryParams))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func kernelHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		iscnIDStr := vars["iscnID"]
		// TODO: proper decode by ISCN ID format
		_, iscnID, err := multibase.Decode(iscnIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		endpoint := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryIscnKernel)

		res, height, err := cliCtx.QueryWithData(endpoint, iscnID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_, cid, err := gocid.CidFromBytes(res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		bz, err := cliCtx.Codec.MarshalJSONIndent(cid, "", "  ")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}

func cidHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cidStr := vars["cid"]
		cid, err := gocid.Decode(cidStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		bz := cid.Bytes()

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		endpoint := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryCID)

		res, height, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
