package withdraw

import (
	"math/big"
	"testing"

	"github.com/tendermint/tendermint/libs/common"

	"github.com/likecoin/likechain/abci/types"

	. "github.com/smartystreets/goconvey/convey"
)

func wrapWithdrawTransaction(tx *types.WithdrawTransaction) *types.Transaction {
	return &types.Transaction{
		Tx: &types.Transaction_WithdrawTx{
			WithdrawTx: tx,
		},
	}
}

func withdrawTx(from, toAddr []byte, value, fee *big.Int, nonce uint64, sig []byte) *types.WithdrawTransaction {
	return &types.WithdrawTransaction{
		From: &types.Identifier{
			Id: &types.Identifier_LikeChainID{
				LikeChainID: &types.LikeChainID{
					Content: from,
				},
			},
		},
		ToAddr: &types.Address{
			Content: toAddr,
		},
		Value: &types.BigInteger{
			Content: value.Bytes(),
		},
		Fee: &types.BigInteger{
			Content: fee.Bytes(),
		},
		Nonce: nonce,
		Sig: &types.Signature{
			Version: 1,
			Content: sig,
		},
	}
}

func withdrawTxWithEthSender(from *types.Address, toAddr []byte, value, fee *big.Int, nonce uint64, sig []byte) *types.WithdrawTransaction {
	return &types.WithdrawTransaction{
		From: &types.Identifier{
			Id: &types.Identifier_Addr{
				Addr: from,
			},
		},
		ToAddr: &types.Address{
			Content: toAddr,
		},
		Value: &types.BigInteger{
			Content: value.Bytes(),
		},
		Fee: &types.BigInteger{
			Content: fee.Bytes(),
		},
		Nonce: nonce,
		Sig: &types.Signature{
			Version: 1,
			Content: sig,
		},
	}
}

func TestValidateWithdrawTransactionFormat(t *testing.T) {
	zero := big.NewInt(0)
	one := big.NewInt(1)

	Convey("Given a Withdraw transaction with valid format and LikeChainID sender", t, func() {
		tx := withdrawTx(common.RandBytes(20), common.RandBytes(20), one, zero, 0, common.RandBytes(65))
		Convey("The transaction should pass the validation", func() {
			So(validateWithdrawTransactionFormat(tx), ShouldBeTrue)
		})
	})

	Convey("Given a Withdraw transaction with valid format and Ethereum address sender", t, func() {
		ethAddr := &types.Address{Content: common.RandBytes(20)}
		tx := withdrawTxWithEthSender(ethAddr, common.RandBytes(20), one, zero, 0, common.RandBytes(65))
		Convey("The transaction should pass the validation", func() {
			So(validateWithdrawTransactionFormat(tx), ShouldBeTrue)
		})
	})

	Convey("Given a Withdraw transaction with invalid LikeChainID sender", t, func() {
		Convey("The transaction should not pass the validation", func() {
			tx := withdrawTx(nil, common.RandBytes(20), one, zero, 0, common.RandBytes(65))
			So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
			for n := 1; n <= 40; n++ {
				if n == 20 {
					continue
				}
				tx = withdrawTx(common.RandBytes(n), common.RandBytes(20), one, zero, 0, common.RandBytes(65))
				So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
			}
		})
	})

	Convey("Given a Withdraw transaction with invalid Ethereum address sender", t, func() {
		Convey("The transaction should not pass the validation", func() {
			tx := withdrawTxWithEthSender(nil, common.RandBytes(20), one, zero, 0, common.RandBytes(65))
			So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
			for n := 1; n <= 40; n++ {
				if n == 20 {
					continue
				}
				ethAddr := &types.Address{Content: common.RandBytes(n)}
				tx = withdrawTxWithEthSender(ethAddr, common.RandBytes(20), one, zero, 0, common.RandBytes(65))
				So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
			}
		})
	})

	Convey("Given a Withdraw transaction with invalid value", t, func() {
		Convey("The transaction should not pass the validation", func() {
			tx := withdrawTx(common.RandBytes(20), common.RandBytes(20), zero, zero, 0, common.RandBytes(65))
			So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)

			value := big.NewInt(2)
			value.Exp(value, big.NewInt(256), nil)
			tx = withdrawTx(common.RandBytes(20), common.RandBytes(20), value, zero, 0, common.RandBytes(65))
			So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
		})
	})

	Convey("Given a Withdraw transaction with maximum possible value", t, func() {
		value := big.NewInt(2)
		value.Exp(value, big.NewInt(256), nil)
		value.Sub(value, one)
		tx := withdrawTx(common.RandBytes(20), common.RandBytes(20), value, zero, 0, common.RandBytes(65))
		Convey("The transaction should pass the validation", func() {
			So(validateWithdrawTransactionFormat(tx), ShouldBeTrue)
		})
	})

	Convey("Given a Withdraw transaction with invalid fee", t, func() {
		Convey("The transaction should not pass the validation", func() {
			fee := big.NewInt(2)
			fee.Exp(fee, big.NewInt(256), nil)
			tx := withdrawTx(common.RandBytes(20), common.RandBytes(20), one, fee, 0, common.RandBytes(65))
			So(validateWithdrawTransactionFormat(tx), ShouldBeFalse)
		})
	})

	Convey("Given a Withdraw transaction with maximum possible fee", t, func() {
		fee := big.NewInt(2)
		fee.Exp(fee, big.NewInt(256), nil)
		fee.Sub(fee, one)
		tx := withdrawTx(common.RandBytes(20), common.RandBytes(20), one, fee, 0, common.RandBytes(65))
		Convey("The transaction should pass the validation", func() {
			So(validateWithdrawTransactionFormat(tx), ShouldBeTrue)
		})
	})
}
