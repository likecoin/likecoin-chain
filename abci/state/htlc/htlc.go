package htlc

import (
	"bytes"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"

	"github.com/tendermint/tendermint/crypto"
)

// IsExpired returns whether the HashedTransfer is expired in current context
func (ht *HashedTransfer) IsExpired(state context.IImmutableState) bool {
	now := state.GetBlockTime()
	return now >= ht.Expiry
}

// CheckCreateHashedTransfer checks if a HashedTransfer is valid to be craeted in current context
func CheckCreateHashedTransfer(state context.IImmutableState, ht *HashedTransfer) response.R {
	if ht.IsExpired(state) {
		return response.HashedTransferInvalidExpiry
	}
	return response.Success
}

// CheckClaimHashedTransfer checks if a HashedTransfer is valid to be claimed in current context
func CheckClaimHashedTransfer(state context.IImmutableState, ht *HashedTransfer, secret []byte) response.R {
	if ht.IsExpired(state) {
		return response.ClaimHashedTransferExpired
	}
	hash := crypto.Sha256(secret)
	if bytes.Compare(hash, ht.HashCommit[:]) != 0 {
		return response.ClaimHashedTransferInvalidSecret
	}
	return response.Success
}

// ClaimHashedTransfer executes a HashedTransfer and send the transfer value to the receiver
func ClaimHashedTransfer(state context.IMutableState, ht *HashedTransfer, txHash []byte) {
	account.AddBalance(state, ht.To, ht.Value.Int)
	RemoveHashedTransfer(state, txHash)
}

// CheckRevokeHashedTransfer checks if a HashedTransfer is valid to be revoked in current context
func CheckRevokeHashedTransfer(state context.IImmutableState, ht *HashedTransfer) response.R {
	if !ht.IsExpired(state) {
		return response.ClaimHashedTransferNotYetExpired
	}
	return response.Success
}

// RevokeHashedTransfer cancels a HashedTransfer and send back the transfer value to the sender
func RevokeHashedTransfer(state context.IMutableState, ht *HashedTransfer, txHash []byte) {
	account.AddBalance(state, ht.From, ht.Value.Int)
	RemoveHashedTransfer(state, txHash)
}
