package init

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/config"
)

func genConfig(configDir string, nodeIndex uint) {
	c := config.DefaultConfig()
	c.TxIndex.IndexAllTags = true
	c.Consensus.BlockTimeIota = 0
	c.Consensus.CreateEmptyBlocksInterval = time.Second * 1800
	c.BaseConfig.Moniker = fmt.Sprintf("tendermint-%d", nodeIndex)
	configFilePath := fmt.Sprintf("%s/config.toml", configDir)
	config.WriteConfigFile(configFilePath, c)
}
