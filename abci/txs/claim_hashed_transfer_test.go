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
		claimHashedTransferTx := ClaimHashedTransferTx(Alice.Address, htlcTxHash, secret, 1, "936578d55fca7cc1028efa5498cb78d206729bb7199b8bbb8a64f56edcc8d1d32514f04505eea1201c93e1c77877f2ba1f5f972868ea24171b619cec93cf011a1c")
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
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("18acf4c455ab70f45e2b136d31f2449b183f1469d1db1cba8747d29f287cb7da23519efcd4b454eb55969cb16c969a8c58072847fd27ed9259e299ac0a89a2f21b")}
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

		claimHashedTransferTx := ClaimHashedTransferTx(Bob.Address, htlcTxHash, secret, 1, "631b6a96c706efcf86dcf4dae0cbc31d0ff0c663c6a805e4848ccad84b9155d84a0c76064804a21c642d47a5fdf10850751734ba6f60e71f2203a99c79eb42cb1b")
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
				rawTx := RawClaimHashedTransferTx(Bob.Address, htlcTxHash, secret, 1, "631b6a96c706efcf86dcf4dae0cbc31d0ff0c663c6a805e4848ccad84b9155d84a0c76064804a21c642d47a5fdf10850751734ba6f60e71f2203a99c79eb42cb1b")
				So(rawTx, ShouldResemble, EncodeTx(claimHashedTransferTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *ClaimHashedTransferTransaction
					err := types.AminoCodec().UnmarshalBinary(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, claimHashedTransferTx)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction from LikeChainID is valid", func() {
			claimHashedTransferTx.From = Bob.ID
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("ae2017223110189a3d95c32d42061ee7c27a119df73eb676ab2a1d5a8d35bee42d6ca8e39e29235d7b27d9cfda769eb1c1a3722c96581142aa520122e6ecb1901b")}
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
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("4d966e88ec4c22c7d5ef9d6e1b1df45910b95b6a2010cb9658ec358d98fb8b6d08aa061244dcfa4f70466e73e86466045a210bd984dc954921186be543259a7c1b")}
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
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("d2bd3fd72c4a1c3ff228f9b8d8b1d237be969401cc3c8f2eb43b880c6157d39300c6703e524c89219dfd8a1cfc82b92d784d308ce047fa7514560a5f0e2b5f2c1b")}
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
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("f67b05108838dfe30379510122c99d3e9d3ca9f027595568e2749d1bf90d10ff65e857f62c5ff0054dfdd14952c64abe006bd88210b0918d0db754be3c79245e1c")}
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
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("6e55f5c035bc7f4e7e39399278efcdb495501052cb39c2069e62195d6aeeffee2a848ab16342945365c05d98d62543e96c6af0297208ab5a12bde19a4b219a411b")}
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
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("9ec382f15b3deeff3ae4415b64984662698f0f495ec4968fb31b00a6f124650b4a210190cc9ee2450b28c573648576cb686d61a057a6d040db79c5ecc9a8df9f1c")}
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
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("0f2f46443a6212f824fd76bf5cfcfd73e7154a5b6fabd71ed14ae920885f01d82bb3c6f209930c0a8bb102926272a8e345f3b9d0c09d347e83f189755ad2d12a1c")}
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

		claimHashedTransferTx := ClaimHashedTransferTx(Alice.Address, htlcTxHash, nil, 2, "265ef8ce2b992cba58a03546fd7d22520de77a4c73bebe8a556c5ab14c594e5920d7d28969c45a8ebab8766f9dc9ea5db33b65a1a435c022745a3c9b82ddf0b81c")
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
				rawTx := RawClaimHashedTransferTx(Alice.Address, htlcTxHash, nil, 2, "265ef8ce2b992cba58a03546fd7d22520de77a4c73bebe8a556c5ab14c594e5920d7d28969c45a8ebab8766f9dc9ea5db33b65a1a435c022745a3c9b82ddf0b81c")
				So(rawTx, ShouldResemble, EncodeTx(claimHashedTransferTx))
				Convey("Decoded transaction should resemble the original transaction", func() {
					var decodedTx *ClaimHashedTransferTransaction
					err := types.AminoCodec().UnmarshalBinary(rawTx, &decodedTx)
					So(err, ShouldBeNil)
					So(decodedTx, ShouldResemble, claimHashedTransferTx)
				})
			})
		})
		Convey("If a ClaimHashedTransfer transaction from LikeChainID is valid", func() {
			claimHashedTransferTx.From = Alice.ID
			claimHashedTransferTx.Sig = &ClaimHashedTransferJSONSignature{Sig("2a7d06dff5db5aba25bd08d1ba447220fcfaf12e989fe8c42534345e7d26011937f4c0b3d848a07a7866bcc098af55258f00cb438e560edaa1154ee40129b5191c")}
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
