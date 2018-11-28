package txs

import (
	"math/big"
	"testing"

	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	. "github.com/likecoin/likechain/abci/fixture"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/htlc"
	"github.com/likecoin/likechain/abci/txstatus"
	"github.com/likecoin/likechain/abci/types"
	"github.com/likecoin/likechain/abci/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashedTrasnferValidateFormat(t *testing.T) {
	Convey("For a HashedTransfer transaction", t, func() {
		commit := make([]byte, 32)
		hashedTransferTx := HashedTransferTx(Alice.Address, Bob.ID, 1, commit, 10, 0, 1, "f8aa5170bab5747a2216f544897adf4f2a3643cd0d6e265d663a416cc5d4a13034d090d9cd48836b2fe0c67825dacd11426f3fb737820665236ca468eb8fb3111c")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(hashedTransferTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the signature is nil", func() {
			hashedTransferTx.Sig = nil
			Convey("The validation should fail", func() {
				So(hashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the fee is nil", func() {
			hashedTransferTx.Fee.Int = nil
			Convey("The validation should fail", func() {
				So(hashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the fee has negative value", func() {
			hashedTransferTx.Fee = types.NewBigInt(-1)
			Convey("The validation should fail", func() {
				So(hashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the fee has value 2^256", func() {
			hashedTransferTx.Fee.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			Convey("The validation should fail", func() {
				So(hashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the fee has value 2^256-1", func() {
			hashedTransferTx.Fee.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
			hashedTransferTx.Fee.Int.Sub(hashedTransferTx.Fee.Int, big.NewInt(1))
			Convey("The validation should success", func() {
				So(hashedTransferTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the HashedTransfer in the tx has invalid format", func() {
			hashedTransferTx.HashedTransfer.To = nil
			Convey("The validation should fail", func() {
				So(hashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestHashedTransferJSONSignature(t *testing.T) {
	Convey("If a HashedTransfer transaction with JSON signature is valid", t, func() {
		commit := make([]byte, 32)
		hashedTransferTx := HashedTransferTx(Alice.Address, Bob.ID, 1, commit, 10, 0, 1, "f8aa5170bab5747a2216f544897adf4f2a3643cd0d6e265d663a416cc5d4a13034d090d9cd48836b2fe0c67825dacd11426f3fb737820665236ca468eb8fb3111c")
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := hashedTransferTx.Sig.RecoverAddress(hashedTransferTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the From address of the HashedTrasnfer in the transaction", func() {
				So(recoveredAddr, ShouldResemble, hashedTransferTx.HashedTransfer.From)
			})
		})
	})
}

func TestHashedTransferEIP712Signature(t *testing.T) {
	Convey("If a HashedTransfer transaction with EIP712 signature is valid", t, func() {
		hashedTransferTx := &HashedTransferTransaction{
			HashedTransfer: htlc.HashedTransfer{
				From:       Alice.Address,
				To:         Bob.ID,
				Value:      types.NewBigInt(1),
				HashCommit: [32]byte{},
				Expiry:     10,
			},
			Fee:   types.NewBigInt(0),
			Nonce: 1,
			Sig:   &HashedTransferEIP712Signature{SigEIP712("bfdce31b9c432249ea61dd095d62425f62c648ba61331fb9621b56e7e5d91e2b283a26fc13ba2b66ee7ee1f300174a38c075bdf64094f760b37175a94aab6c041b")},
		}
		Convey("Address recovery should succeed", func() {
			recoveredAddr, err := hashedTransferTx.Sig.RecoverAddress(hashedTransferTx)
			So(err, ShouldBeNil)
			Convey("The recovered address should be the From address of the HashedTrasnfer in the transaction", func() {
				So(recoveredAddr, ShouldResemble, hashedTransferTx.HashedTransfer.From)
			})
		})
	})
}

func TestCheckAndDeliverHashedTransfer(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		state.SetBlockTime(9)
		account.NewAccountFromID(state, Alice.ID, Alice.Address)
		account.NewAccountFromID(state, Bob.ID, Bob.Address)
		account.AddBalance(state, Alice.ID, big.NewInt(100))
		commit := make([]byte, 32)
		hashedTransferTx := HashedTransferTx(Alice.Address, Bob.ID, 20, commit, 10, 1, 1, "3238d6e98ebf089c4bc470da5b4be53b84b916f97b8fa20e91d1fa2d4fbba1096fc4580d4aa82cc30243eaebd59577dc435d606ca6460d74cff8e15038a91dce1c")
		Convey("If a HashedTransfer transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusPending)
					So(htlc.GetHashedTransfer(state, txHash), ShouldResemble, &hashedTransferTx.HashedTransfer)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(80)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
					Convey("CheckTx again should return Duplicated", func() {
						r := hashedTransferTx.CheckTx(state)
						So(r.Code, ShouldEqual, response.HashedTransferDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := hashedTransferTx.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.HashedTransferDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
							So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
						})
					})
				})
			})
			Convey("RawTransferTx should return the same encoded raw tx", func() {
				rawTx := RawHashedTransferTx(Alice.Address, Bob.ID, 20, commit, 10, 1, 1, "3238d6e98ebf089c4bc470da5b4be53b84b916f97b8fa20e91d1fa2d4fbba1096fc4580d4aa82cc30243eaebd59577dc435d606ca6460d74cff8e15038a91dce1c")
				So(rawTx, ShouldResemble, EncodeTx(hashedTransferTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *HashedTransferTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, hashedTransferTx)
				})
			})
		})
		Convey("If a HashedTransfer transaction from LikeChainID is valid", func() {
			hashedTransferTx.HashedTransfer.From = Alice.ID
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("3d4ed90d366c060544e69468f0e9bdbcccdb238c9becfc4a42cb21e106f8f1aa3193506e41123f46421ef4c6faa38810a7dac8045725a18053990e70ed46e4581b")}
			Convey("CheckTx should succeed", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusPending)
					So(htlc.GetHashedTransfer(state, txHash), ShouldResemble, &hashedTransferTx.HashedTransfer)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(80)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a HashedTransfer transaction has invalid format", func() {
			hashedTransferTx.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidFormat.Code)
				Convey("DeliverTx should return InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a HashedTransfer transaction has invalid sender", func() {
			hashedTransferTx.HashedTransfer.From = Mallory.Address
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("3096c8cd2037b90644474c1e3159639a01dc3a43a96ef7b5c61c74db987252937e3e296b4d8d0c0b17cdea55abde8305e3ba258d41bd87fbb93cf2046af5d7941c")}
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferSenderNotRegistered.Code)
				Convey("DeliverTx should return SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a HashedTransfer transaction is not signed by the sender", func() {
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("484922561daf238c516bc7144a922064dc11ede4249252d8489f78e00610d9ac599fa95b18790c534ecfd7e3721d4913b4d672632b4c8211804ff44aa589586d1c")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a HashedTransfer transaction has invalid nonce", func() {
			hashedTransferTx.Nonce = 2
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("b8918fa1f95e76697ea44bf80465044b8820a9bc03807e4cc9aa4bf054bd9a1b14f7539309b983109bbcf96980b8381db63d07cc2f9a4b82fe354f9a94f457bc1c")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidNonce.Code)
				Convey("DeliverTx should return InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If the receiver of the HashedTransfer transaction is an address without a LikeChain ID", func() {
			hashedTransferTx.HashedTransfer.To = Carol.Address
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("3c6e0aeb32fcecebd2f69c8289e4fac4cd20a271fba33d26853b07c943cf92a105183b7eaa16e545442d65ae5f9367165403d6c7d2d4864fb3fcdf102c0ef8e21b")}
			Convey("CheckTx should return InvalidReceiver", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidReceiver.Code)
				Convey("DeliverTx should return InvalidReceiver", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferInvalidReceiver.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If the receiver of the HashedTransfer transaction is an address with a LikeChain ID", func() {
			hashedTransferTx.HashedTransfer.To = Bob.Address
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("d82b5f9f90fbb2e7b19149694a67c4ef60f738639a04d8e1c2ddde530ba0cec01b93b6ae6d46b7d18a9378b29af43886ead48a0e0a5326aabed710ac81e702111b")}
			Convey("CheckTx should return Success", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should return Success", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusPending)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(80)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If the sender of the HashedTransfer transaction has not enough balance for the sum of the value and the fee", func() {
			hashedTransferTx.Fee.Int = big.NewInt(81)
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("7ed26ceb839ca9713e2464a34ba34dca1474fec17a7baec0268cc04efd20f7f059a6897618409875ff3d4879586cbcc59314a51074c1658888df3f72a5e4cd5d1c")}
			Convey("CheckTx should return NotEnoughBalance", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferNotEnoughBalance.Code)
				Convey("DeliverTx should return NotEnoughBalance", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferNotEnoughBalance.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If the sender of the HashedTransfer transaction has just enough balance for the sum of the value and the fee", func() {
			hashedTransferTx.Fee.Int = big.NewInt(80)
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("f18ddb68035456db30c0b2498780248a95d07920e0568c88e2761c0c8063629c7b2d449278825acbb39ed69ab518c709b270aaf2ab1ce9b0da37abf5219709c61c")}
			Convey("CheckTx should return Success", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should return Success", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusPending)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(80)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If the expiry time of the HashedTransfer transaction is equal to the block time", func() {
			hashedTransferTx.HashedTransfer.Expiry = 9
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("494234e8ee4c86eafb0105f76c5c852978f527a5633f6c1ea52eb0027bb5b5931224f1bd45d9ab68223bbc0c752d4fbedccecda339b3538e5cbe89d5b39d47561c")}
			Convey("CheckTx should return InvalidExpiry", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidExpiry.Code)
				Convey("DeliverTx should return InvalidExpiry", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferInvalidExpiry.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If the expiry time of the HashedTransfer transaction is less than the block time", func() {
			hashedTransferTx.HashedTransfer.Expiry = 8
			hashedTransferTx.Sig = &HashedTransferJSONSignature{Sig("fa9ea06280ea6eb9395f505d3fdc851b530d27e3347faa5f213ee560a479d6a325048d2267389b23fc396ef047958bc1be323249891913f0d6eaedfed7deee021b")}
			Convey("CheckTx should return InvalidExpiry", func() {
				r := hashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.HashedTransferInvalidExpiry.Code)
				Convey("DeliverTx should return InvalidExpiry", func() {
					txHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
					r := hashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.HashedTransferInvalidExpiry.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 2)
				})
			})
		})
	})
}
