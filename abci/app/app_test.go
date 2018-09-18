package app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/likecoin/likechain/abci/account"
	"github.com/likecoin/likechain/abci/context"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/iavl"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	. "github.com/smartystreets/goconvey/convey"
)

func encodeTransaction(tx *types.Transaction) []byte {
	data, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}
	return data
}

func identifierFromID(id []byte) *types.Identifier {
	return &types.Identifier{
		Id: &types.Identifier_LikeChainID{
			LikeChainID: types.NewLikeChainID(id),
		},
	}
}

func identifierFromEthAddr(ethAddr *common.Address) *types.Identifier {
	return &types.Identifier{
		Id: &types.Identifier_Addr{
			Addr: &types.Address{
				Content: ethAddr[:],
			},
		},
	}
}

func signature(sig []byte) *types.Signature {
	return &types.Signature{
		Version: 1,
		Content: sig,
	}
}

func bigInteger(n *big.Int) *types.BigInteger {
	return &types.BigInteger{
		Content: n.Bytes(),
	}
}

func registerTx(ethAddr *common.Address, sig []byte) *types.Transaction {
	return &types.Transaction{
		Tx: &types.Transaction_RegisterTx{
			RegisterTx: &types.RegisterTransaction{
				Addr: &types.Address{
					Content: ethAddr[:],
				},
				Sig: &types.Signature{
					Version: 1,
					Content: sig,
				},
			},
		},
	}
}

type transferTarget struct {
	ToID   []byte
	ToAddr string
	Value  *big.Int
	Remark []byte
}

func transferTx(fromID []byte, fromEthAddr *common.Address, toList []transferTarget, fee *big.Int, nonce uint64, sig []byte) *types.Transaction {
	targetList := []*types.TransferTransaction_TransferTarget{}
	for _, to := range toList {
		target := types.TransferTransaction_TransferTarget{
			Value:  bigInteger(to.Value),
			Remark: to.Remark,
		}
		if to.ToID != nil {
			target.To = identifierFromID(to.ToID)
		} else {
			ethAddr := common.HexToAddress(to.ToAddr)
			target.To = identifierFromEthAddr(&ethAddr)
		}
		targetList = append(targetList, &target)
	}
	tx := types.TransferTransaction{
		ToList: targetList,
		Nonce:  nonce,
		Fee:    bigInteger(fee),
		Sig:    signature(sig),
	}
	if fromID != nil {
		tx.From = identifierFromID(fromID)
	} else {
		tx.From = identifierFromEthAddr(fromEthAddr)
	}
	return &types.Transaction{
		Tx: &types.Transaction_TransferTx{
			TransferTx: &tx,
		},
	}
}

func withdrawTx(fromID []byte, fromEthAddr, toEthAddr *common.Address, value, fee *big.Int, nonce uint64, sig []byte) *types.Transaction {
	tx := types.WithdrawTransaction{
		ToAddr: &types.Address{
			Content: toEthAddr[:],
		},
		Value: bigInteger(value),
		Fee:   bigInteger(fee),
		Nonce: nonce,
		Sig: &types.Signature{
			Version: 1,
			Content: sig,
		},
	}
	if fromID != nil {
		tx.From = identifierFromID(fromID)
	} else {
		tx.From = identifierFromEthAddr(fromEthAddr)
	}
	return &types.Transaction{
		Tx: &types.Transaction_WithdrawTx{
			WithdrawTx: &tx,
		},
	}
}

type accountInfoRes struct {
	Balance   *big.Int
	NextNonce uint64
}

func getAccountInfoRes(data []byte) *accountInfoRes {
	accountInfo := struct {
		Balance   string `json:"balance"`
		NextNonce uint64 `json:"nextNonce"`
	}{}
	err := json.Unmarshal(data, &accountInfo)
	if err != nil {
		return nil
	}
	balance, succ := new(big.Int).SetString(accountInfo.Balance, 10)
	if !succ {
		return nil
	}
	return &accountInfoRes{
		Balance:   balance,
		NextNonce: accountInfo.NextNonce,
	}
}

func getTxStateRes(data []byte) string {
	return string(data)
}

func TestRegistration(t *testing.T) {
	Convey("At initial state", t, func() {
		mockCtx := context.NewMock()
		app := &LikeChainApplication{
			ctx: mockCtx.ApplicationContext,
		}
		app.InitChain(abci.RequestInitChain{})
		Convey("Given a valid RegisterTransaction", func() {
			ethAddr := common.HexToAddress("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9")
			sig := common.Hex2Bytes("65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141b")
			tx := registerTx(&ethAddr, sig)
			rawTx := encodeTransaction(tx)
			Convey("The registration should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				Convey("Duplicated registration in the same block should fail during deliverTx", func() {
					checkTxResDup := app.CheckTx(rawTx)
					So(checkTxResDup.Code, ShouldEqual, response.RegisterCheckTxDuplicated.Code)
					deliverTxResDup := app.DeliverTx(rawTx)
					So(deliverTxResDup.Code, ShouldEqual, response.RegisterDeliverTxDuplicated.Code)
				})
				app.EndBlock(abci.RequestEndBlock{
					Height: 1,
				})
				app.Commit()
				likeChainID := deliverTxRes.Data
				Convey("Query account_info using address should return the corresponding info", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9"),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Query account_info using returned LikeChain ID should return the corresponding info", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainID)
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("But repeated registration should fail", func() {
					checkTxRes = app.CheckTx(rawTx)
					So(checkTxRes.Code, ShouldEqual, response.RegisterCheckTxDuplicated.Code)
					deliverTxRes = app.DeliverTx(rawTx)
					So(deliverTxRes.Code, ShouldEqual, response.RegisterDeliverTxDuplicated.Code)
				})
			})
		})

		Convey("Given a RegisterTransaction with other's signature", func() {
			ethAddr := common.HexToAddress("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9")
			sig := common.Hex2Bytes("b287bb3c420155326e0a7fe3a66fed6c397a4bdb5ddcd54960daa0f06c1fbf06300e862dbd3ae3daeae645630e66962b81cf6aa9ffb258aafde496e0310ab8551c")
			tx := registerTx(&ethAddr, sig)
			rawTx := encodeTransaction(tx)
			Convey("The registration should fail", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.RegisterCheckTxInvalidSignature.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.RegisterDeliverTxInvalidSignature.Code)
			})
		})

		Convey("Given a RegisterTransaction with invalid signature", func() {
			ethAddr := common.HexToAddress("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9")
			sig := common.Hex2Bytes("65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f541400")
			tx := registerTx(&ethAddr, sig)
			rawTx := encodeTransaction(tx)
			Convey("The registration should fail", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.RegisterCheckTxInvalidSignature.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.RegisterDeliverTxInvalidSignature.Code)
			})
		})
	})
}

func TestTransfer(t *testing.T) {
	Convey("Given accounts A, B, C with 100, 200, 300 balance each", t, func() {
		mockCtx := context.NewMock()
		app := &LikeChainApplication{
			ctx: mockCtx.ApplicationContext,
		}
		app.InitChain(abci.RequestInitChain{})

		likeChainIDs := [][]byte{}
		regInfos := []struct {
			Addr string
			Sig  string
		}{
			{"0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9", "65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141b"},
			{"0x833a907efe57af3040039c90f4a59946a0bb3d47", "fb141eca7550c8f6d1f37b20536b06327ba29537a6178ea39e9d7747abdc8c2c4daa4ab23accf2157a2eb5ec1eb54ee68159c5b39f7f4ac17087fd71afd374121b"},
			{"0xaa2f5b6ae13ba7a3d466ffce8cd390519337aade", "e906aaf924d636c9b03160d358ec9a20b2b79770e807e84f4cf7f274149ff0b1185b8508adf8cbbc0436b3215cb6e77fea84e97340cbdacd2bcc0bac4a374b441b"},
		}

		for n, regInfo := range regInfos {
			ethAddr := common.HexToAddress(regInfo.Addr)
			sig := common.Hex2Bytes(regInfo.Sig)
			tx := registerTx(&ethAddr, sig)
			rawTx := encodeTransaction(tx)
			deliverTxRes := app.DeliverTx(rawTx)
			So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
			likeChainID := deliverTxRes.Data
			likeChainIDs = append(likeChainIDs, likeChainID)
			account.SaveBalance(mockCtx.GetMutableState(), identifierFromID(likeChainID), big.NewInt(int64(n+1)*100))
		}
		app.EndBlock(abci.RequestEndBlock{
			Height: 1,
		})
		app.Commit()

		for i, likeChainIDBase64 := range []string{"bDH8FUIuutKKr5CJwwZwL2dUC1M=", "hZ8Rt1VppOsElsUTj9QsxSrujPU=", "1MaeSeg6YEf0bkKy0FOh8MbnDqQ="} {
			likeChainID, _ := base64.StdEncoding.DecodeString(likeChainIDBase64)
			So(bytes.Compare(likeChainIDs[i], likeChainID), ShouldBeZeroValue)
		}

		for n, likeChainID := range likeChainIDs {
			likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainID)
			queryRes := app.Query(abci.RequestQuery{
				Path: "account_info",
				Data: []byte(likeChainIDBase64),
			})
			So(queryRes.Code, ShouldEqual, response.Success.Code)
			accountInfo := getAccountInfoRes(queryRes.Value)
			So(accountInfo, ShouldNotBeNil)
			So(accountInfo.Balance.Cmp(big.NewInt(int64(n+1)*100)), ShouldBeZeroValue)
			So(accountInfo.NextNonce, ShouldEqual, 1)
		}

		Convey("Given a TransferTransaction from A to B value 1", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[1],
					Value: big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("8bdffbad4cc86028e0212477930444f5f3e329ac8f9f23f866bfc70fa5c157ea70d50e073b29662fb11216b1a8d82157b2e1c48185c910c5ac07fb2b238de4651c")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info by LikeChain ID should return the correct balances and nextNonce", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainIDs[0])
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					likeChainIDBase64 = base64.StdEncoding.EncodeToString(likeChainIDs[1])
					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(201)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query account_info by Ethereum address should return the correct balances and nextNonce", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(regInfos[0].Addr),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(regInfos[1].Addr),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(201)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query tx_state should return success", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "success")
				})
				Convey("But repeated transfer with the same transaction should fail", func() {
					checkTxRes := app.CheckTx(rawTx)
					So(checkTxRes.Code, ShouldEqual, response.TransferCheckTxDuplicated.Code)
					deliverTxRes := app.DeliverTx(rawTx)
					So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxDuplicated.Code)
				})
			})
		})

		Convey("Given a TransferTransaction from A's Ethereum address to B's Ethereum address with value 1", func() {
			from := common.HexToAddress(regInfos[0].Addr)
			toList := []transferTarget{
				{
					ToAddr: regInfos[1].Addr,
					Value:  big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("3d7e17b7b95f1462b5ac0b238aed583934619b9da20dcaf71485a66ab3ff086646c8eaca3bf39c1d51d78cffaeb4d2f6f678147aa202bf9c398c42a2d46256f11c")
			tx := transferTx(nil, &from, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info by LikeChain ID address should return the correct balance", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainIDs[0])
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					likeChainIDBase64 = base64.StdEncoding.EncodeToString(likeChainIDs[1])
					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(201)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query account_info by Ethereum address should return the correct balance", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(regInfos[0].Addr),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(regInfos[1].Addr),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(201)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query tx_state should return success", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "success")
				})
			})
		})

		Convey("Given a TransferTransaction from A's LikeChain ID to B's Ethereum address (with value 1) and C's LikeChain ID (with value 2)", func() {
			toList := []transferTarget{
				{
					ToAddr: regInfos[1].Addr,
					Value:  big.NewInt(1),
				},
				{
					ToID:  likeChainIDs[2],
					Value: big.NewInt(2),
				},
			}
			sig := common.Hex2Bytes("695d4935a112f3d3715c873f4205e84b3eb56ad84f155fb21e834a8eb3e9a8d822941d0a8660d4abab6d39653894fa9494095728741fb9b6d9a72594b028853d1b")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info by LikeChain ID address should return the correct balance", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainIDs[0])
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(97)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					likeChainIDBase64 = base64.StdEncoding.EncodeToString(likeChainIDs[1])
					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(201)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)

					likeChainIDBase64 = base64.StdEncoding.EncodeToString(likeChainIDs[2])
					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(302)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query account_info by Ethereum address should return the correct balance", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(regInfos[0].Addr),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(97)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(regInfos[1].Addr),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(201)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)

					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(regInfos[2].Addr),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(302)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query tx_state should return success", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "success")
				})
			})
		})

		Convey("Given a TransferTransaction from unregistered Ethereum address", func() {
			from := common.HexToAddress("0xf6c45b1c4b73ccaeb1d9a37024a6b9fa711d7e7e")
			toList := []transferTarget{
				{
					ToAddr: regInfos[1].Addr,
					Value:  big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("94ae22209628b549d6c84eb345cd448c412e2ecab134ba9dba4457df8e0e0f52460aedbc8214fc6ca79f05eed2ead6d9149824ecf72e238eebc83fb4156989481b")
			tx := transferTx(nil, &from, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should fail with SenderNotRegistered", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.TransferCheckTxSenderNotRegistered.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxSenderNotRegistered.Code)
				Convey("Then query tx_state should return fail", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "fail")
				})
			})
		})

		Convey("Given a TransferTransaction to unregistered LikeChain ID receiver(s)", func() {
			unregLikeChainIDBase64 := "j/FYH9yZaCgTbAuhvdvk+op9Vas="
			unregLikeChainID, _ := base64.StdEncoding.DecodeString(unregLikeChainIDBase64)
			toList := []transferTarget{
				{
					ToAddr: regInfos[1].Addr,
					Value:  big.NewInt(1),
				},
				{
					ToID:  unregLikeChainID,
					Value: big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("a160704e2f191ab623a26edd8e4bc3a3a843270e57b9463dca6697c74a54e332058a8e24b8ffe3765ff91df33b8d8fd7e9f3c690c29335f603a7603fe162d2bf1c")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should fail with InvalidReceiver", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.TransferCheckTxInvalidReceiver.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxInvalidReceiver.Code)
				Convey("Then query tx_state should return fail", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "fail")
				})
			})
		})

		Convey("Given a TransferTransaction to unregistered Ethereum address", func() {
			toList := []transferTarget{
				{
					ToAddr: "0xf6c45b1c4b73ccaeb1d9a37024a6b9fa711d7e7e",
					Value:  big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("869409224bb3e0ca7aac8e2246716895e16b2bef4be5fcd8673ae399a61624d331ae0e3a2b407c0fca3f4627f6bca2a64322408f94eb811607964a8bfc2f37991c")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info for sender's account should return the correct balance", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainIDs[0])
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)
				})
				Convey("Query account_info for receiver's address should return nothing before registration", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte("0xf6c45b1c4b73ccaeb1d9a37024a6b9fa711d7e7e"),
					})
					So(queryRes.Code, ShouldEqual, response.QueryInvalidIdentifier.Code)
				})
				Convey("Then query tx_state should return success", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "success")
				})
				Convey("Registration for the receiver's address should success", func() {
					ethAddr := common.HexToAddress("0xf6c45b1c4b73ccaeb1d9a37024a6b9fa711d7e7e")
					sig := common.Hex2Bytes("5221a47f0c1042f67951e28c513634190a7c4d77703a642d495ac5ef6397c4ec4d6ab2f7d1cda7c05f8e61d781aa2a4fa6e98c4382f741c4a7ab8e4de1d3fee31c")
					tx := registerTx(&ethAddr, sig)
					rawTx := encodeTransaction(tx)
					checkTxRes := app.CheckTx(rawTx)
					So(checkTxRes.Code, ShouldEqual, response.Success.Code)
					deliverTxRes := app.DeliverTx(rawTx)
					So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
					app.EndBlock(abci.RequestEndBlock{
						Height: 3,
					})
					app.Commit()
					Convey("The receiving Ethereum address should have balance after registration", func() {
						likeChainID := deliverTxRes.Data
						likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainID)
						queryRes := app.Query(abci.RequestQuery{
							Path: "account_info",
							Data: []byte(likeChainIDBase64),
						})
						So(queryRes.Code, ShouldEqual, response.Success.Code)
						accountInfo := getAccountInfoRes(queryRes.Value)
						So(accountInfo, ShouldNotBeNil)
						So(accountInfo.Balance.Cmp(big.NewInt(1)), ShouldBeZeroValue)
						So(accountInfo.NextNonce, ShouldEqual, 1)

						queryRes = app.Query(abci.RequestQuery{
							Path: "account_info",
							Data: []byte("0xf6c45b1c4b73ccaeb1d9a37024a6b9fa711d7e7e"),
						})
						So(queryRes.Code, ShouldEqual, response.Success.Code)
						accountInfo = getAccountInfoRes(queryRes.Value)
						So(accountInfo, ShouldNotBeNil)
						So(accountInfo.Balance.Cmp(big.NewInt(1)), ShouldBeZeroValue)
						So(accountInfo.NextNonce, ShouldEqual, 1)
					})
				})
			})
		})

		Convey("Given a TransferTransaction with normal remark", func() {
			toList := []transferTarget{
				{
					ToID:   likeChainIDs[1],
					Value:  big.NewInt(1),
					Remark: []byte("99BottlesOfBeer"),
				},
			}
			sig := common.Hex2Bytes("70f5547ecfd68a66cdd5326da7887146ec83af894f4942361ff30c9d15f742247ed947235e05e637ecb0a08dcd606e1fb0918843a1de7ce6039478416f1b3e361b")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
			})
		})

		Convey("Given a TransferTransaction with 4096 bytes remark", func() {
			zeros := make([]byte, 4096)
			toList := []transferTarget{
				{
					ToID:   likeChainIDs[1],
					Value:  big.NewInt(1),
					Remark: zeros,
				},
			}
			sig := common.Hex2Bytes("d591986da187995d1a709327cb7accc36ec6f1b9ab7fe0aa7238ffcbcfdc8c6d10478b24d3e9d809c55fc74b3a026ecee785e12538acdab3ca21fbc89dc166601c")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
			})
		})

		Convey("Given a TransferTransaction with 4097 remark", func() {
			zeros := make([]byte, 4097)
			toList := []transferTarget{
				{
					ToID:   likeChainIDs[1],
					Value:  big.NewInt(1),
					Remark: zeros,
				},
			}
			sig := common.Hex2Bytes("413ae14d2108d726eda8e5d35eaf366b155835cb23476b1af38436c52796cc4a3a23e20bd3bde88d48e1a1d377582a311d8917da2fb62ed65726a9b2ea8ba4fc1b")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should fail", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.TransferCheckTxInvalidFormat.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxInvalidFormat.Code)
			})
		})

		Convey("Given a TransferTransaction from A to B value 0", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[1],
					Value: big.NewInt(0),
				},
			}
			sig := common.Hex2Bytes("3b345b0fe343d757ee2d0f554ccefbf1d359105522ead2b89e681c43dee79f4518b46c82607f44eb3a2c47ad7e38ff318afd585bb3556094583cd2d7c9cdb6e11b")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info by LikeChain ID should return the correct balances and nextNonce", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainIDs[0])
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(100)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					likeChainIDBase64 = base64.StdEncoding.EncodeToString(likeChainIDs[1])
					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(200)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query tx_state should return success", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "success")
				})
			})
		})

		Convey("Given a TransferTransaction from A to C value 100", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[2],
					Value: big.NewInt(100),
				},
			}
			sig := common.Hex2Bytes("80ee2c9c9a0dcc0d9c131e621961841e3552df477b9242120af28dd039a218910c3d0e17d7f4622f1ecded69005964ff5476016295b44615fdf35ce3031aa8621c")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info by LikeChain ID should return the correct balances and nextNonce", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainIDs[0])
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)

					likeChainIDBase64 = base64.StdEncoding.EncodeToString(likeChainIDs[2])
					queryRes = app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo = getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(400)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 1)
				})
				Convey("Then query tx_state should return success", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "success")
				})
			})
		})

		Convey("Given a TransferTransaction with value sum more than 100", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[1],
					Value: big.NewInt(50),
				},
				{
					ToID:  likeChainIDs[2],
					Value: big.NewInt(51),
				},
			}
			sig := common.Hex2Bytes("e7281027a04b63380e00d2cfd7812543bfccf805f49f37db993786fee675067541f043d47824d0484a5a97998f0f9b1568ce231370a6d069297d8472479093151c")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should fail", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.TransferCheckTxNotEnoughBalance.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxNotEnoughBalance.Code)
			})
		})

		Convey("Given 2 TransferTransactions with value sum more than 100", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[1],
					Value: big.NewInt(50),
				},
			}
			sig := common.Hex2Bytes("5f9ca2aad30ede4d8edde79e438f02706a02bb249c91ebd66f6b8e886797ff8746ed31366582f196b148539f1f35f65d5e57a78d77a8833f85b5dc5f22b65b3b1b")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx1 := encodeTransaction(tx)

			toList = []transferTarget{
				{
					ToID:  likeChainIDs[2],
					Value: big.NewInt(51),
				},
			}
			sig = common.Hex2Bytes("dc53aef2b373c946ea246f94d1c58704b57a0d2f58215eaa5b1ed339407d4e504ab54d3e220882af5a953abee1f3578daa44e4a004e6c0a5e3d3b4935d9f9d0b1b")
			tx = transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 2, sig)
			rawTx2 := encodeTransaction(tx)
			Convey("The first TransferTransactions should success", func() {
				checkTxRes := app.CheckTx(rawTx1)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx1)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				Convey("The second TransferTransactions should fail", func() {
					deliverTxRes := app.DeliverTx(rawTx2)
					So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxNotEnoughBalance.Code)
				})
			})
		})

		Convey("Given a TransferTransaction with invalid nonce", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[1],
					Value: big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("01b5993540d724f67f377f4667d8911b121c4f9ade9f8d87a899d878cec71d336c755ae6f2f066c250a13adfbcd8a3e02e970529e7aceb43fd8ef618ee3b7b9d1b")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 2, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should fail", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.TransferCheckTxInvalidNonce.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxInvalidNonce.Code)
			})
		})

		Convey("Given a TransferTransaction with invalid signature", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[1],
					Value: big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("80a1fd124c4b3f1673ff76295e2660280d48711fb2c81aae78d0a9b2fc521e310f9f2a7e59c266852b9a862e880e2bae91359a86372a307041f9342b9c7715c21b")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should fail", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.TransferCheckTxInvalidSignature.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.TransferDeliverTxInvalidSignature.Code)
			})
		})

		Convey("Given a TransferTransaction from A to A value 1", func() {
			toList := []transferTarget{
				{
					ToID:  likeChainIDs[0],
					Value: big.NewInt(1),
				},
			}
			sig := common.Hex2Bytes("f5282e361d732ba6e175e0fd73cc1c72059df1b9d2b0a6d0259cce37a555063d53b12f5eaa9f5cac7ce6e22bf7f6c76e9aee0b4042b4fe144807741eb986f0271c")
			tx := transferTx(likeChainIDs[0], nil, toList, big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The transfer should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info by LikeChain ID should return the correct balances and nextNonce", func() {
					likeChainIDBase64 := base64.StdEncoding.EncodeToString(likeChainIDs[0])
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(100)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)
				})
				Convey("Then query tx_state should return success", func() {
					txHash := tmhash.Sum(rawTx)
					queryRes := app.Query(abci.RequestQuery{
						Path: "tx_state",
						Data: txHash,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					So(getTxStateRes(queryRes.Value), ShouldEqual, "success")
				})
			})
		})
	})
}

func TestWithdraw(t *testing.T) {
	Convey("Given account A with 100 balance", t, func() {
		mockCtx := context.NewMock()
		app := &LikeChainApplication{
			ctx: mockCtx.ApplicationContext,
		}
		app.InitChain(abci.RequestInitChain{})

		regInfo := struct {
			Addr string
			Sig  string
		}{
			"0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9",
			"65e6d31224fbcec8e41251d7b014e569d4a94c866227637c6b1fcf75a4505f241b2009557e79d5879a8bfbbb5dec86205c3481ed3042ad87f0643778022f54141b",
		}

		ethAddr := common.HexToAddress(regInfo.Addr)
		sig := common.Hex2Bytes(regInfo.Sig)
		tx := registerTx(&ethAddr, sig)
		rawTx := encodeTransaction(tx)
		deliverTxRes := app.DeliverTx(rawTx)
		So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
		likeChainID := deliverTxRes.Data
		account.SaveBalance(mockCtx.GetMutableState(), identifierFromID(likeChainID), big.NewInt(100))

		app.EndBlock(abci.RequestEndBlock{
			Height: 1,
		})
		app.Commit()

		likeChainIDBase64 := "bDH8FUIuutKKr5CJwwZwL2dUC1M="
		likeChainIDParsed, _ := base64.StdEncoding.DecodeString(likeChainIDBase64)
		So(bytes.Compare(likeChainIDParsed, likeChainID), ShouldBeZeroValue)

		queryRes := app.Query(abci.RequestQuery{
			Path: "account_info",
			Data: []byte(likeChainIDBase64),
		})
		So(queryRes.Code, ShouldEqual, response.Success.Code)
		accountInfo := getAccountInfoRes(queryRes.Value)
		So(accountInfo, ShouldNotBeNil)
		So(accountInfo.Balance.Cmp(big.NewInt(100)), ShouldBeZeroValue)
		So(accountInfo.NextNonce, ShouldEqual, 1)

		Convey("Given a WithdrawTransaction from A to a certain address with value 1", func() {
			toEthAddr := common.HexToAddress("0x833a907efe57af3040039c90f4a59946a0bb3d47")
			sig := common.Hex2Bytes("d2354ea2e358bfd8e40d7afeaf6dbc79f6241d5517c398b5901f5162b7d9a09e58d2bdaaaf577ed28d1b871fea7a20572f2bf3865d6bad7e82687967c5cb63dd1c")
			tx := withdrawTx(likeChainID, nil, &toEthAddr, big.NewInt(1), big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The withdraw should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				packedTx := deliverTxRes.Data
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info should return the correct balance", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)
				})
				Convey("Then query withdraw_proof should return a proof", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path:   "withdraw_proof",
						Data:   packedTx,
						Height: 2,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					proof := iavl.RangeProof{}
					err := json.Unmarshal(queryRes.Value, &proof)
					So(err, ShouldBeNil)
					Convey("The proof should be corresponding to the withdraw tree hash", func() {
						err := proof.Verify(mockCtx.GetMutableState().GetAppHash()[:20])
						So(err, ShouldBeNil)
					})
				})
				Convey("Then query withdraw_proof with changed info should return fail", func() {
					packedTx[0]++
					queryRes := app.Query(abci.RequestQuery{
						Path:   "withdraw_proof",
						Data:   packedTx,
						Height: 2,
					})
					So(queryRes.Code, ShouldEqual, response.QueryWithdrawProofNotExist.Code)
				})
				Convey("But repeated withdraw with the same transaction should fail", func() {
					checkTxRes := app.CheckTx(rawTx)
					So(checkTxRes.Code, ShouldEqual, response.WithdrawCheckTxDuplicated.Code)
					deliverTxRes := app.DeliverTx(rawTx)
					So(deliverTxRes.Code, ShouldEqual, response.WithdrawDeliverTxDuplicated.Code)
				})
			})
		})

		Convey("Given a WithdrawTransaction from A's Ethereum address to a certain address with value 1", func() {
			fromEthAddr := common.HexToAddress("0x539c17e9e5fd1c8e3b7506f4a7d9ba0a0677eae9")
			toEthAddr := common.HexToAddress("0x833a907efe57af3040039c90f4a59946a0bb3d47")
			sig := common.Hex2Bytes("cfd63e8ff3991492c7eb56723ec12fdcc2e145b20c0de2a578ce63c268ad770f4f3361e27a8ae34fdf7b897f13a09b2e544eca7a8d533db28af42d54ff4df08d1c")
			tx := withdrawTx(nil, &fromEthAddr, &toEthAddr, big.NewInt(1), big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The withdraw should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				packedTx := deliverTxRes.Data
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info should return the correct balance", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(99)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)
				})
				Convey("Then query withdraw_proof should return a proof", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path:   "withdraw_proof",
						Data:   packedTx,
						Height: 2,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					proof := iavl.RangeProof{}
					err := json.Unmarshal(queryRes.Value, &proof)
					So(err, ShouldBeNil)
					Convey("The proof should be corresponding to the withdraw tree hash", func() {
						err := proof.Verify(mockCtx.GetMutableState().GetAppHash()[:20])
						So(err, ShouldBeNil)
					})
				})
			})
		})

		Convey("Given a WithdrawTransaction from A to a certain address with value 100", func() {
			toEthAddr := common.HexToAddress("0x833a907efe57af3040039c90f4a59946a0bb3d47")
			sig := common.Hex2Bytes("3b0ea1e2e032d01b559f6d27a92c6be0372fb4d5d54ee6707835b6f217d1fa7226e9d2e1180331dfd12a880639e98bc8aa10349fba1da467cb2784eddfa903d41b")
			tx := withdrawTx(likeChainID, nil, &toEthAddr, big.NewInt(100), big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The withdraw should success", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.Success.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.Success.Code)
				packedTx := deliverTxRes.Data
				app.EndBlock(abci.RequestEndBlock{
					Height: 2,
				})
				app.Commit()
				Convey("Then query account_info should return the correct balance", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path: "account_info",
						Data: []byte(likeChainIDBase64),
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					accountInfo := getAccountInfoRes(queryRes.Value)
					So(accountInfo, ShouldNotBeNil)
					So(accountInfo.Balance.Cmp(big.NewInt(0)), ShouldBeZeroValue)
					So(accountInfo.NextNonce, ShouldEqual, 2)
				})
				Convey("Then query withdraw_proof should return a proof", func() {
					queryRes := app.Query(abci.RequestQuery{
						Path:   "withdraw_proof",
						Data:   packedTx,
						Height: 2,
					})
					So(queryRes.Code, ShouldEqual, response.Success.Code)
					proof := iavl.RangeProof{}
					err := json.Unmarshal(queryRes.Value, &proof)
					So(err, ShouldBeNil)
					Convey("The proof should be corresponding to the withdraw tree hash", func() {
						err := proof.Verify(mockCtx.GetMutableState().GetAppHash()[:20])
						So(err, ShouldBeNil)
					})
				})
			})
		})

		Convey("Given a WithdrawTransaction from A to a certain address with value 101", func() {
			toEthAddr := common.HexToAddress("0x833a907efe57af3040039c90f4a59946a0bb3d47")
			sig := common.Hex2Bytes("d7abbd0ffeca27528cf28816faaf6b9e412f020d1f453250880071a7c3515fea12b1ac8594c7b893946efd723efe62915122e662da261da7336fce90623f7c8e1b")
			tx := withdrawTx(likeChainID, nil, &toEthAddr, big.NewInt(101), big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The withdraw should fail", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.WithdrawCheckTxNotEnoughBalance.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.WithdrawDeliverTxNotEnoughBalance.Code)
			})
		})

		Convey("Given a WithdrawTransaction with invalid signature", func() {
			toEthAddr := common.HexToAddress("0x833a907efe57af3040039c90f4a59946a0bb3d47")
			sig := common.Hex2Bytes("e828d630862be9e3564d0723c875ea93b1ec6be17c42f2a7345909d55f0b403024a1471b1000339e2a9f026d8e47d9f0afa856f899e671328b0fe63436e555911c")
			tx := withdrawTx(likeChainID, nil, &toEthAddr, big.NewInt(101), big.NewInt(0), 1, sig)
			rawTx := encodeTransaction(tx)
			Convey("The withdraw should fail with InvalidSignature", func() {
				checkTxRes := app.CheckTx(rawTx)
				So(checkTxRes.Code, ShouldEqual, response.WithdrawCheckTxInvalidSignature.Code)
				deliverTxRes := app.DeliverTx(rawTx)
				So(deliverTxRes.Code, ShouldEqual, response.WithdrawDeliverTxInvalidSignature.Code)
			})
		})
	})
}
