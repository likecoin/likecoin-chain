package response

import (
	"github.com/likecoin/likechain/abci/txstatus"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// R is a struct for response
type R struct {
	Code                 uint32
	Data                 []byte
	Log                  string
	Info                 string
	Tags                 []common.KVPair
	Status               txstatus.TxStatus
	ShouldIncrementNonce bool
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
	if r.Code != 0 {
		r.Code++
	}
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
func (r R) Merge(r2 R) R {
	if r2.Code > 0 {
		r.Code = r2.Code
	}
	if len(r2.Data) > 0 {
		r.Data = r2.Data
	}
	if r2.Log != "" {
		r.Log = r2.Log
	}
	if r2.Info != "" {
		r.Info = r2.Info
	}
	if len(r2.Tags) > 0 {
		r.Tags = append(r.Tags, r2.Tags...)
	}
	if r2.Status != 0 {
		r.Status = r2.Status
	}
	return r
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
	Code:                 0,
	Info:                 "OK",
	Status:               txstatus.TxStatusSuccess,
	ShouldIncrementNonce: true,
}

// Transactions

var RegisterInvalidFormat = R{
	Code:                 10000,
	Info:                 "Invalid RegisterTransaction format",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var RegisterInvalidSignature = R{
	Code:                 10010,
	Info:                 "Invalid RegisterTransaction signature",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var RegisterDuplicated = R{
	Code:                 10020,
	Info:                 "Duplicated RegisterTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var DepositInvalidFormat = R{
	Code:                 11000,
	Info:                 "Invalid DepositTransaction format",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var DepositInvalidSignature = R{
	Code:                 11010,
	Info:                 "Invalid DepositTransaction signature",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var DepositDuplicated = R{
	Code:                 11020,
	Info:                 "Duplicated DepositTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var DepositSenderNotRegistered = R{
	Code:                 11030,
	Info:                 "Sender of DepositTransaction not register",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var DepositInvalidNonce = R{
	Code:                 11040,
	Info:                 "Invalid DepositTransaction nonce",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var DepositNotApprover = R{
	Code:                 11050,
	Info:                 "User is not a DepositApprover",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var DepositDoubleApproval = R{
	Code:                 11060,
	Info:                 "User already approved another DepositTransaction for the same block numner",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var DepositAlreadyExecuted = R{
	Code:                 11070,
	Info:                 "The deposit proposal has already executed for the given block number",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var TransferInvalidFormat = R{
	Code:                 12000,
	Info:                 "Invalid TransferTransaction format",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var TransferInvalidSignature = R{
	Code:                 12010,
	Info:                 "Invalid TransferTransaction signature",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var TransferDuplicated = R{
	Code:                 12020,
	Info:                 "Duplicated TransferTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var TransferSenderNotRegistered = R{
	Code:                 12030,
	Info:                 "Sender of TransferTransaction not register",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var TransferNotEnoughBalance = R{
	Code:                 12040,
	Info:                 "Sender's balance of TransferTransaction not enough",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var TransferInvalidReceiver = R{
	Code:                 12050,
	Info:                 "One or more receivers in TransferTransaction are invalid",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var TransferInvalidNonce = R{
	Code:                 12060,
	Info:                 "Invalid TransferTransaction nonce",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var WithdrawInvalidFormat = R{
	Code:                 13000,
	Info:                 "Invalid WithdrawTransaction format",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var WithdrawInvalidSignature = R{
	Code:                 13010,
	Info:                 "Invalid WithdrawTransaction signature",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var WithdrawDuplicated = R{
	Code:                 13020,
	Info:                 "Duplicated WithdrawTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var WithdrawSenderNotRegistered = R{
	Code:                 13030,
	Info:                 "Sender of WithdrawTransaction not register",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var WithdrawNotEnoughBalance = R{
	Code:                 13040,
	Info:                 "Sender's balance of WithdrawTransaction not enough",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var WithdrawInvalidNonce = R{
	Code:                 13050,
	Info:                 "Invalid WithdrawTransaction nonce",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var HashedTransferInvalidFormat = R{
	Code:                 15000,
	Info:                 "Invalid HashedTransferTransaction format",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var HashedTransferInvalidSignature = R{
	Code:                 15010,
	Info:                 "Invalid HashedTransferTransaction signature",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var HashedTransferDuplicated = R{
	Code:                 15020,
	Info:                 "Duplicated HashedTransferTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var HashedTransferSenderNotRegistered = R{
	Code:                 15030,
	Info:                 "Sender of HashedTransferTransaction not register",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var HashedTransferNotEnoughBalance = R{
	Code:                 15040,
	Info:                 "Sender's balance of HashedTransferTransaction not enough",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var HashedTransferInvalidReceiver = R{
	Code:                 15050,
	Info:                 "The receiver in HashedTransferTransaction is invalid",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var HashedTransferInvalidNonce = R{
	Code:                 15060,
	Info:                 "Invalid HashedTransferTransaction nonce",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var HashedTransferInvalidExpiry = R{
	Code:                 15070,
	Info:                 "Invalid HashedTransferTransaction expiry time",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var ClaimHashedTransferInvalidFormat = R{
	Code:                 16000,
	Info:                 "Invalid ClaimHashedTransferTransaction format",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var ClaimHashedTransferInvalidSignature = R{
	Code:                 16010,
	Info:                 "Invalid ClaimHashedTransferTransaction signature",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var ClaimHashedTransferDuplicated = R{
	Code:                 16020,
	Info:                 "Duplicated ClaimHashedTransferTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var ClaimHashedTransferSenderNotRegistered = R{
	Code:                 16030,
	Info:                 "Sender of ClaimHashedTransferTransaction not register",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var ClaimHashedTransferTxNotExist = R{
	Code:                 16040,
	Info:                 "The HashedTransferTransaction does not exist",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var ClaimHashedTransferExpired = R{
	Code:                 16050,
	Info:                 "The HashedTransferTransaction has already expired",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var ClaimHashedTransferInvalidNonce = R{
	Code:                 16060,
	Info:                 "Invalid ClaimHashedTransferTransaction nonce",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var ClaimHashedTransferInvalidSecret = R{
	Code:                 16070,
	Info:                 "The secret does not match the committed hash of the HashedTransferTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var ClaimHashedTransferNotYetExpired = R{
	Code:                 16080,
	Info:                 "The HashedTransferTransaction is not yet expired",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var ClaimHashedTransferInvalidSender = R{
	Code:                 16090,
	Info:                 "The sender is neither the sender or receiver of the HashedTransferTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var SimpleTransferInvalidFormat = R{
	Code:                 17000,
	Info:                 "Invalid SimpleTransferTransaction format",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var SimpleTransferInvalidSignature = R{
	Code:                 17010,
	Info:                 "Invalid SimpleTransferTransaction signature",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var SimpleTransferDuplicated = R{
	Code:                 17020,
	Info:                 "Duplicated SimpleTransferTransaction",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var SimpleTransferSenderNotRegistered = R{
	Code:                 17030,
	Info:                 "Sender of SimpleTransferTransaction not register",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

var SimpleTransferNotEnoughBalance = R{
	Code:                 17040,
	Info:                 "Sender's balance of SimpleTransferTransaction not enough",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var SimpleTransferInvalidReceiver = R{
	Code:                 17050,
	Info:                 "One or more receivers in SimpleTransferTransaction are invalid",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: true,
}

var SimpleTransferInvalidNonce = R{
	Code:                 17060,
	Info:                 "Invalid SimpleTransferTransaction nonce",
	Status:               txstatus.TxStatusFail,
	ShouldIncrementNonce: false,
}

// Queries

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

var QueryTxNotExist = R{
	Code: 62000,
	Info: "Transaction status record does not exist",
}
