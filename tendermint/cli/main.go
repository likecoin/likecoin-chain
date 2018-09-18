package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/privval"
)

func main() {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)

	var typ string
	var outputDir string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize config files for Tendermint",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			pv := privval.FilePV{}
			switch typ {
			case "secp256k1":
				pv.PrivKey = secp256k1.GenPrivKey()
			case "ed25519":
				pv.PrivKey = ed25519.GenPrivKey()
			default:
				panic(fmt.Sprintf("Unknown key type %s", typ))
			}
			pv.PubKey = pv.PrivKey.PubKey()
			pv.Address = pv.PubKey.Address()
			jsonBytes, err := cdc.MarshalJSONIndent(&pv, "", "  ")
			if err != nil {
				panic(err)
			}
			os.Mkdir(outputDir, 0755)
			err = common.WriteFileAtomic(outputDir+"/priv_validator.json", jsonBytes, 0600)
			if err != nil {
				panic(err)
			}
		},
	}

	initCmd.Flags().StringVar(&typ, "type", "secp256k1", "private key type [secp256k1, ed25519]")
	initCmd.Flags().StringVar(&outputDir, "output_dir", ".", "output directory for generated files")

	initCmd.Execute()
}
