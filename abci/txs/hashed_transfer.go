package txs

import (
	"math/big"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/htlc"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
)

// HashedTransferTransaction represents a Hashed TimeLock Transfer transaction
type HashedTransferTransaction struct {
	HashedTransfer htlc.HashedTransfer
	Nonce          uint64
	Fee            types.BigInt
	Sig            HashedTransferSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *HashedTransferTransaction) ValidateFormat() bool {
	if tx.Fee.Int == nil || tx.Sig == nil {
		return false
	}
	if !tx.Fee.IsWithinRange() {
		return false
	}
	return tx.HashedTransfer.Validate()
}

func (tx *HashedTransferTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.HashedTransferInvalidFormat.Info)
		return response.HashedTransferInvalidFormat, nil
	}
	senderID = account.IdentifierToLikeChainID(state, tx.HashedTransfer.From)
	if senderID == nil {
		logTx(tx).Info(response.HashedTransferSenderNotRegistered.Info)
		return response.HashedTransferSenderNotRegistered, senderID
	}
	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.HashedTransferInvalidSignature.Info)
		return response.HashedTransferInvalidSignature, senderID
	}
	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.HashedTransferInvalidNonce.Info)
		return response.HashedTransferInvalidNonce, senderID
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.HashedTransferDuplicated.Info)
		return response.HashedTransferDuplicated, senderID
	}
	toID := account.IdentifierToLikeChainID(state, tx.HashedTransfer.To)
	if toID == nil {
		logTx(tx).
			WithField("to", tx.HashedTransfer.To).
			Info(response.HashedTransferInvalidReceiver.Info)
		return response.HashedTransferInvalidReceiver, senderID
	}
	senderBalance := account.FetchBalance(state, senderID)
	total := new(big.Int).Set(tx.HashedTransfer.Value.Int)
	total.Add(total, tx.Fee.Int)
	if senderBalance.Cmp(total) < 0 {
		logTx(tx).
			WithField("total", total.String()).
			WithField("balance", senderBalance.String()).
			Info(response.HashedTransferNotEnoughBalance.Info)
		return response.HashedTransferNotEnoughBalance, senderID
	}
	return htlc.CheckCreateHashedTransfer(state, &tx.HashedTransfer), senderID
}

// CheckTx checks the transaction to see if it should be executed
func (tx *HashedTransferTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *HashedTransferTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID := tx.checkTx(state)
	if checkTxRes.Code != 0 {
		switch checkTxRes.Code {
		case response.HashedTransferInvalidReceiver.Code:
			fallthrough
		case response.HashedTransferNotEnoughBalance.Code:
			fallthrough
		case response.HashedTransferInvalidExpiry.Code:
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	account.IncrementNextNonce(state, senderID)

	total := new(big.Int).Set(tx.HashedTransfer.Value.Int)
	total.Add(total, tx.Fee.Int)
	account.MinusBalance(state, senderID, total)
	htlc.CreateHashedTransfer(state, &tx.HashedTransfer, txHash)

	return response.Success.Merge(response.R{
		Status: txstatus.TxStatusPending,
	})
}

// HashedTransferTx returns a HashedTransferTransaction
func HashedTransferTx(from, to types.Identifier, value int64, commit []byte, expiry int64, fee int64, nonce uint64, sigHex string) *HashedTransferTransaction {
	if len(commit) != 32 {
		panic("Wrong commit length")
	}
	hashCommit := [32]byte{}
	copy(hashCommit[:], commit)
	sig := &HashedTransferJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &HashedTransferTransaction{
		HashedTransfer: htlc.HashedTransfer{
			From:       from,
			To:         to,
			Value:      types.NewBigInt(value),
			HashCommit: hashCommit,
			Expiry:     expiry,
		},
		Fee:   types.NewBigInt(fee),
		Nonce: nonce,
		Sig:   sig,
	}
}

// RawHashedTransferTx returns raw bytes of a HashedTransferTransaction
func RawHashedTransferTx(from, to types.Identifier, value int64, commit []byte, expiry int64, fee int64, nonce uint64, sigHex string) []byte {
	return EncodeTx(HashedTransferTx(from, to, value, commit, expiry, fee, nonce, sigHex))
}
