package register

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers/table"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	"github.com/sirupsen/logrus"
)

var log = logger.L

func logTx(tx *types.RegisterTransaction) *logrus.Entry {
	return log.WithField("tx", tx.ToString())
}

func checkRegister(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		log.Panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransactionFormat(tx) {
		logTx(tx).Info(response.RegisterCheckTxInvalidFormat.Info)
		return response.RegisterCheckTxInvalidFormat
	}

	if !validateRegisterSignature(state, tx) {
		logTx(tx).Info(response.RegisterCheckTxInvalidSignature.Info)
		return response.RegisterCheckTxInvalidSignature
	}

	if account.IsAddressRegistered(state, tx.Addr.ToEthereum()) {
		logTx(tx).Info(response.RegisterCheckTxDuplicated.Info)
		return response.RegisterCheckTxDuplicated
	}

	return response.Success
}

func deliverRegister(state context.IMutableState, rawTx *types.Transaction, txHash []byte) response.R {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		log.Panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransactionFormat(tx) {
		logTx(tx).Info(response.RegisterDeliverTxInvalidFormat.Info)
		return response.RegisterDeliverTxInvalidFormat
	}

	if !validateRegisterSignature(state, tx) {
		logTx(tx).Info(response.RegisterDeliverTxInvalidSignature.Info)
		return response.RegisterDeliverTxInvalidSignature
	}

	if account.IsAddressRegistered(state, tx.Addr.ToEthereum()) {
		logTx(tx).Info(response.RegisterDeliverTxDuplicated.Info)
		return response.RegisterDeliverTxDuplicated
	}

	ethAddr := tx.Addr.ToEthereum()
	id, _ := account.NewAccount(state, ethAddr)

	return response.Success.Merge(response.R{
		Data: id.Content,
	})
}

// validateRegisterSignature validates register transaction
func validateRegisterSignature(state context.IImmutableState, tx *types.RegisterTransaction) bool {
	hashedMsg := tx.GenerateSigningMessageHash()
	sigAddr, err := utils.RecoverSignature(hashedMsg, tx.Sig)
	if err != nil {
		log.WithError(err).Info("Unable to recover signature when validating signature")
		return false
	}

	if tx.Addr.ToEthereum() != sigAddr {
		log.WithFields(logrus.Fields{
			"tx_addr":  tx.Addr.ToHex(),
			"sig_addr": sigAddr.Hex(),
		}).Info("Recovered address is not match")
		return false
	}

	return true
}

// validateRegisterTransactionFormat validates register transaction
func validateRegisterTransactionFormat(tx *types.RegisterTransaction) bool {
	return tx.Addr.IsValidFormat() && tx.Sig.IsValidFormat()
}

func init() {
	log.Info("Init register handlers")
	_type := reflect.TypeOf((*types.Transaction_RegisterTx)(nil))
	table.RegisterCheckTxHandler(_type, checkRegister)
	table.RegisterDeliverTxHandler(_type, deliverRegister)
}
