package txs

import (
	"encoding/base64"
	"math/big"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

// TransferOutput represents one of the outputs in a Transfer transaction
type TransferOutput struct {
	To     types.Identifier
	Value  types.BigInt
	Remark TransferRemark
}

// Validate checks if a transfer output is valid
func (output TransferOutput) Validate() bool {
	return output.To != nil && output.Value.Int != nil && output.Value.IsWithinRange() && output.Remark.Validate()
}

// TransferRemark is the remark of a transfer output
type TransferRemark []byte

func (remark TransferRemark) String() string {
	return base64.StdEncoding.EncodeToString(remark)
}

// Validate checks if the remark is valid
func (remark TransferRemark) Validate() bool {
	return len(remark) <= 4096
}

// TransferTransaction represents a Transfer transaction
type TransferTransaction struct {
	From    types.Identifier
	Outputs []TransferOutput
	Nonce   uint64
	Fee     types.BigInt
	Sig     TransferSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *TransferTransaction) ValidateFormat() bool {
	if tx.From == nil || tx.Fee.Int == nil || tx.Sig == nil {
		return false
	}
	if len(tx.Outputs) == 0 {
		return false
	}
	if !tx.Fee.IsWithinRange() {
		return false
	}
	for _, output := range tx.Outputs {
		if !output.Validate() {
			return false
		}
	}
	return true
}

func (tx *TransferTransaction) checkTx(state context.IImmutableState) (
	r response.R, senderID *types.LikeChainID, toIdens map[types.Identifier]*big.Int,
) {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.TransferInvalidFormat.Info)
		return response.TransferInvalidFormat, nil, nil
	}
	senderID = account.IdentifierToLikeChainID(state, tx.From)
	if senderID == nil {
		logTx(tx).Info(response.TransferSenderNotRegistered.Info)
		return response.TransferSenderNotRegistered, senderID, nil
	}
	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !account.IsLikeChainIDHasAddress(state, senderID, addr) {
		logTx(tx).
			WithField("recovered_addr", addr).
			WithError(err).
			Info(response.TransferInvalidSignature.Info)
		return response.TransferInvalidSignature, senderID, nil
	}
	nextNonce := account.FetchNextNonce(state, senderID)
	if tx.Nonce > nextNonce {
		logTx(tx).Info(response.TransferInvalidNonce.Info)
		return response.TransferInvalidNonce, senderID, nil
	} else if tx.Nonce < nextNonce {
		logTx(tx).Info(response.TransferDuplicated.Info)
		return response.TransferDuplicated, senderID, nil
	}
	senderBalance := account.FetchBalance(state, senderID)
	toIdens = make(map[types.Identifier]*big.Int)
	total := new(big.Int).Set(tx.Fee.Int)
	for _, output := range tx.Outputs {
		if _, ok := output.To.(*types.LikeChainID); ok {
			toID := account.IdentifierToLikeChainID(state, output.To)
			if toID == nil {
				logTx(tx).
					WithField("to", output.To).
					Info(response.TransferInvalidReceiver.Info)
				return response.TransferInvalidReceiver, senderID, nil
			}
			toIdens[toID] = output.Value.Int
		} else {
			toIdens[output.To] = output.Value.Int
		}
		total.Add(total, output.Value.Int)
	}
	if senderBalance.Cmp(total) < 0 {
		logTx(tx).
			WithField("total", total.String()).
			WithField("balance", senderBalance.String()).
			Info(response.TransferNotEnoughBalance.Info)
		return response.TransferNotEnoughBalance, senderID, nil
	}
	return response.Success, senderID, toIdens
}

// CheckTx checks the transaction to see if it should be executed
func (tx *TransferTransaction) CheckTx(state context.IImmutableState) response.R {
	r, _, _ := tx.checkTx(state)
	return r
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *TransferTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxRes, senderID, toIdens := tx.checkTx(state)
	if checkTxRes.Code != 0 {
		if checkTxRes.ShouldIncrementNonce {
			account.IncrementNextNonce(state, senderID)
		}
		return checkTxRes
	}

	account.IncrementNextNonce(state, senderID)

	total := new(big.Int).Set(tx.Fee.Int)
	for identity, value := range toIdens {
		account.AddBalance(state, identity, value)
		total.Add(total, value)
	}
	account.MinusBalance(state, senderID, total)

	return response.Success
}

// TransferTx returns a TransferTransaction
func TransferTx(from types.Identifier, outputs []TransferOutput, fee types.BigInt, nonce uint64, sigHex string) *TransferTransaction {
	sig := &TransferJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &TransferTransaction{
		From:    from,
		Outputs: outputs,
		Fee:     fee,
		Nonce:   nonce,
		Sig:     sig,
	}
}

// RawTransferTx returns raw bytes of a TransferTransaction
func RawTransferTx(from types.Identifier, outputs []TransferOutput, fee types.BigInt, nonce uint64, sigHex string) []byte {
	return EncodeTx(TransferTx(from, outputs, fee, nonce, sigHex))
}
