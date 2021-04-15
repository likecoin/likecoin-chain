package oldgenesis

import (
	"github.com/tendermint/tendermint/crypto/ed25519"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"

	"github.com/pkg/errors"
)

// ConsensusParams contains consensus critical parameters that determine the
// validity of blocks.
type ConsensusParams struct {
	Block     tmproto.BlockParams     `json:"block"`
	Evidence  EvidenceParams          `json:"evidence"`
	Validator tmproto.ValidatorParams `json:"validator"`
}

// EvidenceParams determine how we handle evidence of malfeasance.
type EvidenceParams struct {
	MaxAge int64 `json:"max_age"` // only accept new evidence more recent than this
}

// DefaultEvidenceParams Params returns a default EvidenceParams.
func DefaultEvidenceParams() EvidenceParams {
	return EvidenceParams{
		MaxAge: 100000, // 27.8 hrs at 1block/s
	}
}

func DefaultConsensusParams() *ConsensusParams {
	return &ConsensusParams{
		types.DefaultBlockParams(),
		DefaultEvidenceParams(),
		types.DefaultValidatorParams(),
	}
}

const (
	ABCIPubKeyTypeEd25519 = "ed25519"
)

var ABCIPubKeyTypesToAminoNames = map[string]string{
	ABCIPubKeyTypeEd25519: ed25519.PubKeyName,
}

// Validate validates the ConsensusParams to ensure all values are within their
// allowed limits, and returns an error if they are not.
func (params *ConsensusParams) Validate() error {
	if params.Block.MaxBytes <= 0 {
		return errors.Errorf("Block.MaxBytes must be greater than 0. Got %d",
			params.Block.MaxBytes)
	}
	if params.Block.MaxBytes > types.MaxBlockSizeBytes {
		return errors.Errorf("Block.MaxBytes is too big. %d > %d",
			params.Block.MaxBytes, types.MaxBlockSizeBytes)
	}

	if params.Block.MaxGas < -1 {
		return errors.Errorf("Block.MaxGas must be greater or equal to -1. Got %d",
			params.Block.MaxGas)
	}

	if params.Block.TimeIotaMs <= 0 {
		return errors.Errorf("Block.TimeIotaMs must be greater than 0. Got %v",
			params.Block.TimeIotaMs)
	}

	if params.Evidence.MaxAge <= 0 {
		return errors.Errorf("EvidenceParams.MaxAge must be greater than 0. Got %d",
			params.Evidence.MaxAge)
	}

	if len(params.Validator.PubKeyTypes) == 0 {
		return errors.New("len(Validator.PubKeyTypes) must be greater than 0")
	}

	// Check if keyType is a known ABCIPubKeyType
	for i := 0; i < len(params.Validator.PubKeyTypes); i++ {
		keyType := params.Validator.PubKeyTypes[i]
		if _, ok := ABCIPubKeyTypesToAminoNames[keyType]; !ok {
			return errors.Errorf("params.Validator.PubKeyTypes[%d], %s, is an unknown pubkey type",
				i, keyType)
		}
	}

	return nil
}
