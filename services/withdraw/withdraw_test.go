package withdraw

import (
	"bytes"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"

	"github.com/likecoin/likechain/services/tendermint"

	. "github.com/smartystreets/goconvey/convey"
)

func testContractProof(signedHeader *types.SignedHeader, tmToEthAddr map[int]common.Address, validators []types.Validator) {
	// Below mimics the logic in Relay contract
	contractProof := genContractProofPayload(signedHeader, tmToEthAddr)
	i := 0
	suffixLen := int(contractProof[i])
	i++
	suffix := contractProof[i : i+suffixLen]
	i += suffixLen
	votesCount := int(contractProof[i])
	So(votesCount, ShouldEqual, len(validators))
	i++
	for j := 0; j < votesCount; j++ {
		timeLen := int(contractProof[i])
		i++
		timeBytes := contractProof[i : i+timeLen]
		i += timeLen
		sig := contractProof[i : i+65]
		i += 65
		buf := new(bytes.Buffer)
		buf.WriteByte(0) // length place holder
		buf.WriteByte(8)
		buf.WriteByte(2)
		if signedHeader.Header.Height > 0 {
			buf.WriteByte(0x11)
			amino.EncodeInt64(buf, signedHeader.Header.Height)
		}
		if signedHeader.Commit.Precommits[0].Round > 0 {
			buf.WriteByte(0x19)
			amino.EncodeInt64(buf, int64(signedHeader.Commit.Precommits[0].Round))
		}
		buf.Write(timeBytes)
		buf.Write(suffix)
		signBytes := buf.Bytes()
		signBytes[0] = byte(len(signBytes) - 1)
		So(signBytes, ShouldResemble, signedHeader.Commit.Precommits[j].SignBytes(signedHeader.Header.ChainID))
		recoveryID := sig[0:1]
		rs := sig[1:]
		ethSig := make([]byte, 0, 65)
		ethSig = append(ethSig, rs...)
		ethSig = append(ethSig, recoveryID...)
		ethSig[64] -= 27
		recoveredPubKeyBytes, err := ethCrypto.Ecrecover(tmhash.Sum(signBytes), ethSig)
		So(err, ShouldBeNil)
		recoveredPubKey, err := ethCrypto.UnmarshalPubkey(recoveredPubKeyBytes)
		So(err, ShouldBeNil)
		recoveredAddr := ethCrypto.PubkeyToAddress(*recoveredPubKey)
		So(recoveredAddr, ShouldResemble, tmToEthAddr[j])
	}
	blockHash := suffix[4:36]
	So(blockHash, ShouldResemble, []byte(signedHeader.Commit.BlockID.Hash))
	appHashLen := int(contractProof[i])
	i++
	appHash := contractProof[i : i+appHashLen]
	So(appHash, ShouldResemble, []byte(signedHeader.Header.AppHash))
	i += appHashLen
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(appHashLen))
	buf.Write(appHash)
	hash := tmhash.Sum(buf.Bytes())
	buf = new(bytes.Buffer)
	buf.WriteByte(32)
	buf.Write(hash)
	buf.WriteByte(32)
	buf.Write(contractProof[i : i+32])
	hash = tmhash.Sum(buf.Bytes())
	i += 32
	buf = new(bytes.Buffer)
	buf.WriteByte(32)
	buf.Write(hash)
	buf.WriteByte(32)
	buf.Write(contractProof[i : i+32])
	hash = tmhash.Sum(buf.Bytes())
	i += 32
	buf = new(bytes.Buffer)
	buf.WriteByte(32)
	buf.Write(contractProof[i : i+32])
	buf.WriteByte(32)
	buf.Write(hash)
	hash = tmhash.Sum(buf.Bytes())
	i += 32
	buf = new(bytes.Buffer)
	buf.WriteByte(32)
	buf.Write(contractProof[i : i+32])
	buf.WriteByte(32)
	buf.Write(hash)
	hash = tmhash.Sum(buf.Bytes())
	So(hash, ShouldResemble, blockHash)
}

func TestContractProofPayload(t *testing.T) {
	Convey("Given a Tendermint signed block header", t, func() {
		const validatorsCount = 10
		validators := make([]types.Validator, validatorsCount)
		privKeys := make([]crypto.PrivKey, validatorsCount)
		for i := 0; i < validatorsCount; i++ {
			privKey := secp256k1.GenPrivKey()
			privKeys[i] = privKey
			validators[i] = *types.NewValidator(privKey.PubKey(), cmn.RandInt63n(0xFFFFFFFF))
		}
		tmToEthAddr := tendermint.MapValidatorIndexToEthAddr(validators)
		header := types.Header{
			Version: version.Consensus{
				Block: 123,
				App:   456,
			},
			ChainID:  "some chain ID",
			Height:   789,
			Time:     time.Now(),
			NumTxs:   10,
			TotalTxs: 1337,
			LastBlockID: types.BlockID{
				Hash: cmn.RandBytes(32),
				PartsHeader: types.PartSetHeader{
					Total: 1,
					Hash:  cmn.RandBytes(32),
				},
			},
			LastCommitHash:     cmn.RandBytes(32),
			DataHash:           cmn.RandBytes(32),
			ValidatorsHash:     cmn.RandBytes(32),
			NextValidatorsHash: cmn.RandBytes(32),
			ConsensusHash:      cmn.RandBytes(32),
			AppHash:            cmn.RandBytes(64),
			LastResultsHash:    cmn.RandBytes(32),
			EvidenceHash:       cmn.RandBytes(32),
			ProposerAddress:    cmn.RandBytes(20),
		}
		blockID := types.BlockID{
			Hash: header.Hash(),
			PartsHeader: types.PartSetHeader{
				Total: 1,
				Hash:  cmn.RandBytes(32),
			},
		}
		precommits := make([]*types.Vote, validatorsCount)
		for i := 0; i < validatorsCount; i++ {
			vote := types.Vote{
				Type:             2,
				Height:           header.Height,
				Round:            0,
				Timestamp:        time.Now().Add(time.Duration(cmn.RandInt31()) * time.Nanosecond),
				BlockID:          blockID,
				ValidatorAddress: validators[i].Address,
				ValidatorIndex:   i,
			}
			var err error
			vote.Signature, err = privKeys[i].Sign(vote.SignBytes(header.ChainID))
			So(err, ShouldBeNil)
			So(vote.Verify(header.ChainID, privKeys[i].PubKey()), ShouldBeNil)
			precommits[i] = &vote
		}
		signedHeader := types.SignedHeader{
			Header: &header,
			Commit: &types.Commit{
				BlockID:    blockID,
				Precommits: precommits,
			},
		}
		Convey("If round number is 0", func() {
			Convey("Contract proof should be valid", func() {
				testContractProof(&signedHeader, tmToEthAddr, validators)
			})
		})
		Convey("If round number is not 0", func() {
			round := int(cmn.RandInt31())
			if round == 0 {
				// Really bad luck
				round++
			}
			for i := 0; i < validatorsCount; i++ {
				vote := precommits[i]
				vote.Round = round
				var err error
				vote.Signature, err = privKeys[i].Sign(vote.SignBytes(header.ChainID))
				So(err, ShouldBeNil)
				So(vote.Verify(header.ChainID, privKeys[i].PubKey()), ShouldBeNil)
			}
			signedHeader := types.SignedHeader{
				Header: &header,
				Commit: &types.Commit{
					BlockID:    blockID,
					Precommits: precommits,
				},
			}
			Convey("Contract proof should be valid", func() {
				testContractProof(&signedHeader, tmToEthAddr, validators)
			})
		})
	})
}
