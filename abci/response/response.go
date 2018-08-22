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
