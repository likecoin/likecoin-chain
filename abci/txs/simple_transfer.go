package txs

import (
	"math/big"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

// SimpleTransferTransaction represents a SimpleTransfer transaction
type SimpleTransferTransaction struct {
	From   types.Identifier
	To     types.Identifier
	Value  types.BigInt
	Remark string
	Fee    types.BigInt
	Nonce  uint64
	Sig    SimpleTransferSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *SimpleTransferTransaction) ValidateFormat() bool {
	if tx.From == nil || tx.To == nil || tx.Value.Int == nil || tx.Fee.Int == nil || tx.Sig == nil {
		return false
	}
	if !tx.Value.IsWithinRange() || !tx.Fee.IsWithinRange() {
		return false
	}
	if len(tx.Remark) > 4096 {
		return false
	}
	return true
}

func (tx *SimpleTransferTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.SimpleTransferInvalidFormat.Info)
		return response.SimpleTransferInvalidFormat, nil
	}
	senderID = account.IdentifierToLikeChainID(state, tx.From)
	if senderID == nil {
		logTx(tx).Info(response.SimpleTransferSenderNotRegistered.Info)
		return response.SimpleTransferSenderNotRegistered, senderID
	}
	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.SimpleTransferInvalidSignature.Info)
		return response.SimpleTransferInvalidSignature, senderID
	}
	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.SimpleTransferInvalidNonce.Info)
		return response.SimpleTransferInvalidNonce, senderID
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.SimpleTransferDuplicated.Info)
		return response.SimpleTransferDuplicated, senderID
	}
	if toID, ok := tx.To.(*types.LikeChainID); ok {
		if !(account.IsLikeChainIDRegistered(state, toID)) {
			logTx(tx).
				WithField("to", tx.To).
				Info(response.SimpleTransferInvalidReceiver.Info)
			return response.SimpleTransferInvalidReceiver, senderID
		}
	}
	senderBalance := account.FetchBalance(state, senderID)
	total := new(big.Int).Set(tx.Value.Int)
	total.Add(total, tx.Fee.Int)
	if senderBalance.Cmp(total) < 0 {
		logTx(tx).
			WithField("total", total.String()).
			WithField("balance", senderBalance.String()).
			Info(response.SimpleTransferNotEnoughBalance.Info)
		return response.SimpleTransferNotEnoughBalance, senderID
	}
	return response.Success, senderID
}

// CheckTx checks the transaction to see if it should be executed
func (tx *SimpleTransferTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *SimpleTransferTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID := tx.checkTx(state)
	if checkTxRes.Code != response.Success.Code {
		if checkTxRes.ShouldIncrementNonce {
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	account.IncrementNextNonce(state, senderID)

	total := new(big.Int).Set(tx.Value.Int)
	total.Add(total, tx.Fee.Int)
	account.AddBalance(state, tx.To, tx.Value.Int)
	account.MinusBalance(state, senderID, total)

	return response.Success
}

// SimpleTransferTx returns a TransferTransaction
func SimpleTransferTx(from, to types.Identifier, value types.BigInt, remark string, fee types.BigInt, nonce uint64, sigHex string) *SimpleTransferTransaction {
	sig := &SimpleTransferJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &SimpleTransferTransaction{
		From:   from,
		To:     to,
		Value:  value,
		Remark: remark,
		Fee:    fee,
		Nonce:  nonce,
		Sig:    sig,
	}
}

// RawSimpleTransferTx returns raw bytes of a TransferTransaction
func RawSimpleTransferTx(from, to types.Identifier, value types.BigInt, remark string, fee types.BigInt, nonce uint64, sigHex string) []byte {
	return EncodeTx(SimpleTransferTx(from, to, value, remark, fee, nonce, sigHex))
}
