package init

import (
	"fmt"
	"os"

	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

// Run starts the init process, generates configs and keys, then collect generated info to generate genesis file
// and Docker Compose file
func Run(profileDir, dockerDir string, nodeCount uint) {
	// 1. initialize configs and keys for each node
	// 2. collect keys to generate genesis
	// 3. write genesis to each node
	// 4. generate docker-compose.yml and docker-compose.production.yml
	pubInfos := make([]publicInfo, 0, nodeCount)
	for nodeIndex := uint(1); nodeIndex <= nodeCount; nodeIndex++ {
		configDir := fmt.Sprintf("%s/%d/config", profileDir, nodeIndex)
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			panic(err)
		}
		stateDir := fmt.Sprintf("%s/%d/data", profileDir, nodeIndex)
		err = os.MkdirAll(stateDir, 0755)
		if err != nil {
			panic(err)
		}
		pubInfo := genKeys(configDir, stateDir)
		pubInfos = append(pubInfos, pubInfo)
		genConfig(configDir, nodeIndex)
	}
	genesis := genGenesis(pubInfos)
	for nodeIndex := uint(1); nodeIndex <= nodeCount; nodeIndex++ {
		genesisPath := fmt.Sprintf("%s/%d/config/genesis.json", profileDir, nodeIndex)
		err := genesis.SaveAs(genesisPath)
		if err != nil {
			panic(err)
		}
	}
	genDockerComposeFiles(dockerDir, pubInfos)
}

func init() {
	cryptoAmino.RegisterAmino(cdc)
}
