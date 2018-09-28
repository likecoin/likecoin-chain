package deposit

import (
	"reflect"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers/table"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	"github.com/sirupsen/logrus"
)

var log = logger.L

func logTx(tx *types.DepositTransaction) *logrus.Entry {
	return log.WithField("tx", tx)
}

func checkDeposit(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetDepositTx()
	if tx == nil {
		log.Panic("Expect DepositTx but got nil")
	}

	_ = tx.BlockNumber

	if !validateDepositTransactionFormat(tx) {
		logTx(tx).Info(response.DepositCheckTxInvalidFormat.Info)
		return response.DepositCheckTxInvalidFormat
	}

	return response.Success // TODO
}

func deliverDeposit(state context.IMutableState, rawTx *types.Transaction, txHash []byte) response.R {
	tx := rawTx.GetDepositTx()
	if tx == nil {
		log.Panic("Expect DepositTx but got nil")
	}

	_ = tx.BlockNumber

	if !validateDepositTransactionFormat(tx) {
		logTx(tx).Info(response.DepositDeliverTxInvalidFormat.Info)
		return response.DepositDeliverTxInvalidFormat
	}

	return response.Success // TODO
}

func validateDepositTransactionFormat(tx *types.DepositTransaction) bool {
	return false // TODO
}

func deposit(state context.IMutableState, tx *types.DepositTransaction) {
	// TODO
}

func init() {
	log.Info("Init transfer handlers")
	_type := reflect.TypeOf((*types.Transaction_DepositTx)(nil))
	table.RegisterCheckTxHandler(_type, checkDeposit)
	table.RegisterDeliverTxHandler(_type, deliverDeposit)
}
