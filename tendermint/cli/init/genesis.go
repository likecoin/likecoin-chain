package init

import (
	"fmt"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

func genGenesis(pubInfos []publicInfo) types.GenesisDoc {
	consensusParams := types.DefaultConsensusParams()
	consensusParams.Validator.PubKeyTypes = []string{types.ABCIPubKeyTypeSecp256k1}
	genDoc := types.GenesisDoc{
		ChainID:         fmt.Sprintf("like-%v", cmn.RandStr(6)),
		GenesisTime:     tmtime.Now(),
		ConsensusParams: consensusParams,
	}
	validators := make([]types.GenesisValidator, 0, len(pubInfos))
	for i, pubInfo := range pubInfos {
		validators = append(validators, types.GenesisValidator{
			Name:    fmt.Sprintf("tendermint-%d", i+1),
			Address: pubInfo.PubKey.Address(),
			PubKey:  pubInfo.PubKey,
			Power:   10,
		})
	}
	genDoc.Validators = validators
	return genDoc
}
