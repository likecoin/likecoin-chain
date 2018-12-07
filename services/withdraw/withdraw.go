package withdraw

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/tendermint/tendermint/crypto"
	tmRPC "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"

	"github.com/likecoin/likechain/services/abi/relay"
	logger "github.com/likecoin/likechain/services/log"
	"github.com/likecoin/likechain/services/tendermint"
)

var log = logger.L

// type AppHashContractProof struct {
// 	Height     uint64
// 	Round      uint64
// 	Payload struct {
// 		SuffixLen  uint8
// 		Suffix     []byte
// 		VotesCount uint8
// 		Votes      []struct {
// 			TimeLen uint8
// 			Time    []byte
// 			Sig     [65]byte
// 		}
// 		AppHashLen   uint8
// 		AppHash      []byte
// 		AppHashProof [4][32]byte
// 	}
// }

func encodeTimestamp(vote *types.CanonicalVote) []byte {
	cdc := types.GetCodec()
	buf := new(bytes.Buffer)
	// Field 4, typ3 = variable length struct (2), so field index = 00100|010 = 0x22
	// See comments in the Relay smart contract  for more details
	buf.WriteByte(0x22)
	buf.Write(cdc.MustMarshalBinaryLengthPrefixed(vote.Timestamp))
	return buf.Bytes()
}

func encodeSuffix(vote *types.CanonicalVote) []byte {
	cdc := types.GetCodec()
	buf := new(bytes.Buffer)
	// Field 5, typ3 = variable length struct (2), so field index = 00101|010 = 0x2A
	// See comments in the Relay smart contract  for more details
	buf.WriteByte(0x2A)
	buf.Write(cdc.MustMarshalBinaryLengthPrefixed(vote.BlockID))
	// Field 6, typ3 = variable length struct (2), so field index = 00110|010 = 0x32
	buf.WriteByte(0x32)
	buf.Write(cdc.MustMarshalBinaryBare(vote.ChainID))
	return buf.Bytes()
}

func genContractProofPayload(signedHeader *types.SignedHeader, tmToEthAddr map[int]common.Address) []byte {
	header := signedHeader.Header
	rawVotes := signedHeader.Commit.Precommits
	votes := []*types.Vote{}

	for _, vote := range rawVotes {
		if vote != nil {
			votes = append(votes, vote)
		}
	}

	votesCount := len(votes)
	if votesCount == 0 {
		return nil
	}

	cv := types.CanonicalizeVote(header.ChainID, votes[0])
	suffix := encodeSuffix(&cv)

	buf := new(bytes.Buffer)
	buf.WriteByte(uint8(len(suffix)))
	buf.Write(suffix)
	buf.WriteByte(uint8(votesCount))

	for _, vote := range votes {
		cv := types.CanonicalizeVote(header.ChainID, vote)
		time := encodeTimestamp(&cv)
		buf.WriteByte(uint8(len(time)))
		buf.Write(time)

		signBytes := vote.SignBytes(header.ChainID)
		ethAddr := tmToEthAddr[vote.ValidatorIndex]
		ethSig := tendermint.SignatureToEthereumSig(vote.Signature, crypto.Sha256(signBytes), ethAddr)
		buf.Write(ethSig[64:])
		buf.Write(ethSig[:64])
	}

	buf.WriteByte(uint8(len(header.AppHash)))
	buf.Write(header.AppHash)
	_, proof := Proof(header)
	for _, pf := range proof {
		buf.Write(pf)
	}
	return buf.Bytes()
}

func waitForReceipt(ethClient *ethclient.Client, txHash common.Hash) (*ethTypes.Receipt, error) {
	for {
		receipt, err := ethClient.TransactionReceipt(context.Background(), txHash)
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		time.Sleep(15 * time.Second)
	}
}

func doWithdraw(tmClient *tmRPC.HTTP, ethClient *ethclient.Client, auth *bind.TransactOpts, contractAddr common.Address, callData withdrawCallData) {
	contract, err := relay.NewRelay(contractAddr, ethClient)
	if err != nil {
		panic(err)
	}

	log.
		WithField("withdraw_info", common.Bytes2Hex(callData.WithdrawInfo)).
		WithField("contract_proof", common.Bytes2Hex(callData.ContractProof)).
		Info("Calling withdraw on Ethereum")
	tx, err := contract.Withdraw(auth, callData.WithdrawInfo, callData.ContractProof)
	if err != nil {
		panic(err)
	}

	receipt, err := waitForReceipt(ethClient, tx.Hash())
	if err != nil {
		panic(err)
	}
	log.
		WithField("gas_used", receipt.GasUsed).
		WithField("status", receipt.Status).
		Info("withdraw call executed on Ethereum")
}

func getContractHeight(ethClient *ethclient.Client, contractAddr common.Address) int64 {
	contract, err := relay.NewRelay(contractAddr, ethClient)
	if err != nil {
		panic(err)
	}
	height, err := contract.LatestBlockHeight(nil)
	if err != nil {
		panic(err)
	}
	return height.Int64()
}

func commitWithdrawHash(tmClient *tmRPC.HTTP, ethClient *ethclient.Client, auth *bind.TransactOpts, contractAddr common.Address, height int64) {
	validators := tendermint.GetValidators(tmClient)
	tmToEthAddr := tendermint.MapValidatorIndexToEthAddr(validators)

	signedHeader := tendermint.GetSignedHeader(tmClient, height)

	log.
		WithField("header_block_hash", signedHeader.Commit.BlockID.Hash).
		Debug("Got SignedHeader")
	contractPayload := genContractProofPayload(&signedHeader, tmToEthAddr)
	contract, err := relay.NewRelay(contractAddr, ethClient)
	if err != nil {
		panic(err)
	}

	round := uint64(signedHeader.Commit.Round())
	log.
		WithField("height", height).
		WithField("round", round).
		WithField("contract_payload", common.Bytes2Hex(contractPayload)).
		Info("Calling commitWithdrawHash on Ethereum")

	tx, err := contract.CommitWithdrawHash(auth, uint64(height), round, contractPayload)
	if err != nil {
		panic(err)
	}

	receipt, err := waitForReceipt(ethClient, tx.Hash())
	if err != nil {
		panic(err)
	}
	log.
		WithField("gas_used", receipt.GasUsed).
		WithField("status", receipt.Status).
		Info("commitWithdrawHash call executed on Ethereum")
}

type withdrawCallData struct {
	WithdrawInfo  []byte
	ContractProof []byte
}

func getWithdrawCallDataArr(tmClient *tmRPC.HTTP, lastHeight, newHeight int64) []withdrawCallData {
	log.
		WithField("last_height", lastHeight).
		WithField("new_height", newHeight).
		Info("Searching withdraws on LikeChain")
	queryString := fmt.Sprintf("withdraw.height>%d AND withdraw.height<=%d", lastHeight, newHeight)
	// TODO: may need pagination
	searchResult, err := tmClient.TxSearch(queryString, true, 1, 100)
	if err != nil {
		panic(err)
	}
	if searchResult.TotalCount <= 0 {
		log.
			WithField("new_height", newHeight).
			Info("No withdraw search result")
		return nil
	}
	callDataArr := make([]withdrawCallData, searchResult.TotalCount)
	for i := 0; i < searchResult.TotalCount; i++ {
		packedTx := searchResult.Txs[i].TxResult.Data
		log.
			WithField("result_index", i).
			WithField("tx_hash", searchResult.Txs[i].Hash).
			WithField("packed_tx", common.Bytes2Hex(packedTx)).
			Debug("Withdraw search result")
		queryResult, err := tmClient.ABCIQueryWithOptions("withdraw_proof", packedTx, tmRPC.ABCIQueryOptions{Height: newHeight})
		if err != nil {
			log.
				WithField("packed_tx", common.Bytes2Hex(packedTx)).
				WithError(err).
				Panic("Cannot get withdraw_proof from LikeChain")
		}
		proof := ParseRangeProof(queryResult.Response.Value)
		if proof == nil {
			log.
				WithField("range_proof_json", string(queryResult.Response.Value)).
				Panic("Cannot parse RangeProof")
		}
		log.
			WithField("root_hash", common.Bytes2Hex(proof.ComputeRootHash())).
			Debug("Computed RangeProof root hash")
		contractProof := proof.ContractProof()
		callDataArr[i] = withdrawCallData{packedTx, contractProof}
	}
	return callDataArr
}

// Run starts the subscription to the withdraws on LikeChain and commits proofs onto Ethereum
func Run(tmClient *tmRPC.HTTP, ethClient *ethclient.Client, auth *bind.TransactOpts, contractAddr common.Address) {
	lastHeight := getContractHeight(ethClient, contractAddr)
	for ; ; time.Sleep(time.Minute) {
		// TODO: load lastHeight from database?
		newHeight := tendermint.GetHeight(tmClient)
		if newHeight == lastHeight {
			log.
				WithField("last_height", lastHeight).
				Info("No new LikeChain block since last height")
			continue
		}
		withdrawCallDataArr := getWithdrawCallDataArr(tmClient, lastHeight, newHeight)
		if len(withdrawCallDataArr) <= 0 {
			continue
		}
		contractHeight := getContractHeight(ethClient, contractAddr)
		if contractHeight < newHeight {
			commitWithdrawHash(tmClient, ethClient, auth, contractAddr, newHeight)
		} else if contractHeight > newHeight {
			log.
				WithField("contract_height", contractHeight).
				WithField("new_height", newHeight).
				Panic("New height is less than contract height")
		}
		// TODO: save callDataArr in database
		// TODO: save lastHeight in database?
		lastHeight = newHeight
		for _, callData := range withdrawCallDataArr {
			doWithdraw(tmClient, ethClient, auth, contractAddr, callData)
		}
	}
}
