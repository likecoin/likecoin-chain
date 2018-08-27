package withdraw

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	handler "github.com/likecoin/likechain/abci/handlers"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	"github.com/sirupsen/logrus"
)

var log = logger.L

func logTx(tx *types.WithdrawTransaction) *logrus.Entry {
	return log.WithField("tx", tx)
}

func checkWithdraw(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetWithdrawTx()
	if tx == nil {
		log.Panic("Expect WithdrawTx but got nil")
	}

	_ = tx.From

	if !validateWithdrawTransactionFormat(tx) {
		logTx(tx).Info(response.WithdrawCheckTxInvalidFormat.Info)
		return response.WithdrawCheckTxInvalidFormat
	}

	if !validateWithdrawSignature(tx.Sig) {
		logTx(tx).Info(response.WithdrawCheckTxInvalidSignature.Info)
		return response.WithdrawCheckTxInvalidSignature
	}

	return response.Success // TODO
}

func deliverWithdraw(state context.IMutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetWithdrawTx()
	if tx == nil {
		log.Panic("Expect WithdrawTx but got nil")
	}

	if !validateWithdrawTransactionFormat(tx) {
		logTx(tx).Info(response.WithdrawDeliverTxInvalidFormat.Info)
		return response.WithdrawDeliverTxInvalidFormat
	}

	if !validateWithdrawSignature(tx.Sig) {
		logTx(tx).Info(response.WithdrawDeliverTxInvalidSignature.Info)
		return response.WithdrawDeliverTxInvalidSignature
	}

	return response.Success // TODO
}

func validateWithdrawSignature(sig *types.Signature) bool {
	return false // TODO
}

func validateWithdrawTransactionFormat(tx *types.WithdrawTransaction) bool {
	return false // TODO
}

func withdraw(state context.IMutableState, tx *types.WithdrawTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_WithdrawTx)(nil))
	handler.RegisterCheckTxHandler(t, checkWithdraw)
	handler.RegisterDeliverTxHandler(t, deliverWithdraw)
}
