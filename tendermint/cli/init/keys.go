package init

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

type publicInfo struct {
	NodeID p2p.ID
	PubKey crypto.PubKey
}

func genKeys(keyDir, stateDir string) publicInfo {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	address := pubKey.Address()
	keyFilePath := fmt.Sprintf("%s/priv_validator_key.json", keyDir)
	stateFilePath := fmt.Sprintf("%s/priv_validator_state.json", stateDir)
	pv := privval.GenFilePV(keyFilePath, stateFilePath)
	pv.Key.PrivKey = privKey
	pv.Key.PubKey = pubKey
	pv.Key.Address = address
	pv.Save()
	nodeKeyPath := fmt.Sprintf("%s/node_key.json", keyDir)
	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyPath)
	if err != nil {
		panic(err)
	}
	pubInfo := publicInfo{
		NodeID: nodeKey.ID(),
		PubKey: pubKey,
	}
	jsonBytes, err := cdc.MarshalJSONIndent(&pubInfo, "", "  ")
	if err != nil {
		panic(err)
	}
	publicInfoPath := fmt.Sprintf("%s/public.json", keyDir)
	err = common.WriteFileAtomic(publicInfoPath, jsonBytes, 0644)
	if err != nil {
		panic(err)
	}
	return pubInfo
}
