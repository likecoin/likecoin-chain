package txs

import (
	"github.com/likecoin/likechain/abci/context"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/crypto"

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
	cdc.RegisterConcrete(&HashedTransferTransaction{}, "github.com/likecoin/likechain/HashedTransferTransaction", nil)
	cdc.RegisterConcrete(&ClaimHashedTransferTransaction{}, "github.com/likecoin/likechain/ClaimHashedTransferTransaction", nil)
	cdc.RegisterConcrete(&SimpleTransferTransaction{}, "github.com/likecoin/likechain/SimpleTransferTransaction", nil)
	cdc.RegisterConcrete(&ContractUpdateTransaction{}, "github.com/likecoin/likechain/ContractUpdateTransaction", nil)
	cdc.RegisterInterface((*RegisterSignature)(nil), nil)
	cdc.RegisterConcrete(&RegisterJSONSignature{}, "github.com/likecoin/likechain/RegisterJSONSignature", nil)
	cdc.RegisterConcrete(&RegisterEIP712Signature{}, "github.com/likecoin/likechain/RegisterEIP712Signature", nil)
	cdc.RegisterInterface((*TransferSignature)(nil), nil)
	cdc.RegisterConcrete(&TransferJSONSignature{}, "github.com/likecoin/likechain/TransferJSONSignature", nil)
	cdc.RegisterInterface((*WithdrawSignature)(nil), nil)
	cdc.RegisterConcrete(&WithdrawJSONSignature{}, "github.com/likecoin/likechain/WithdrawJSONSignature", nil)
	cdc.RegisterConcrete(&WithdrawEIP712Signature{}, "github.com/likecoin/likechain/WithdrawEIP712Signature", nil)
	cdc.RegisterInterface((*DepositSignature)(nil), nil)
	cdc.RegisterConcrete(&DepositJSONSignature{}, "github.com/likecoin/likechain/DepositJSONSignature", nil)
	cdc.RegisterInterface((*HashedTransferSignature)(nil), nil)
	cdc.RegisterConcrete(&HashedTransferJSONSignature{}, "github.com/likecoin/likechain/HashedTransferJSONSignature", nil)
	cdc.RegisterConcrete(&HashedTransferEIP712Signature{}, "github.com/likecoin/likechain/HashedTransferEIP712Signature", nil)
	cdc.RegisterInterface((*ClaimHashedTransferSignature)(nil), nil)
	cdc.RegisterConcrete(&ClaimHashedTransferJSONSignature{}, "github.com/likecoin/likechain/ClaimHashedTransferJSONSignature", nil)
	cdc.RegisterConcrete(&ClaimHashedTransferEIP712Signature{}, "github.com/likecoin/likechain/ClaimHashedTransferEIP712Signature", nil)
	cdc.RegisterInterface((*SimpleTransferSignature)(nil), nil)
	cdc.RegisterConcrete(&SimpleTransferJSONSignature{}, "github.com/likecoin/likechain/SimpleTransferJSONSignature", nil)
	cdc.RegisterConcrete(&SimpleTransferEIP712Signature{}, "github.com/likecoin/likechain/SimpleTransferEIP712Signature", nil)
	cdc.RegisterInterface((*ContractUpdateSignature)(nil), nil)
	cdc.RegisterConcrete(&ContractUpdateJSONSignature{}, "github.com/likecoin/likechain/ContractUpdateJSONSignature", nil)
}

// EncodeTx encodes a transaction into raw bytes
func EncodeTx(tx Transaction) []byte {
	bs, err := types.AminoCodec().MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		log.
			WithField("tx", tx).
			WithError(err).
			Panic("Cannot encode transaction")
	}
	return bs
}

func recoverEthSignature(hash []byte, sig [65]byte) (*types.Address, error) {
	// Transform yellow paper V from 27/28 to 0/1
	sig[64] -= 27
	pubKeyBytes, err := crypto.Ecrecover(hash, sig[:])
	if err != nil {
		return nil, err
	}
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	ethAddr := crypto.PubkeyToAddress(*pubKey)
	addr := types.Address(ethAddr)
	return &addr, nil
}
