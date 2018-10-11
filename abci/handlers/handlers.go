package handlers

import (
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers/table"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"
	abci "github.com/tendermint/tendermint/abci/types"

	// Init handlers
	_ "github.com/likecoin/likechain/abci/handlers/register"
	_ "github.com/likecoin/likechain/abci/handlers/transfer"
	_ "github.com/likecoin/likechain/abci/handlers/withdraw"
)

var log = logger.L

// CheckTx handles CheckTx
func CheckTx(state context.IImmutableState, tx *types.Transaction) abci.ResponseCheckTx {
	_type, handler, exist := table.GetCheckTxHandlerFromTx(tx)
	if !exist {
		log.WithField("type", _type).Debug("CheckTx handler not exist")
		return abci.ResponseCheckTx{} // TODO
	}
	return handler(state, tx).ToResponseCheckTx()
}

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

// DeliverTx handles DeliverTx
func DeliverTx(state context.IMutableState, tx *types.Transaction, txHash []byte) abci.ResponseDeliverTx {
	_type, handler, exist := table.GetDeliverTxHandlerFromTx(tx)
	if !exist {
		log.WithField("type", _type).Debug("Deliver handler not exist")
		return abci.ResponseDeliverTx{} // TODO
	}
	r := handler(state, tx, txHash)
	SetStatus(state, txHash, r.Status)
	return r.ToResponseDeliverTx()
}
