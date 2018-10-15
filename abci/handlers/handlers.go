package handlers

import (
	"fmt"

	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/handlers/table"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/transaction"
	"github.com/likecoin/likechain/abci/types"
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
		return abci.ResponseCheckTx{
			Code: 10,
			Log:  fmt.Sprintf("No CheckTx handler for type %v", _type),
		}
	}
	return handler(state, tx).ToResponseCheckTx()
}

// DeliverTx handles DeliverTx
func DeliverTx(state context.IMutableState, tx *types.Transaction, txHash []byte) abci.ResponseDeliverTx {
	_type, handler, exist := table.GetDeliverTxHandlerFromTx(tx)
	if !exist {
		log.WithField("type", _type).Debug("Deliver handler not exist")
		return abci.ResponseDeliverTx{
			Code: 11,
			Log:  fmt.Sprintf("No DeliverTx handler for type %v", _type),
		}
	}
	r := handler(state, tx, txHash)
	oldStatus := transaction.GetStatus(state, txHash)
	if oldStatus == types.TxStatusNotSet || oldStatus == types.TxStatusPending {
		transaction.SetStatus(state, txHash, r.Status)
	}
	return r.ToResponseDeliverTx()
}
