package response

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// R is a struct for response
type R struct {
	Code uint32
	Data []byte
	Log  string
	Info string
	Tags []common.KVPair
}

// ToResponseCheckTx converts R to abci ResponseCheckTx
func (r R) ToResponseCheckTx() abci.ResponseCheckTx {
	return abci.ResponseCheckTx{
		Code: r.Code,
		Data: r.Data,
		Log:  r.Log,
		Info: r.Info,
		Tags: r.Tags,
	}
}

// ToResponseDeliverTx converts R to abci ResponseDeliverTx
func (r R) ToResponseDeliverTx() abci.ResponseDeliverTx {
	return abci.ResponseDeliverTx{
		Code: r.Code,
		Data: r.Data,
		Log:  r.Log,
		Info: r.Info,
		Tags: r.Tags,
	}
}

// ToResponseQuery converts R to abci ResponseQuery
func (r R) ToResponseQuery() abci.ResponseQuery {
	return abci.ResponseQuery{
		Code:  r.Code,
		Log:   r.Log,
		Info:  r.Info,
		Value: r.Data,
	}
}

// Merge merges R2 into R1
func (r1 R) Merge(r2 R) R {
	if r2.Code > 0 {
		r1.Code = r2.Code
	}
	if len(r2.Data) > 0 {
		r1.Data = r2.Data
	}
	if r2.Log != "" {
		r1.Log = r2.Log
	}
	if r2.Info != "" {
		r1.Info = r2.Info
	}
	if len(r2.Tags) > 0 {
		r1.Tags = append(r1.Tags, r2.Tags...)
	}
	return r1
}

// Error Code Definition (5 digits)
// 1 0 0 0 0
// | | | | |
// | | | | Type of transaction
// | | Case
// | | - R case 10-59
// | | - Other case  60-99
// Type
// - Before parsing request 00-09
// - Transaction            10-59
// - Query                  60-99

var Success = R{
	Code: 0,
	Info: "OK",
}

var RegisterCheckTxInvalidFormat = R{
	Code: 10000,
	Info: "Invalid RegisterTransaction format in CheckTx",
}

var RegisterDeliverTxInvalidFormat = R{
	Code: 10001,
	Info: "Invalid RegisterTransaction format in DeliverTx",
}

var RegisterCheckTxInvalidSignature = R{
	Code: 10010,
	Info: "Invalid RegisterTransaction signature in CheckTx",
}

var RegisterDeliverTxInvalidSignature = R{
	Code: 10011,
	Info: "Invalid RegisterTransaction signature in DeliverTx",
}

var RegisterCheckTxDuplicated = R{
	Code: 10020,
	Info: "Duplicated RegisterTransaction in CheckTx",
}

var RegisterDeliverTxDuplicated = R{
	Code: 10021,
	Info: "Duplicated RegisterTransaction in DeliverTx",
}

var DepositCheckTxInvalidFormat = R{
	Code: 11000,
	Info: "Invalid DepositTransaction format in CheckTx",
}

var DepositDeliverTxInvalidFormat = R{
	Code: 11001,
	Info: "Invalid DepositTransaction format in DeliverTx",
}

var DepositCheckTxDuplicated = R{
	Code: 11020,
	Info: "Duplicated DepositTransaction in CheckTx",
}

var DepositDeliverTxDuplicated = R{
	Code: 11021,
	Info: "Duplicated DepositTransaction in DeliverTx",
}

var TransferCheckTxInvalidFormat = R{
	Code: 12000,
	Info: "Invalid TransferTransaction format in CheckTx",
}

var TransferDeliverTxInvalidFormat = R{
	Code: 12001,
	Info: "Invalid TransferTransaction format in DeliverTx",
}

var TransferCheckTxInvalidSignature = R{
	Code: 12010,
	Info: "Invalid TransferTransaction signature in CheckTx",
}

var TransferDeliverTxInvalidSignature = R{
	Code: 12011,
	Info: "Invalid TransferTransaction signature in DeliverTx",
}

var TransferCheckTxDuplicated = R{
	Code: 12020,
	Info: "Duplicated TransferTransaction in CheckTx",
}

var TransferDeliverTxDuplicated = R{
	Code: 12021,
	Info: "Duplicated TransferTransaction in DeliverTx",
}

var TransferCheckTxSenderNotRegistered = R{
	Code: 12030,
	Info: "Sender of TransferTransaction not register in CheckTx",
}

var TransferDeliverTxSenderNotRegistered = R{
	Code: 12031,
	Info: "Sender of TransferTransaction not register in DeliverTx",
}

var TransferCheckTxNotEnoughBalance = R{
	Code: 12040,
	Info: "Sender's balance of TransferTransaction not enough in CheckTx",
}

var TransferDeliverTxNotEnoughBalance = R{
	Code: 12041,
	Info: "Sender's balance of TransferTransaction not enough in DeliverTx",
}

var TransferCheckTxInvalidNonce = R{
	Code: 12600,
	Info: "Invalid TransferTransaction nonce in CheckTx",
}

var TransferDeliverTxInvalidNonce = R{
	Code: 12601,
	Info: "Invalid TransferTransaction nonce in DeliverTx",
}

var WithdrawCheckTxInvalidFormat = R{
	Code: 13000,
	Info: "Invalid WithdrawTransaction format in CheckTx",
}

var WithdrawDeliverTxInvalidFormat = R{
	Code: 13001,
	Info: "Invalid WithdrawTransaction format in DeliverTx",
}

var WithdrawCheckTxInvalidSignature = R{
	Code: 13010,
	Info: "Invalid WithdrawTransaction signature in CheckTx",
}

var WithdrawDeliverTxInvalidSignature = R{
	Code: 13011,
	Info: "Invalid WithdrawTransaction signature in DeliverTx",
}

var WithdrawCheckTxDuplicated = R{
	Code: 13020,
	Info: "Duplicated WithdrawTransaction in CheckTx",
}

var WithdrawDeliverTxDuplicated = R{
	Code: 13021,
	Info: "Duplicated WithdrawTransaction in DeliverTx",
}

var WithdrawCheckTxSenderNotRegistered = R{
	Code: 13030,
	Info: "Sender of WithdrawTransaction not register in CheckTx",
}

var WithdrawDeliverTxSenderNotRegistered = R{
	Code: 13031,
	Info: "Sender of WithdrawTransaction not register in DeliverTx",
}

var WithdrawCheckTxNotEnoughBalance = R{
	Code: 13040,
	Info: "Sender's balance of WithdrawTransaction not enough in CheckTx",
}

var WithdrawDeliverTxNotEnoughBalance = R{
	Code: 13041,
	Info: "Sender's balance of WithdrawTransaction not enough in DeliverTx",
}

var WithdrawCheckTxInvalidNonce = R{
	Code: 13600,
	Info: "Invalid WithdrawTransaction nonce in CheckTx",
}

var WithdrawDeliverTxInvalidNonce = R{
	Code: 13601,
	Info: "Invalid WithdrawTransaction nonce in DeliverTx",
}

var QueryPathNotExist = R{
	Code: 60010,
	Info: "Invalid query path",
}

var QueryParsingRequestError = R{
	Code: 60020,
	Info: "Unable to parse request data in Query",
}

var QueryParsingResponseError = R{
	Code: 60030,
	Info: "Unable to parse response data in Query",
}

var QueryInvalidIdentifier = R{
	Code: 60040,
	Info: "Identifier is invalid in Query",
}

var QueryWithdrawProofInvalidHeight = R{
	Code: 61000,
	Info: "Invalid height in withdraw proof",
}

var QueryWithdrawProofNotExist = R{
	Code: 61010,
	Info: "Withdraw record does not exist",
}
