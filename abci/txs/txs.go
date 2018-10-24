package txs

import (
	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	"github.com/sirupsen/logrus"
)

var log = logger.L

// Transaction represents a Tendermint transaction
type Transaction interface {
	ValidateFormat() bool
	CheckTx(state context.IImmutableState) response.R
	DeliverTx(state context.IMutableState, txHash []byte) response.R
}

func logTx(tx Transaction) *logrus.Entry {
	return log.WithField("tx", tx)
}

func init() {
	cdc := types.AminoCodec()
	cdc.RegisterInterface((*Transaction)(nil), nil)
	cdc.RegisterConcrete(&RegisterTransaction{}, "github.com/likecoin/likechain/RegisterTransaction", nil)
	cdc.RegisterConcrete(&TransferTransaction{}, "github.com/likecoin/likechain/TransferTransaction", nil)
	cdc.RegisterConcrete(&WithdrawTransaction{}, "github.com/likecoin/likechain/WithdrawTransaction", nil)
	cdc.RegisterConcrete(&DepositTransaction{}, "github.com/likecoin/likechain/DepositTransaction", nil)
	cdc.RegisterInterface((*TransferSignature)(nil), nil)
	cdc.RegisterConcrete(&RegisterJSONSignature{}, "github.com/likecoin/likechain/RegisterJSONSignature", nil)
	cdc.RegisterInterface((*RegisterSignature)(nil), nil)
	cdc.RegisterConcrete(&TransferJSONSignature{}, "github.com/likecoin/likechain/TransferJSONSignature", nil)
	cdc.RegisterInterface((*WithdrawSignature)(nil), nil)
	cdc.RegisterConcrete(&WithdrawJSONSignature{}, "github.com/likecoin/likechain/WithdrawJSONSignature", nil)
}

// EncodeTx encodes a transaction into raw bytes
func EncodeTx(tx Transaction) []byte {
	bs, err := types.AminoCodec().MarshalBinary(tx)
	if err != nil {
		panic(err)
	}
	return bs
}
