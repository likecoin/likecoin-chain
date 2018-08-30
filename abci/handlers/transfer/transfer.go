package transfer

import (
	"math/big"
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

func logTx(tx *types.TransferTransaction) *logrus.Entry {
	return log.WithField("tx", tx)
}

func checkTransfer(state context.IImmutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetTransferTx()
	if tx == nil {
		log.Panic("Expect TransferTx but got nil")
	}

	if !validateTransferTransactionFormat(state, tx) {
		logTx(tx).Info(response.TransferCheckTxInvalidFormat.Info)
		return response.TransferCheckTxInvalidFormat
	}

	if !validateTransferSignature(state, tx) {
		logTx(tx).Info(response.TransferCheckTxInvalidSignature.Info)
		return response.TransferCheckTxInvalidSignature
	}

	return response.Success // TODO
}

func deliverTransfer(state context.IMutableState, rawTx *types.Transaction) response.R {
	tx := rawTx.GetTransferTx()
	if tx == nil {
		log.Panic("Expect TransferTx but got nil")
	}

	if !validateTransferTransactionFormat(state, tx) {
		logTx(tx).Info(response.TransferDeliverTxInvalidFormat.Info)
		return response.TransferDeliverTxInvalidFormat
	}

	if !validateTransferSignature(state, tx) {
		logTx(tx).Info(response.TransferDeliverTxInvalidSignature.Info)
		return response.TransferDeliverTxInvalidSignature
	}

	fromID, exist := account.GetLikeChainID(state, *tx.From)
	if !exist {
		return response.Success // TODO: error code for sender account does not exist
	}

	_ = account.FetchBalance(state, fromID)
	_ = account.FetchNextNonce(state, fromID)
	// Increment nonce
	// Adjust balance of sender and receiver

	return response.Success // TODO
}

func validateTransferSignature(state context.IImmutableState, tx *types.TransferTransaction) bool {
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

	senderAddr := tx.From.GetAddr()
	if senderAddr != nil {
		if senderAddr.ToEthereum() == sigAddr {
			return true
		}
		log.WithFields(logrus.Fields{
			"txAddr":  senderAddr.ToHex(),
			"sigAddr": sigAddr.Hex(),
		}).Info("Recovered address is not match")
	} else {
		id := tx.From.GetLikeChainID()
		if id != nil {
			if account.IsLikeChainIDHasAddress(state, id, &sigAddr) {
				return true
			}
			log.WithFields(logrus.Fields{
				"likeChainID": id.ToString(),
				"sigAddr":     sigAddr.Hex(),
			}).Info("Recovered address is not bind to the LikeChain ID of the sender")
		}
	}

	return false
}

func validateTransferTransactionFormat(state context.IImmutableState, tx *types.TransferTransaction) bool {
	if !tx.From.IsValidFormat() {
		log.Debug("Invalid sender format in transfer transaction")
		return false
	}

	if len(tx.ToList) > 0 {
		for _, target := range tx.ToList {
			if !target.IsValidFormat() {
				log.Debug("Invalid receiver format in transfer transaction")
				return false
			}
		}
	} else {
		log.Debug("No receiver in transfer transaction")
		return false
	}

	if tx.Fee.ToBigInt().Cmp(big.NewInt(0)) < 0 {
		log.Debug("Invalid fee in transfer transaction")
		return false
	}

	if !tx.Sig.IsValidFormat() {
		log.Debug("Invalid signature format in transfer transaction")
		return false
	}

	return true
}

func transfer(state context.IMutableState, tx *types.TransferTransaction) {
	// TODO
}

func init() {
	t := reflect.TypeOf((*types.Transaction_TransferTx)(nil))
	handler.RegisterCheckTxHandler(t, checkTransfer)
	handler.RegisterDeliverTxHandler(t, deliverTransfer)
}
