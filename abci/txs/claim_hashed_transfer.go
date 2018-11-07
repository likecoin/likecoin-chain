package txs

import (
	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/htlc"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// ClaimHashedTransferTransaction represents a claim to a previous Hashed TimeLock Transfer transaction
type ClaimHashedTransferTransaction struct {
	From       types.Identifier
	HTLCTxHash []byte
	Secret     []byte
	Nonce      uint64
	Sig        ClaimHashedTransferSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *ClaimHashedTransferTransaction) ValidateFormat() bool {
	if tx.From == nil || tx.Sig == nil {
		return false
	}
	if len(tx.HTLCTxHash) != tmhash.Size {
		return false
	}
	secretLen := len(tx.Secret)
	if secretLen != 0 && secretLen != 32 {
		return false
	}
	return true
}

func (tx *ClaimHashedTransferTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID, ht *htlc.HashedTransfer,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.ClaimHashedTransferInvalidFormat.Info)
		return response.ClaimHashedTransferInvalidFormat, nil, nil
	}
	senderID = account.IdentifierToLikeChainID(state, tx.From)
	if senderID == nil {
		logTx(tx).Info(response.ClaimHashedTransferSenderNotRegistered.Info)
		return response.ClaimHashedTransferSenderNotRegistered, senderID, nil
	}
	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.ClaimHashedTransferInvalidSignature.Info)
		return response.ClaimHashedTransferInvalidSignature, senderID, nil
	}
	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.ClaimHashedTransferInvalidNonce.Info)
		return response.ClaimHashedTransferInvalidNonce, senderID, nil
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.ClaimHashedTransferDuplicated.Info)
		return response.ClaimHashedTransferDuplicated, senderID, nil
	}

	ht = htlc.GetHashedTransfer(state, tx.HTLCTxHash)
	if ht == nil {
		return response.ClaimHashedTransferTxNotExist, senderID, ht
	}
	htlcReceiverID := account.IdentifierToLikeChainID(state, ht.To)
	if senderID.Equals(htlcReceiverID) {
		return htlc.CheckClaimHashedTransfer(state, ht, tx.Secret), senderID, ht
	}
	htlcSenderID := account.IdentifierToLikeChainID(state, ht.From)
	if senderID.Equals(htlcSenderID) {
		return htlc.CheckRevokeHashedTransfer(state, ht), senderID, ht
	}
	return response.ClaimHashedTransferInvalidSender, senderID, ht
}

// CheckTx checks the transaction to see if it should be executed
func (tx *ClaimHashedTransferTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *ClaimHashedTransferTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID, ht := tx.checkTx(state)
	if checkTxRes.Code != 0 {
		switch checkTxRes.Code {
		case response.ClaimHashedTransferTxNotExist.Code:
			fallthrough
		case response.ClaimHashedTransferExpired.Code:
			fallthrough
		case response.ClaimHashedTransferInvalidSecret.Code:
			fallthrough
		case response.ClaimHashedTransferNotYetExpired.Code:
			fallthrough
		case response.ClaimHashedTransferInvalidSender.Code:
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	account.IncrementNextNonce(state, senderID)

	var tags []cmn.KVPair

	htlcReceiverID := account.IdentifierToLikeChainID(state, ht.To)
	if senderID.Equals(htlcReceiverID) {
		htlc.ClaimHashedTransfer(state, ht, tx.HTLCTxHash)
		tags = []cmn.KVPair{
			{
				Key:   []byte("claim_hashed_transfer.htlc_tx_hash"),
				Value: tx.HTLCTxHash,
			},
			{
				Key:   []byte("claim_hashed_transfer.secret"),
				Value: tx.Secret,
			},
		}
	} else {
		htlcSenderID := account.IdentifierToLikeChainID(state, ht.From)
		if senderID.Equals(htlcSenderID) {
			htlc.RevokeHashedTransfer(state, ht, tx.HTLCTxHash)
		} else {
			logTx(tx).
				WithField("from", tx.From.String()).
				WithField("htlc_tx_hash", cmn.HexBytes(tx.HTLCTxHash)).
				Panic("tx.From is neither sender or receiver of the HashedTransferTransaction")
		}
	}

	txstatus.SetStatus(state, tx.HTLCTxHash, txstatus.TxStatusSuccess)

	return response.Success.Merge(response.R{
		Tags: tags,
	})
}

// ClaimHashedTransferTx returns a ClaimHashedTransferTransaction
func ClaimHashedTransferTx(from types.Identifier, htlcTxHash []byte, secret []byte, nonce uint64, sigHex string) *ClaimHashedTransferTransaction {
	if len(secret) != 32 && len(secret) != 0 {
		panic("Wrong secret length")
	}
	if len(htlcTxHash) != tmhash.Size {
		panic("Wrong htlcTxHash length")
	}
	sig := &ClaimHashedTransferJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &ClaimHashedTransferTransaction{
		From:       from,
		HTLCTxHash: htlcTxHash,
		Secret:     secret,
		Nonce:      nonce,
		Sig:        sig,
	}
}

// RawClaimHashedTransferTx returns raw bytes of a ClaimHashedTransferTransaction
func RawClaimHashedTransferTx(from types.Identifier, htlcTxHash []byte, secret []byte, nonce uint64, sigHex string) []byte {
	return EncodeTx(ClaimHashedTransferTx(from, htlcTxHash, secret, nonce, sigHex))
}
