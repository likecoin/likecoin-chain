package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

func main() {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)

	var typ string
	var keyOutputDir string
	var stateOutputDir string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize config files for Tendermint",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var privKey crypto.PrivKey
			switch typ {
			case "secp256k1":
				privKey = secp256k1.GenPrivKey()
			case "ed25519":
				privKey = ed25519.GenPrivKey()
			default:
				panic(fmt.Sprintf("Unknown key type %s", typ))
			}
			pubKey := privKey.PubKey()
			address := pubKey.Address()
			keyFilePath := keyOutputDir + "/priv_validator_key.json"
			stateFilePath := stateOutputDir + "/priv_validator_state.json"
			pv := privval.GenFilePV(keyFilePath, stateFilePath)
			pv.Key.PrivKey = privKey
			pv.Key.PubKey = pubKey
			pv.Key.Address = address
			pv.Save()
			nodeKey, err := p2p.LoadOrGenNodeKey(keyOutputDir + "/node_key.json")
			if err != nil {
				panic(err)
			}
			publicInfo := struct {
				NodeID p2p.ID
				PubKey crypto.PubKey
			}{
				NodeID: nodeKey.ID(),
				PubKey: pubKey,
			}
			jsonBytes, err := cdc.MarshalJSONIndent(&publicInfo, "", "  ")
			if err != nil {
				panic(err)
			}
			err = common.WriteFileAtomic(keyOutputDir+"/public.json", jsonBytes, 0755)
			if err != nil {
				panic(err)
			}
		},
	}

	initCmd.Flags().StringVar(&typ, "type", "secp256k1", "private key type [secp256k1, ed25519]")
	initCmd.Flags().StringVar(&keyOutputDir, "key_output_dir", "./config", "output directory for generated key files")
	initCmd.Flags().StringVar(&stateOutputDir, "state_output_dir", "./data", "output directory for generated state file")

	initCmd.Execute()
}
