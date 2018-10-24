package txstatus

import (
	"bytes"
	"encoding/binary"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/utils"
)

// TxStatus is an integer representation of transaction status
type TxStatus int8

// List of TxStatus
const (
	TxStatusNotSet TxStatus = iota - 1
	TxStatusFail
	TxStatusSuccess
	TxStatusPending
)

// BytesToTxStatus converts []byte to TxStatus
func BytesToTxStatus(b []byte) TxStatus {
	if b != nil {
		var status TxStatus
		err := binary.Read(bytes.NewReader(b), binary.BigEndian, &status)
		if err == nil {
			return status
		}
	}
	return TxStatusNotSet
}

// Bytes converts TxStatus to []byte
func (status TxStatus) Bytes() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, status)
	if err == nil {
		return buf.Bytes()
	}
	return nil
}

func (status TxStatus) String() string {
	switch status {
	case TxStatusNotSet:
		return "not found"
	case TxStatusFail:
		return "fail"
	case TxStatusSuccess:
		return "success"
	case TxStatusPending:
		return "pending"
	}
	return ""
}

func getStatusKey(txHash []byte) []byte {
	return utils.DbTxHashKey(txHash, "status")
}

// GetStatus returns transaction status by txHash
func GetStatus(state context.IImmutableState, txHash []byte) TxStatus {
	_, statusBytes := state.ImmutableStateTree().Get(getStatusKey(txHash))
	return BytesToTxStatus(statusBytes)
}

// SetStatus set the transaction status of the given txHash
func SetStatus(state context.IMutableState, txHash []byte, status TxStatus) {
	state.MutableStateTree().Set(getStatusKey(txHash), status.Bytes())
}
