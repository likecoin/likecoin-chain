package txs

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"strconv"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// WithdrawTransaction represents a Withdraw transaction
type WithdrawTransaction struct {
	From   types.Identifier
	ToAddr types.Address
	Value  types.BigInt
	Fee    types.BigInt
	Nonce  uint64
	Sig    WithdrawSignature
}

// Pack returns the packed form of the withdraw transaction, which is used for withdraw in Ethereum contract
// Packing format: 20-byte From, 20-byte ToAddr, 32-byte Value, 32-Byte Fee, 8-byte Nonce
// All numbers are in big endian
func (tx *WithdrawTransaction) Pack() []byte {
	buf := new(bytes.Buffer)
	buf.Write(tx.From.Bytes())
	buf.Write(tx.ToAddr.Bytes())
	valueBytes := tx.Value.ToUint256Bytes()
	if valueBytes == nil {
		return nil
	}
	buf.Write(valueBytes)
	feeBytes := tx.Fee.ToUint256Bytes()
	if feeBytes == nil {
		return nil
	}
	buf.Write(feeBytes)
	binary.Write(buf, binary.BigEndian, tx.Nonce)
	return buf.Bytes()
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *WithdrawTransaction) ValidateFormat() bool {
	if tx.From == nil || tx.Value.Int == nil || tx.Fee.Int == nil || tx.Sig == nil {
		return false
	}
	if !tx.Value.IsWithinRange() {
		return false
	}
	if !tx.Fee.IsWithinRange() {
		return false
	}
	return true
}

func (tx *WithdrawTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.WithdrawInvalidFormat.Info)
		return response.WithdrawInvalidFormat, nil
	}

	senderID = account.IdentifierToLikeChainID(state, tx.From)
	if senderID == nil {
		logTx(tx).Info(response.WithdrawSenderNotRegistered.Info)
		return response.WithdrawSenderNotRegistered, nil
	}

	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.WithdrawInvalidSignature.Info)
		return response.WithdrawInvalidSignature, senderID
	}

	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.WithdrawInvalidNonce.Info)
		return response.WithdrawInvalidNonce, senderID
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.WithdrawDuplicated.Info)
		return response.WithdrawDuplicated, senderID
	}

	senderBalance := account.FetchBalance(state, tx.From)

	total := new(big.Int).Add(tx.Value.Int, tx.Fee.Int)
	if senderBalance.Cmp(total) < 0 {
		logTx(tx).Info(response.WithdrawNotEnoughBalance.Info)
		return response.WithdrawNotEnoughBalance, senderID
	}

	return response.Success, senderID
}

// CheckTx checks the transaction to see if it should be executed
func (tx *WithdrawTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *WithdrawTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID := tx.checkTx(state)
	if checkTxRes.Code != response.Success.Code {
		if checkTxRes.ShouldIncrementNonce {
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	tx.From = senderID
	account.IncrementNextNonce(state, senderID)

	// The fee is distributed to the one who do withdraw on Ethereum
	total := new(big.Int).Add(tx.Value.Int, tx.Fee.Int)
	account.MinusBalance(state, senderID, total)
	packedTx := tx.Pack()

	withdrawTree := state.MutableWithdrawTree()
	withdrawTree.Set(crypto.Sha256(packedTx), []byte{1})

	height := state.GetHeight() + 1
	log.
		WithField("height", height).
		WithField("packed_tx", cmn.HexBytes(packedTx)).
		Debug("Saved Withdraw proof")

	return response.Success.Merge(response.R{
		Tags: []cmn.KVPair{
			{
				Key:   []byte("withdraw.height"),
				Value: []byte(strconv.FormatInt(height, 10)),
			},
		},
		Data: packedTx,
	})
}

// WithdrawTx returns a WithdrawTransaction
func WithdrawTx(from types.Identifier, toAddrHex string, value, fee types.BigInt, nonce uint64, sigHex string) *WithdrawTransaction {
	toAddr := *types.Addr(toAddrHex)
	sig := &WithdrawJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &WithdrawTransaction{
		From:   from,
		ToAddr: toAddr,
		Value:  value,
		Fee:    fee,
		Nonce:  nonce,
		Sig:    sig,
	}
}

// RawWithdrawTx returns raw bytes of a WithdrawTransaction
func RawWithdrawTx(from types.Identifier, toAddrHex string, value, fee types.BigInt, nonce uint64, sigHex string) []byte {
	return EncodeTx(WithdrawTx(from, toAddrHex, value, fee, nonce, sigHex))
}
