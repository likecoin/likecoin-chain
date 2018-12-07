package txs

import (
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWithdrawValidateFormat(t *testing.T) {
	Convey("For a withdraw transaction", t, func() {
		withdrawTx := WithdrawTx(types.Addr("0x0000000000000000000000000000000000000000"), "0x0000000000000000000000000000000000000000", types.NewBigInt(1), types.NewBigInt(0), 1, "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(withdrawTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the transaction has nil signature", func() {
			withdrawTx.Sig = nil
			Convey("The validation should fail", func() {
				So(withdrawTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil From field", func() {
			withdrawTx.From = nil
			Convey("The validation should fail", func() {
				So(withdrawTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil value", func() {
			withdrawTx.Value.Int = nil
			Convey("The validation should fail", func() {
				So(withdrawTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has negative value", func() {
			withdrawTx.Value.Int = big.NewInt(-1)
			Convey("The validation should fail", func() {
				So(withdrawTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has nil fee", func() {
			withdrawTx.Fee.Int = nil
			Convey("The validation should fail", func() {
				So(withdrawTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the transaction has negative fee", func() {
			withdrawTx.Fee.Int = big.NewInt(-1)
			Convey("The validation should fail", func() {
				So(withdrawTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestWithdrawJSONSignature(t *testing.T) {
	Convey("If a withdraw transaction with JSON signature is valid", t, func() {
		withdrawTx := WithdrawTx(types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"), "0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", types.NewBigInt(1), types.NewBigInt(0), 1, "fe1d84d34083c4e051e789ce2b4ba07fde811f6fc10c23e715b866bdc01c9a8f0cfd7995ed461d7240e3df28667f518c1d54475da9dc5b9567b60d10bc9aacaa1c")
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := withdrawTx.Sig.RecoverAddress(withdrawTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the From address of the withdraw transaction", func() {
				So(recoveredAddr, ShouldResemble, withdrawTx.From)
			})
		})
	})
}

func TestWithdrawEIP712Signature(t *testing.T) {
	Convey("If a withdraw transaction with EIP-712 signature is valid", t, func() {
		withdrawTx := &WithdrawTransaction{
			From:   types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"),
			ToAddr: *types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"),
			Value:  types.NewBigInt(1),
			Fee:    types.NewBigInt(0),
			Nonce:  1,
			Sig:    &WithdrawEIP712Signature{SigEIP712("707c63febd50a4ee228ecce1ae8b3d2c00cf32a62c59df5fe6582f44a368871d0c3f32cc29903c00fa9eb278ae85a43074a8a99e7d0d8c9c0cd14cc401fa84351b")},
		}
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := withdrawTx.Sig.RecoverAddress(withdrawTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the From address of the withdraw transaction", func() {
				So(recoveredAddr, ShouldResemble, withdrawTx.From)
			})
		})
	})
}

func TestPackWithdraw(t *testing.T) {
	Convey("In the beginning", t, func() {
		withdrawTx := &WithdrawTransaction{
			From:   types.Addr("0x1111111111111111111111111111111111111111"),
			ToAddr: *types.Addr("0x2222222222222222222222222222222222222222"),
			Value:  types.NewBigInt(1),
			Fee:    types.NewBigInt(2),
			Nonce:  3,
			Sig:    nil,
		}
		Convey("If the withdraw value and fee are valid", func() {
			Convey("pack should return the correctly packed bytes", func() {
				expectedBytes, _ := utils.Hex2Bytes("11111111111111111111111111111111111111112222222222222222222222222222222222222222000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000003")
				bs := withdrawTx.Pack()
				So(bs, ShouldResemble, expectedBytes)
			})
		})
		Convey("If the withdraw value is more than 256 bits", func() {
			withdrawTx.Value.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			Convey("pack should return nil", func() {
				bs := withdrawTx.Pack()
				So(bs, ShouldBeNil)
			})
		})
		Convey("If the withdraw fee is more than 256 bits", func() {
			withdrawTx.Fee.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			Convey("pack should return nil", func() {
				bs := withdrawTx.Pack()
				So(bs, ShouldBeNil)
			})
		})
	})
}

func TestCheckAndDeliverWithdraw(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()

		id := account.NewAccount(state, types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"))
		account.AddBalance(state, id, big.NewInt(100))

		withdrawTx := WithdrawTx(types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"), "0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", types.NewBigInt(100), types.NewBigInt(0), 1, "423e4afe124422da689f213453d3b020f3b91817c09c7009274307d7dff1cbcb3a79c8470a159102ebd002ea6ecd40911cf58c8a06badef694338c6da06d72581c")
		Convey("If a withdraw transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				withdrawTx.CheckTx(state)
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(account.FetchBalance(state, id).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					Convey("CheckTx again should return Duplicated", func() {
						// The logic of withdrawTx.DeliverTx will change the transaction object, so need to reconstruct
						// it for testing
						withdrawTx := WithdrawTx(types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"), "0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", types.NewBigInt(100), types.NewBigInt(0), 1, "423e4afe124422da689f213453d3b020f3b91817c09c7009274307d7dff1cbcb3a79c8470a159102ebd002ea6ecd40911cf58c8a06badef694338c6da06d72581c")
						r := withdrawTx.CheckTx(state)
						So(r.Code, ShouldEqual, response.WithdrawDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := withdrawTx.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.WithdrawDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
						})
					})
				})
			})
			Convey("RawWithdrawTx should return the same encoded raw tx", func() {
				rawTx := RawWithdrawTx(types.Addr("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"), "0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", types.NewBigInt(100), types.NewBigInt(0), 1, "423e4afe124422da689f213453d3b020f3b91817c09c7009274307d7dff1cbcb3a79c8470a159102ebd002ea6ecd40911cf58c8a06badef694338c6da06d72581c")
				So(rawTx, ShouldResemble, EncodeTx(withdrawTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *WithdrawTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx.From, ShouldResemble, withdrawTx.From)
					So(decodedTx.ToAddr, ShouldResemble, withdrawTx.ToAddr)
					So(decodedTx.Value, ShouldResemble, withdrawTx.Value)
					So(decodedTx.Fee.Cmp(withdrawTx.Fee.Int), ShouldBeZeroValue)
				})
			})
		})
		Convey("If a withdraw transaction from LikeChainID is valid", func() {
			withdrawTx.From = id
			withdrawTx.Sig = &WithdrawJSONSignature{Sig("0f8aef0f07ccca41d6f88015db88bc3930a8dda38473c52cbac529adfb20c37d050541843ab17f40b18d56a18be96bf55bd1ac44d21e1be33ea38aec93d795ed1c")}
			Convey("CheckTx should succeed", func() {
				withdrawTx.CheckTx(state)
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(account.FetchBalance(state, id).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
				})
			})
		})
		Convey("If a withdraw transaction has invalid format", func() {
			withdrawTx.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.WithdrawInvalidFormat.Code)
				Convey("DeliverTx should InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.WithdrawInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a withdraw transaction has unregistered sender", func() {
			withdrawTx.From = types.Addr("0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade")
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.WithdrawSenderNotRegistered.Code)
				Convey("DeliverTx should SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.WithdrawSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a withdraw transaction has invalid signature", func() {
			withdrawTx.Sig = &WithdrawJSONSignature{Sig("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.WithdrawInvalidSignature.Code)
				Convey("DeliverTx should InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.WithdrawInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a withdraw transaction has signature not from the sender", func() {
			withdrawTx.Sig = &WithdrawJSONSignature{Sig("557c7f5b1f3e445866bf3a695f0fc326d62dd3a9bc840abf69133fd116053afe27f4000027c0890846781a451c0641e01ba99e99eab2acea82e296d1ee2f7e5a1b")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.WithdrawInvalidSignature.Code)
				Convey("DeliverTx should InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.WithdrawInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a withdraw transaction has nonce exceeding next nonce", func() {
			withdrawTx.Nonce = 2
			withdrawTx.Sig = &WithdrawJSONSignature{Sig("b94501e53296c9ee970700b35d6700a5c8b5720db577af68afe28ff9666815960910d979f59ac1028509d97f05d78a4566bc8d6a491c15c5cc986141c10e72661c")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.WithdrawInvalidNonce.Code)
				Convey("DeliverTx should InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.WithdrawInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a withdraw transaction is transferring more balance than the sender has", func() {
			withdrawTx.Value.Int = big.NewInt(101)
			withdrawTx.Sig = &WithdrawJSONSignature{Sig("fffcf4232f06acdc9fef331f5bb816f131563747e25d8b5a5bf1784b64c49825274e8de719884622f20db1548a89d36ced79498d34726a0c6a65f099e139027b1b")}
			Convey("CheckTx should return NotEnoughBalance", func() {
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.WithdrawNotEnoughBalance.Code)
				Convey("DeliverTx should NotEnoughBalance", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.WithdrawNotEnoughBalance.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a withdraw transaction has not enough balance for fee", func() {
			withdrawTx.Fee.Int = big.NewInt(1)
			withdrawTx.Sig = &WithdrawJSONSignature{Sig("3c3610afc4005df69121eb4b45b5a3dd534c7415e4e2464375e2befc065b51dd44bda0d3e63eca2b1fa947c8fd56d0e270e385d503cf802a3a941c223fdc1a8e1b")}
			Convey("CheckTx should return NotEnoughBalance", func() {
				r := withdrawTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.WithdrawNotEnoughBalance.Code)
				Convey("DeliverTx should NotEnoughBalance", func() {
					txHash := utils.HashRawTx(EncodeTx(withdrawTx))
					r := withdrawTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.WithdrawNotEnoughBalance.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
	})
}
