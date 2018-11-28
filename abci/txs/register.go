package txs

import (
	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"
)

// RegisterTransaction represents a Register transaction
type RegisterTransaction struct {
	Addr types.Address
	Sig  RegisterSignature
}

// ValidateFormat checks if a transaction has invalid format, e.g. nil fields, negative transfer amounts
func (tx *RegisterTransaction) ValidateFormat() bool {
	return tx.Sig != nil
}

// CheckTx checks the transaction to see if it should be executed
func (tx *RegisterTransaction) CheckTx(state context.IImmutableState) response.R {
	if !tx.ValidateFormat() {
		logTx(tx).Info(response.RegisterInvalidFormat.Info)
		return response.RegisterInvalidFormat
	}
	addr, err := tx.Sig.RecoverAddress(tx)
	if err != nil || !addr.Equals(&tx.Addr) {
		logTx(tx).
			WithError(err).
			Info(response.RegisterInvalidSignature.Info)
		return response.RegisterInvalidSignature
	}
	if account.IsAddressRegistered(state, addr) {
		logTx(tx).Info(response.RegisterDuplicated.Info)
		return response.RegisterDuplicated
	}
	return response.Success
}

// DeliverTx checks the transaction to see if it should be executed
func (tx *RegisterTransaction) DeliverTx(state context.IMutableState, txHash []byte) response.R {
	checkTxResult := tx.CheckTx(state)
	if checkTxResult.Code != response.Success.Code {
		return checkTxResult
	}

	id := account.NewAccount(state, &tx.Addr)
	return response.Success.Merge(response.R{
		Data: id.Bytes(),
	})
}

// RegisterTx returns a RegisterTransaction
func RegisterTx(addrHex, sigHex string) *RegisterTransaction {
	addr := *types.Addr(addrHex)
	sig := &RegisterJSONSignature{
		JSONSignature: Sig(sigHex),
	}
	return &RegisterTransaction{
		Addr: addr,
		Sig:  sig,
	}
}

// RawRegisterTx returns raw bytes of a RegisterTransaction
func RawRegisterTx(addrHex, sigHex string) []byte {
	return EncodeTx(RegisterTx(addrHex, sigHex))
}
