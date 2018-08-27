package register

import (
	"reflect"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	handler "github.com/likecoin/likechain/abci/handlers"
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

	if !validateRegisterTransaction(tx) {
		logTx(tx).Info(response.RegisterCheckTxInvalidFormat.Info)
		return response.RegisterCheckTxInvalidFormat
	}

	if !validateRegisterSignature(state, tx) {
		logTx(tx).Info(response.RegisterCheckTxInvalidSignature.Info)
		return response.RegisterCheckTxInvalidSignature
	}

	_, existed := state.ImmutableStateTree().Get(utils.DbAddrKey(tx.Addr.ToEthereum()))
	if existed != nil {
		logTx(tx).Info(response.RegisterCheckTxDuplicated.Info)
		return response.RegisterCheckTxDuplicated
	}

	return response.Success
}

func deliverRegister(state context.IMutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetRegisterTx()
	if tx == nil {
		log.Panic("Expect RegisterTx but got nil")
	}

	if !validateRegisterTransaction(tx) {
		logTx(tx).Info(response.RegisterDeliverTxInvalidFormat.Info)
		return response.RegisterDeliverTxInvalidFormat
	}

	if !validateRegisterSignature(state, tx) {
		logTx(tx).Info(response.RegisterDeliverTxInvalidSignature.Info)
		return response.RegisterDeliverTxInvalidSignature
	}

	_, existed := state.ImmutableStateTree().Get(utils.DbAddrKey(tx.Addr.ToEthereum()))
	if existed != nil {
		logTx(tx).Info(response.RegisterDeliverTxDuplicated.Info)
		return response.RegisterDeliverTxDuplicated
	}

	id, err := register(state, tx)
	if err != nil {
		log.WithError(err).Panic("Error occurs during registration")
	}

	return response.Success.Merge(response.R{
		Data: id.Content,
	})
}

// validateRegisterSignature validates register transaction
func validateRegisterSignature(state context.IImmutableState, tx *types.RegisterTransaction) bool {
	hashedMsg, err := tx.GenerateSigningMessageHash()
	if err != nil {
		log.WithError(err).Info("Unable to generate signing message hash when validating signature")
		return false
	}

	sigAddr, err := utils.RecoverSignature(hashedMsg, tx.Sig)
	if err != nil {
		log.WithError(err).Info("Unable to recover signature when validating signature")
		return false
	}

	if tx.Addr.ToEthereum() != sigAddr {
		log.WithFields(logrus.Fields{
			"txAddr":  tx.Addr.ToHex(),
			"sigAddr": sigAddr.Hex(),
		}).Info("Recovered address is not match")
		return false
	}

	return true
}

// validateRegisterTransaction validates register transaction
func validateRegisterTransaction(tx *types.RegisterTransaction) bool {
	return tx.Addr.IsValidFormat() && tx.Sig.IsValidFormat()
}

// register creates a new LikeChain account
func register(state context.IMutableState, tx *types.RegisterTransaction) (types.LikeChainID, error) {
	ethAddr := tx.Addr.ToEthereum()
	return account.NewAccount(state, ethAddr)
}

func init() {
	t := reflect.TypeOf((*types.Transaction_RegisterTx)(nil))
	handler.RegisterCheckTxHandler(t, checkRegister)
	handler.RegisterDeliverTxHandler(t, deliverRegister)
}
