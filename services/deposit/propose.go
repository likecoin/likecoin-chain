package deposit

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/likecoin/likechain/abci/query"
	"github.com/likecoin/likechain/abci/response"
	"github.com/likecoin/likechain/abci/state/deposit"
	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"

	"github.com/likecoin/likechain/services/eth"
	"github.com/likecoin/likechain/services/tendermint"
	"github.com/likecoin/likechain/services/utils"

	tmRPC "github.com/tendermint/tendermint/rpc/client"
)

var proposeLock = &sync.Mutex{}

func fillSig(tx *txs.DepositTransaction, privKey *ecdsa.PrivateKey) {
	tx.Proposal.Sort()
	jsonMap := tx.GenerateJSONMap()
	hash, err := jsonMap.Hash()
	if err != nil {
		log.
			WithField("tx", tx).
			WithError(err).
			Panic("Cannot hash deposit transaction")
	}
	sig, err := crypto.Sign(hash, privKey)
	if err != nil {
		log.
			WithField("tx", tx).
			WithError(err).
			Panic("Cannot sign deposit transaction")
	}
	sig[64] += 27
	jsonSig := txs.DepositJSONSignature{}
	copy(jsonSig.JSONSignature[:], sig)
	tx.Sig = &jsonSig
}

// Propose a deposit proposal to LikeChain, return true if succeed, false if timeout, panic otherwise
func propose(tmClient *tmRPC.HTTP, tmPrivKey *ecdsa.PrivateKey, proposal deposit.Proposal) bool {
	proposeLock.Lock()
	defer proposeLock.Unlock()
	log.
		WithField("block_number", proposal.BlockNumber).
		WithField("event_count", len(proposal.Inputs)).
		Info("Proposing new proposal")
	ethAddr := crypto.PubkeyToAddress(tmPrivKey.PublicKey)
	addr, err := types.NewAddress(ethAddr[:])
	if err != nil {
		log.
			WithField("eth_addr", ethAddr.Hex()).
			WithError(err).
			Panic("Cannot convert Ethereum address to LikeChain address")
	}
	queryResult, err := tmClient.ABCIQuery("account_info", []byte(addr.String()))
	if err != nil {
		log.
			WithField("addr", addr.String()).
			WithError(err).
			Panic("Cannot query account info from ABCI")
	}
	accInfo := query.GetAccountInfoRes(queryResult.Response.Value)
	if accInfo == nil {
		log.
			WithField("account_info_res", string(queryResult.Response.Value)).
			WithField("account_info_res_raw", queryResult.Response.Value).
			Panic("Cannot parse account info result")
	}
	log.
		WithField("nonce", accInfo.NextNonce).
		Debug("Got account info")
	tx := &txs.DepositTransaction{
		Proposer: addr,
		Proposal: proposal,
		Nonce:    accInfo.NextNonce,
	}
	fillSig(tx, tmPrivKey)
	rawTx := txs.EncodeTx(tx)
	log.
		WithField("raw_tx", common.Bytes2Hex(rawTx)).
		Debug("Broadcasting transaction onto LikeChain")
	result, err := tendermint.BroadcastTxCommit(tmClient, rawTx)
	if err != nil {
		switch err.(type) {
		case tendermint.BroadcastTimeout:
			return false
		default:
			log.
				WithField("raw_tx", common.Bytes2Hex(rawTx)).
				WithError(err).
				Panic("Broadcast transaction onto LikeChain failed")
		}
	}
	if result.Code != response.Success.Code {
		switch result.Code {
		case response.DepositDoubleApproval.ToResponseCheckTx().Code:
			fallthrough
		case response.DepositDoubleApproval.ToResponseDeliverTx().Code:
			fallthrough
		case response.DepositAlreadyExecuted.ToResponseCheckTx().Code:
			fallthrough
		case response.DepositAlreadyExecuted.ToResponseDeliverTx().Code:
			log.
				WithField("code", result.Code).
				WithField("info", result.Info).
				WithField("log", result.Log).
				Info("Deposit transaction unnecessary and rejected, skipping")
		default:
			log.
				WithField("code", result.Code).
				WithField("info", result.Info).
				WithField("log", result.Log).
				Panic("Deposit transaction executed but failed")
		}
	} else {
		log.Info("Successfully executed deposit transaction onto LikeChain")
	}
	return true
}

func scanAndProposeForRange(config *Config, from, to uint64) (abandonedBlocks []uint64) {
	log.
		WithField("begin_block", from).
		WithField("end_block", to).
		Debug("Searching blocks in range")
	proposals := eth.GetTransfersFromBlocks(
		config.LoadBalancer,
		config.TokenAddr,
		config.RelayAddr,
		from,
		to,
	)
	if len(proposals) == 0 {
		log.
			WithField("begin_block", from).
			WithField("end_block", to).
			Info("No transfer events in range")
		return nil
	}
	for _, proposal := range proposals {
		log.
			WithField("block", proposal.BlockNumber).
			Info("Proposing proposal")
		utils.RetryIfPanic(5, func() {
			succ := propose(config.TMClient, config.TMPrivKey, proposal)
			if !succ {
				log.
					WithField("block", proposal.BlockNumber).
					Warn("Propose block timeout, putting block into lower priority poller")
				abandonedBlocks = append(abandonedBlocks, proposal.BlockNumber)
			}
		})
	}
	return abandonedBlocks
}
