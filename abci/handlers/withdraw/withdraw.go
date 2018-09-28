package withdraw

import (
	"math/big"
	"reflect"
	"strconv"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers/table"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
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

	if !validateWithdrawTransactionFormat(tx) {
		logTx(tx).Info(response.WithdrawCheckTxInvalidFormat.Info)
		return response.WithdrawCheckTxInvalidFormat
	}

	senderID := account.IdentifierToLikeChainID(state, tx.From)
	if senderID == nil {
		logTx(tx).Info(response.TransferCheckTxSenderNotRegistered.Info)
		return response.TransferCheckTxSenderNotRegistered
	}

	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.WithdrawCheckTxInvalidNonce.Info)
		return response.WithdrawCheckTxInvalidNonce
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.WithdrawCheckTxDuplicated.Info)
		return response.WithdrawCheckTxDuplicated
	}

	if !validateWithdrawSignature(state, tx) {
		logTx(tx).Info(response.WithdrawCheckTxInvalidSignature.Info)
		return response.WithdrawCheckTxInvalidSignature
	}

	senderBalance := account.FetchBalance(state, tx.From)
	amount := tx.Value.ToBigInt()
	if senderBalance.Cmp(amount) < 0 {
		logTx(tx).Info(response.WithdrawCheckTxNotEnoughBalance.Info)
		return response.WithdrawCheckTxNotEnoughBalance
	}

	// TODO: check fee

	return response.Success
}

func deliverWithdraw(state context.IMutableState, rawTx *types.Transaction, txHash []byte) response.R {
	tx := rawTx.GetWithdrawTx()
	if tx == nil {
		log.Panic("Expect WithdrawTx but got nil")
	}

	if !validateWithdrawTransactionFormat(tx) {
		logTx(tx).Info(response.WithdrawDeliverTxInvalidFormat.Info)
		return response.WithdrawDeliverTxInvalidFormat
	}

	senderID := account.IdentifierToLikeChainID(state, tx.From)
	if senderID == nil {
		logTx(tx).Info(response.TransferDeliverTxSenderNotRegistered.Info)
		return response.TransferDeliverTxSenderNotRegistered
	}

	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.WithdrawDeliverTxInvalidNonce.Info)
		return response.WithdrawDeliverTxInvalidNonce
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.WithdrawDeliverTxDuplicated.Info)
		return response.WithdrawDeliverTxDuplicated
	}

	if !validateWithdrawSignature(state, tx) {
		logTx(tx).Info(response.WithdrawDeliverTxInvalidSignature.Info)
		return response.WithdrawDeliverTxInvalidSignature
	}

	// TODO: check fee

	senderBalance := account.FetchBalance(state, tx.From)
	amount := tx.Value.ToBigInt()
	amount.Add(amount, tx.Fee.ToBigInt())
	if senderBalance.Cmp(amount) < 0 {
		logTx(tx).Info(response.WithdrawDeliverTxNotEnoughBalance.Info)
		return response.WithdrawDeliverTxNotEnoughBalance
	}

	account.MinusBalance(state, tx.From, amount)
	account.IncrementNextNonce(state, senderID)

	// Normalize Identifier for packedTx
	tx.From = senderID.ToIdentifier()
	packedTx := tx.Pack()

	withdrawTree := state.MutableWithdrawTree()
	withdrawTree.Set(crypto.Sha256(packedTx), []byte{1})

	return response.Success.Merge(response.R{
		Tags: []common.KVPair{
			{
				Key:   []byte("withdraw.height"),
				Value: []byte(strconv.FormatInt(state.GetHeight()+1, 10)),
			},
		},
		Data: packedTx,
	})
}

func validateWithdrawSignature(state context.IImmutableState, tx *types.WithdrawTransaction) bool {
	hashedMsg := tx.GenerateSigningMessageHash()
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
			"tx_addr":  senderAddr.ToHex(),
			"sig_addr": sigAddr.Hex(),
		}).Info("Recovered address is not match")
	} else {
		id := tx.From.GetLikeChainID()
		if id != nil {
			if account.IsLikeChainIDHasAddress(state, id, sigAddr) {
				return true
			}
			log.WithFields(logrus.Fields{
				"likechain_id": id.ToString(),
				"sig_addr":     sigAddr.Hex(),
			}).Info("Recovered address is not bind to the LikeChain ID of the sender")
		}
	}

	return false
}

func validateWithdrawTransactionFormat(tx *types.WithdrawTransaction) bool {
	if !tx.From.IsValidFormat() {
		log.Debug("Invalid sender format in withdraw transaction")
		return false
	}
	if !tx.ToAddr.IsValidFormat() {
		log.Debug("Invalid receiver Ethereum address format in withdraw transaction")
		return false
	}
	limit := big.NewInt(2)
	limit.Exp(limit, big.NewInt(256), nil)
	zero := big.NewInt(0)
	value := tx.Value.ToBigInt()
	if value.Cmp(zero) <= 0 || value.Cmp(limit) >= 0 {
		log.Debug("Invalid value range in wighdraw transaction")
		return false
	}
	fee := tx.Fee.ToBigInt()
	if fee.Cmp(zero) < 0 || fee.Cmp(limit) >= 0 {
		log.Debug("Invalid fee range in wighdraw transaction")
		return false
	}
	if !tx.Sig.IsValidFormat() {
		log.Debug("Invalid signature format in withdraw transaction")
		return false
	}
	return true
}

func init() {
	log.Info("Init withdraw handlers")
	_type := reflect.TypeOf((*types.Transaction_WithdrawTx)(nil))
	table.RegisterCheckTxHandler(_type, checkWithdraw)
	table.RegisterDeliverTxHandler(_type, deliverWithdraw)
}
