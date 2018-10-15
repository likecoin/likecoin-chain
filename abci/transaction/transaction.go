package transaction

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
)

func getStatusKey(txHash []byte) []byte {
	return utils.DbTxHashKey(txHash, "status")
}

// GetStatus returns transaction status by txHash
func GetStatus(state context.IImmutableState, txHash []byte) types.TxStatus {
	_, statusBytes := state.ImmutableStateTree().Get(getStatusKey(txHash))
	return types.BytesToTxStatus(statusBytes)
}

// SetStatus set the transaction status of the given txHash
func SetStatus(
	state context.IMutableState,
	txHash []byte,
	status types.TxStatus,
) {
	state.MutableStateTree().Set(getStatusKey(txHash), status.Bytes())
}
