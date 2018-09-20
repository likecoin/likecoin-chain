package types

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/sirupsen/logrus"
)

var log = logger.L

// NewAddressFromHex creates Address from hex string
func NewAddressFromHex(hex string) *Address {
	return &Address{Content: common.FromHex(hex)}
}

// NewZeroAddress creates Address from all zeros
func NewZeroAddress() *Address {
	return NewAddressFromHex("0x0000000000000000000000000000000000000000")
}

// IsValidFormat checks the length of the address
func (addr *Address) IsValidFormat() bool {
	return len(addr.Content) == 20
}

// ToIdentifier converts Address to Identifier
func (addr *Address) ToIdentifier() *Identifier {
	return &Identifier{
		Id: &Identifier_Addr{
			Addr: addr,
		},
	}
}

// ToEthereum converts Address struct to Ethereum address
func (addr *Address) ToEthereum() common.Address {
	addrBytes := addr.Content
	return common.BytesToAddress(addrBytes)
}

// ToHex converts Address to hex string
func (addr *Address) ToHex() string {
	if addr != nil {
		return addr.ToEthereum().Hex()
	}
	return ""
}

// NewSignatureFromHex creates Signature from hex string
func NewSignatureFromHex(hex string) *Signature {
	return &Signature{
		Content: common.FromHex(hex),
		Version: 1,
	}
}

// NewZeroSignature creates a Signature with all zeros
func NewZeroSignature() *Signature {
	return NewSignatureFromHex("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
}

// ToHex converts Signature to hex string
func (sig *Signature) ToHex() string {
	return common.ToHex(sig.Content)
}

// IsValidFormat checks the signature version and its length
func (sig *Signature) IsValidFormat() bool {
	switch sig.Version {
	case 1:
		content := sig.Content
		if len(content) != 65 {
			return false
		}
		return true
	default:
		return false
	}
}

// NewBigInteger creates BigInteger from int64
func NewBigInteger(s string) *BigInteger {
	i := new(big.Int)
	i.SetString(s, 10)
	return &BigInteger{Content: i.Bytes()}
}

// ToBigInt converts BigInteger struct to big Int
func (i *BigInteger) ToBigInt() *big.Int {
	bigInt := new(big.Int)
	return bigInt.SetBytes(i.Content)
}

// ToString converts BigInteger struct to string
func (i *BigInteger) ToString() string {
	bigInt := new(big.Int)
	bigInt.SetBytes(i.Content)
	return bigInt.String()
}

// IsValidFormat checks Identifier format
func (id *Identifier) IsValidFormat() bool {
	return (id.GetLikeChainID() != nil && len(id.GetLikeChainID().Content) > 0) ||
		(id.GetAddr() != nil && len(id.GetAddr().Content) > 0)
}

// ToString converts Identifier to hex string
func (id *Identifier) ToString() string {
	if likeChainID := id.GetLikeChainID(); likeChainID != nil {
		return likeChainID.ToString()
	} else if addr := id.GetAddr(); addr != nil {
		return strings.ToLower(addr.ToHex())
	}
	return ""
}

// NewLikeChainID creates a LikeChainID from bytes
func NewLikeChainID(content []byte) *LikeChainID {
	return &LikeChainID{Content: content}
}

// NewLikeChainIDFromString creates a LikeChainID from string
func NewLikeChainIDFromString(s string) (*LikeChainID, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return NewLikeChainID(b), nil
}

// ToString converts LikeChain ID to base64-encoded strings
func (id *LikeChainID) ToString() string {
	return base64.StdEncoding.EncodeToString(id.Content)
}

// ToIdentifier converts LikeChainID to Identifier
func (id *LikeChainID) ToIdentifier() *Identifier {
	return &Identifier{
		Id: &Identifier_LikeChainID{
			LikeChainID: id,
		},
	}
}

// NewTransferTarget creates a new TransferTarget
func NewTransferTarget(id *Identifier, value string, remark string) *TransferTransaction_TransferTarget {
	return &TransferTransaction_TransferTarget{
		To:     id,
		Value:  NewBigInteger(value),
		Remark: []byte(remark),
	}
}

// IsValidFormat checks the format of TransferTarget ensuring the `value` and
// `to` are correct
func (t *TransferTransaction_TransferTarget) IsValidFormat() bool {
	value := t.Value
	to := t.To
	if value == nil ||
		to == nil {
		return false
	}
	return value.ToBigInt().Cmp(big.NewInt(0)) >= 0 &&
		to.IsValidFormat()
}

var sigPrefix = []byte("\x19Ethereum Signed Message:\n")

// generateSigningMessageHash wraps a JSON in map in follwing format
// `\x19Ethereum Signed Message:\n" + len(message) + message`
// and return Keccak256 hash
func generateSigningMessageHash(jsonMap map[string]interface{}) (hash []byte) {
	msg, err := json.Marshal(jsonMap)
	if err == nil {
		hashingMsg := []byte(fmt.Sprintf("%s%d%s", sigPrefix, len(msg), msg))
		hash = crypto.Keccak256(hashingMsg)
	}
	return hash
}

// GenerateSigningMessageHash generates a signature from a RegisterTx
func (tx *RegisterTransaction) GenerateSigningMessageHash() []byte {
	return generateSigningMessageHash(map[string]interface{}{
		"addr": strings.ToLower(tx.Addr.ToEthereum().Hex()),
	})
}

// ToString converts RegisterTransaction to formatted string
func (tx *RegisterTransaction) ToString() string {
	return fmt.Sprintf(
		"<Addr: %s, Sig: %s>",
		tx.Addr.ToHex(),
		tx.Sig.ToHex(),
	)
}

// ToTransaction wraps RegisterTransaction into Transaction
func (tx *RegisterTransaction) ToTransaction() *Transaction {
	return &Transaction{
		Tx: &Transaction_RegisterTx{
			RegisterTx: tx,
		},
	}
}

// GenerateSigningMessageHash generates a signature from a TransferTx
func (tx *TransferTransaction) GenerateSigningMessageHash() (hash []byte) {
	to := make([]map[string]interface{}, len(tx.ToList))
	for i, target := range tx.ToList {
		to[i] = map[string]interface{}{
			"identity": target.To.ToString(),
			"remark":   string(target.Remark),
			"value":    target.Value.ToString(),
		}
	}
	return generateSigningMessageHash(map[string]interface{}{
		"fee":      tx.Fee.ToString(),
		"identity": tx.From.ToString(),
		"nonce":    tx.Nonce,
		"to":       to,
	})
}

// GenerateSigningMessageHash generates a signature from a WithdrawTx
func (tx *WithdrawTransaction) GenerateSigningMessageHash() (hash []byte) {
	return generateSigningMessageHash(map[string]interface{}{
		"fee":      tx.Fee.ToString(),
		"identity": tx.From.ToString(),
		"nonce":    tx.Nonce,
		"to_addr":  strings.ToLower(tx.ToAddr.ToHex()),
		"value":    tx.Value.ToString(),
	})
}

func bigIntToUint256Bytes(n *big.Int) []byte {
	result := make([]byte, 32)
	b := n.Bytes()
	l := len(b)
	if l > 32 {
		log.WithFields(logrus.Fields{
			"number": n.String(),
		}).Panic("Transforming big.Int to uint256 bytes with overflowing value")
	}
	copy(result[32-l:], b)
	return result
}

// Pack generates a packed version for the WithdrawTx, used when storing and proving withdraw records in withdraw tree
// Packing format: 20-byte From, 20-byte ToAddr, 32-byte Value, 32-Byte Fee, 8-byte Nonce
// All numbers are in big endian
func (tx *WithdrawTransaction) Pack() []byte {
	buf := new(bytes.Buffer)
	buf.Write(tx.From.GetLikeChainID().Content)
	buf.Write(tx.ToAddr.Content)
	buf.Write(bigIntToUint256Bytes(tx.Value.ToBigInt()))
	buf.Write(bigIntToUint256Bytes(tx.Fee.ToBigInt()))
	binary.Write(buf, binary.BigEndian, tx.Nonce)
	return buf.Bytes()
}

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
	case TxStatusFail:
		return "fail"
	case TxStatusSuccess:
		return "success"
	case TxStatusPending:
		return "pending"
	}
	return ""
}
