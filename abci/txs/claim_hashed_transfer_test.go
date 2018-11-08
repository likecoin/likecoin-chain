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

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClaimHashedTrasnferValidateFormat(t *testing.T) {
	Convey("For a ClaimHashedTransfer transaction", t, func() {
		htlcTxHash := make([]byte, tmhash.Size)
		secret := make([]byte, 32)
		claimHashedTransferTx := ClaimHashedTransferTx(Alice.Address, htlcTxHash, secret, 1, "936578d55fca7cc1028efa5498cb78d206729bb7199b8bbb8a64f56edcc8d1d32514f04505eea1201c93e1c77877f2ba1f5f972868ea24171b619cec93cf011a1c")
		Convey("If the transaction has valid format", func() {
			Convey("The validation should succeed", func() {
				So(claimHashedTransferTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the signature is nil", func() {
			claimHashedTransferTx.Sig = nil
			Convey("The validation should fail", func() {
				So(claimHashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the HTLC TxHash has invalid length", func() {
			claimHashedTransferTx.HTLCTxHash = make([]byte, tmhash.Size-1)
			Convey("The validation should fail", func() {
				So(claimHashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
		Convey("If the secret is nil", func() {
			claimHashedTransferTx.Secret = nil
			Convey("The validation should succeed", func() {
				So(claimHashedTransferTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the secret has length 32", func() {
			Convey("The validation should succeed", func() {
				So(claimHashedTransferTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the secret has length 0", func() {
			claimHashedTransferTx.Secret = make([]byte, 0, 16)
			Convey("The validation should succeed", func() {
				So(claimHashedTransferTx.ValidateFormat(), ShouldBeTrue)
			})
		})
		Convey("If the secret has other length", func() {
			claimHashedTransferTx.Secret = make([]byte, 1)
			Convey("The validation should fail", func() {
				So(claimHashedTransferTx.ValidateFormat(), ShouldBeFalse)
			})
		})
	})
}

func TestClaimHashedTransferSignature(t *testing.T) {
	Convey("For a ClaimHashedTransfer transaction", t, func() {
		htlcTxHash := make([]byte, tmhash.Size)
		secret := make([]byte, 32)
		claimHashedTransferTx := ClaimHashedTransferTx(Alice.Address, htlcTxHash, secret, 1, "33cbe33180cd8bd7056c13f4f955d804702138680a977e6d979449057fb9e0a507597a542dc1fbe43410c0bc18915b8a1aab66366ac124b1a17cd223e260095f1b")
		Convey("If a ClaimHashedTransfer transaction is valid with secret length 32", func() {
			Convey("Address recovery should succeed", func() {
				recoveredAddr, err := claimHashedTransferTx.Sig.RecoverAddress(claimHashedTransferTx)
				So(err, ShouldBeNil)
				Convey("The recovered address should be the From address in the transaction", func() {
					So(recoveredAddr, ShouldResemble, claimHashedTransferTx.From)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction is valid with secret length 0", func() {
			claimHashedTransferTx.Secret = nil
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("c011b32ed5b9df48b22a476dbc4f243a47520facb2bd3ca0ca7b2fa644257240073e895e87418194ea39a3b236c35b7623c42d04f3bf38553188f1fc1deedb301c")}
			Convey("Address recovery should succeed", func() {
				recoveredAddr, err := claimHashedTransferTx.Sig.RecoverAddress(claimHashedTransferTx)
				So(err, ShouldBeNil)
				Convey("The recovered address should be the From address in the transaction", func() {
					So(recoveredAddr, ShouldResemble, claimHashedTransferTx.From)
				})
			})
		})
	})
}

func TestCheckAndDeliverClaimHashedTransfer(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		state.SetBlockTime(0)
		account.NewAccountFromID(state, Alice.ID, Alice.Address)
		account.NewAccountFromID(state, Bob.ID, Bob.Address)
		account.NewAccountFromID(state, Carol.ID, Carol.Address)
		account.AddBalance(state, Alice.ID, big.NewInt(100))
		secret := make([]byte, 32)
		commit := crypto.Sha256(secret)
		hashedTransferTx := HashedTransferTx(Alice.Address, Bob.ID, 20, commit, 10, 1, 1, "1e81458788aec37fcb4934f7624e941d953535a08fd849828ac697deb1fedb432e941d59472bfa3f9ec6932c5c0d371b14515f2b4d64bb9b827c6a92129afb6e1c")
		htlcTxHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
		r := hashedTransferTx.DeliverTx(state, htlcTxHash)
		So(r.Code, ShouldEqual, response.Success.Code)

		claimHashedTransferTx := ClaimHashedTransferTx(Bob.Address, htlcTxHash, secret, 1, "3d24816f574711e5a5ebcc50be5cc75a6a84b6b1363197e43a2bfd419961dc816f6a3c5fca6c203e554e1102ba76712b767f32998970d17eee4e194e45e49b981c")
		Convey("If a ClaimHashedTransfer transaction from address is valid", func() {
			Convey("CheckTx should succeed", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(txstatus.GetStatus(state, htlcTxHash), ShouldEqual, txstatus.TxStatusSuccess)
					So(htlc.GetHashedTransfer(state, htlcTxHash), ShouldBeNil)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(79)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(20)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 2)
					Convey("CheckTx again should return Duplicated", func() {
						r := claimHashedTransferTx.CheckTx(state)
						So(r.Code, ShouldEqual, response.ClaimHashedTransferDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := claimHashedTransferTx.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.ClaimHashedTransferDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
							So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 2)
						})
					})
				})
			})
			Convey("RawTransferTx should return the same encoded raw tx", func() {
				rawTx := RawClaimHashedTransferTx(Bob.Address, htlcTxHash, secret, 1, "3d24816f574711e5a5ebcc50be5cc75a6a84b6b1363197e43a2bfd419961dc816f6a3c5fca6c203e554e1102ba76712b767f32998970d17eee4e194e45e49b981c")
				So(rawTx, ShouldResemble, EncodeTx(claimHashedTransferTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *ClaimHashedTransferTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, claimHashedTransferTx)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction from LikeChainID is valid", func() {
			claimHashedTransferTx.From = Bob.ID
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("5669a83fd580fdb832219883a19bf1ddb5b5974bbdee0f06510b431232a832855600cf000ed8589cfde82357ecd28888de422bd4aba79ed4c8395924a4d32a8b1c")}
			Convey("CheckTx should succeed", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(txstatus.GetStatus(state, htlcTxHash), ShouldEqual, txstatus.TxStatusSuccess)
					So(htlc.GetHashedTransfer(state, htlcTxHash), ShouldBeNil)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(79)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(20)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a HashedTransfer transaction has invalid format", func() {
			claimHashedTransferTx.Sig = nil
			Convey("CheckTx should return InvalidFormat", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidFormat.Code)
				Convey("DeliverTx should return InvalidFormat", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidFormat.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction has unregistered sender", func() {
			claimHashedTransferTx.From = Mallory.Address
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("8962956bce009d54d6b8f21560b4e4411f63e2f654463d9e99ae73a63a07b9a030172c69a228145647cf5f3a87ada932f1ca5d5f77b3c0ab1d07d0adefad7b6f1c")}
			Convey("CheckTx should return SenderNotRegistered", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferSenderNotRegistered.Code)
				Convey("DeliverTx should return SenderNotRegistered", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferSenderNotRegistered.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction is not signed by the sender", func() {
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("090db602613657e02878d2f3b8cd6d70bf1ae1cf7f1f42db1a97533b20cb4c9927f3d9be7f29fce3a2c6750f524b95fa405b5be8630848cbb635299aa3a3fa281b")}
			Convey("CheckTx should return InvalidSignature", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidSignature.Code)
				Convey("DeliverTx should return InvalidSignature", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidSignature.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction has invalid nonce", func() {
			claimHashedTransferTx.Nonce = 2
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("ce0cb306833263e182d89edfb937ebdde43c5a87d81ed9461a6dd224fa4c8da74a5897da9b407c9fd31987021c8f6047bd763c11efd9eadef9ca8f409d10f79d1c")}
			Convey("CheckTx should return InvalidNonce", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidNonce.Code)
				Convey("DeliverTx should return InvalidNonce", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidNonce.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 1)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction has invalid HTLC TxHash", func() {
			claimHashedTransferTx.HTLCTxHash = make([]byte, tmhash.Size)
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("9cb82f2ca5dff26beda3e0a97a9ab753610bb6cf03b665ce57734f7ed35ad803629220e41178689a4807c7263178ad0619230e7b8ac0e5d39098bd084fef14671c")}
			Convey("CheckTx should return TxNotExist", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferTxNotExist.Code)
				Convey("DeliverTx should return TxNotExist", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferTxNotExist.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction is neither from the sender or receiver of the HashedTransfer transaction", func() {
			claimHashedTransferTx.From = Carol.ID
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("a19c9672abcc4b08b464065230ca08114359e418ed30aa9ee949ff6b83dca05f41799cc71917dd336daab8c5fd0b611d7def4ef61cf18eba87137a02383a35351c")}
			Convey("CheckTx should return InvalidSender", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidSender.Code)
				Convey("DeliverTx should return InvalidSender", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidSender.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Carol.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction has invalid secret", func() {
			claimHashedTransferTx.Secret[0]++
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("5c8663bd2d2a0040cb3fed3c2c3fcfd1437845d1e51df15056a80d4cee35242634dd927d3af87d0f2719e62390d6cd112260e1f2dece439ca0304713a553fe201c")}
			Convey("CheckTx should return InvalidSecret", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidSecret.Code)
				Convey("DeliverTx should return InvalidSecret", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferInvalidSecret.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction is executed exactly at expiry time", func() {
			state.SetBlockTime(10)
			Convey("CheckTx should return Expired", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferExpired.Code)
				Convey("DeliverTx should return Expired", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferExpired.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 2)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction is executed after expiry time", func() {
			state.SetBlockTime(11)
			Convey("CheckTx should return Expired", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferExpired.Code)
				Convey("DeliverTx should return Expired", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferExpired.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Bob.ID), ShouldEqual, 2)
				})
			})
		})
	})
}

func TestCheckAndDeliverRevokeHashedTransfer(t *testing.T) {
	Convey("In the beginning", t, func() {
		appCtx := context.NewMock()
		state := appCtx.GetMutableState()
		state.SetBlockTime(0)
		account.NewAccountFromID(state, Alice.ID, Alice.Address)
		account.NewAccountFromID(state, Bob.ID, Bob.Address)
		account.NewAccountFromID(state, Carol.ID, Carol.Address)
		account.AddBalance(state, Alice.ID, big.NewInt(100))
		secret := make([]byte, 32)
		commit := crypto.Sha256(secret)
		hashedTransferTx := HashedTransferTx(Alice.Address, Bob.ID, 20, commit, 10, 1, 1, "1e81458788aec37fcb4934f7624e941d953535a08fd849828ac697deb1fedb432e941d59472bfa3f9ec6932c5c0d371b14515f2b4d64bb9b827c6a92129afb6e1c")
		htlcTxHash := utils.HashRawTx(EncodeTx(hashedTransferTx))
		r := hashedTransferTx.DeliverTx(state, htlcTxHash)
		So(r.Code, ShouldEqual, response.Success.Code)

		claimHashedTransferTx := ClaimHashedTransferTx(Alice.Address, htlcTxHash, nil, 2, "a06d1e079ff87508a0a9186489a0f6ecba5af57ee1186b171e56b1f22766df864b52d0e5d08cddaef3aa322a86169e0772e16e486a8991912f3d7a351e4798fe1c")
		Convey("If a ClaimHashedTransfer transaction (revoke) from address is valid", func() {
			state.SetBlockTime(11)
			Convey("CheckTx should succeed", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(txstatus.GetStatus(state, htlcTxHash), ShouldEqual, txstatus.TxStatusSuccess)
					So(htlc.GetHashedTransfer(state, htlcTxHash), ShouldBeNil)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 3)
					Convey("CheckTx again should return Duplicated", func() {
						r := claimHashedTransferTx.CheckTx(state)
						So(r.Code, ShouldEqual, response.ClaimHashedTransferDuplicated.Code)
						Convey("DeliverTx should return Duplicated", func() {
							r := claimHashedTransferTx.DeliverTx(state, txHash)
							So(r.Code, ShouldEqual, response.ClaimHashedTransferDuplicated.Code)
							So(r.Status, ShouldEqual, txstatus.TxStatusFail)
							So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 3)
						})
					})
				})
			})
			Convey("RawTransferTx should return the same encoded raw tx", func() {
				rawTx := RawClaimHashedTransferTx(Alice.Address, htlcTxHash, nil, 2, "a06d1e079ff87508a0a9186489a0f6ecba5af57ee1186b171e56b1f22766df864b52d0e5d08cddaef3aa322a86169e0772e16e486a8991912f3d7a351e4798fe1c")
				So(rawTx, ShouldResemble, EncodeTx(claimHashedTransferTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *ClaimHashedTransferTransaction
					err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, claimHashedTransferTx)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction from LikeChainID is valid", func() {
			claimHashedTransferTx.From = Alice.ID
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("47ab91ffbca34d0d247bdcb6c5cb5c4670ad540d8951d9e66d5e6cabb213768c50bae76e5778639ae6aad1848280093d0382451b5ad90818a00a5dac98b1cc901b")}
			state.SetBlockTime(11)
			Convey("CheckTx should succeed", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(txstatus.GetStatus(state, htlcTxHash), ShouldEqual, txstatus.TxStatusSuccess)
					So(htlc.GetHashedTransfer(state, htlcTxHash), ShouldBeNil)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 3)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction (revoke) is executed exactly at expiry time", func() {
			state.SetBlockTime(10)
			Convey("CheckTx should succeed", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.Success.Code)
				Convey("DeliverTx should succeed", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.Success.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusSuccess)
					So(txstatus.GetStatus(state, htlcTxHash), ShouldEqual, txstatus.TxStatusSuccess)
					So(htlc.GetHashedTransfer(state, htlcTxHash), ShouldBeNil)
					So(account.FetchBalance(state, Alice.ID).Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(account.FetchBalance(state, Bob.ID).Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 3)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction (revoke) is executed before expiry time", func() {
			state.SetBlockTime(9)
			Convey("CheckTx should return NotYetExpired", func() {
				r := claimHashedTransferTx.CheckTx(state)
				So(r.Code, ShouldEqual, response.ClaimHashedTransferNotYetExpired.Code)
				Convey("DeliverTx should return NotYetExpired", func() {
					txHash := utils.HashRawTx(EncodeTx(claimHashedTransferTx))
					r := claimHashedTransferTx.DeliverTx(state, txHash)
					So(r.Code, ShouldEqual, response.ClaimHashedTransferNotYetExpired.Code)
					So(r.Status, ShouldEqual, txstatus.TxStatusFail)
					So(account.FetchNextNonce(state, Alice.ID), ShouldEqual, 3)
				})
			})
		})
	})
}
