package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	tmRPC "github.com/tendermint/tendermint/rpc/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/likecoin/likechain/abci/txs"
	"github.com/likecoin/likechain/abci/types"
)

const spaces = "                                                               "

func shouldCallString(v reflect.Value) bool {
	if v.Type().PkgPath() == "github.com/likecoin/likechain/abci/types" {
		return true
	}
	if strings.HasSuffix(v.Type().Name(), "Signature") {
		return true
	}
	return false
}

func simpleValueToString(v reflect.Value) (string, bool) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true", true
		}
		return "false", true
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return fmt.Sprintf("%d", v.Int()), true
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint()), true
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return fmt.Sprintf("%f", v.Float()), true
	case reflect.String:
		return fmt.Sprintf("\"%s\"", v.String()), true
	default:
		if shouldCallString(v) {
			m := v.MethodByName("String")
			if m.IsValid() {
				return m.Call(nil)[0].String(), true
			} else if v.CanAddr() {
				m = v.Addr().MethodByName("String")
				if m.IsValid() {
					return m.Call(nil)[0].String(), true
				}
			}
		}
		switch v.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:
			return simpleValueToString(v.Elem())
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			t := v.Type().Elem()
			if t.Kind() == reflect.Uint8 {
				return hex.EncodeToString(v.Slice(0, v.Len()).Bytes()), true
			}
			fallthrough
		default:
			return "", false
		}
	}
}

func printValue(v reflect.Value, indent int) {
	s, ok := simpleValueToString(v)
	if ok {
		fmt.Printf("%s\n", s)
		return
	}
	switch v.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		fmt.Println()
		for i := 0; i < v.Len(); i++ {
			fmt.Printf("%.*s%d: ", indent+4, spaces, i)
			printValue(v.Index(i), indent+4)
		}
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		printValue(v.Elem(), indent)
	case reflect.Struct:
		t := v.Type()
		fmt.Printf("%s { \n", t.Name())
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			fmt.Printf("%.*s%s: ", indent+4, spaces, f.Name)
			printValue(v.Field(i), indent+4)
		}
		fmt.Printf("%.*s}\n", indent, spaces)
	default:
		fmt.Printf("%.*s%s\n", indent, spaces, v.String())
	}
}

var decodeCmd = &cobra.Command{
	Use:   "decode [tx-hash or base64-raw-tx]",
	Short: "decode LikeChain transactions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txHashMatcher := regexp.MustCompile("^[a-fA-F0-9]{64}$")
		log.
			WithField("args", args[0]).
			Debug("Read input args")

		var rawTx []byte
		if txHashMatcher.Match([]byte(args[0])) {
			tmEndPoint := viper.GetString("tmEndPoint")
			log.
				WithField("args", args[0]).
				WithField("tm_endpoint", tmEndPoint).
				Debug("Treated args as txHash")
			txHashHex := args[0]
			txHash, _ := hex.DecodeString(txHashHex)
			tmClient := tmRPC.NewHTTP(tmEndPoint, "/websocket")
			txResult, err := tmClient.Tx(txHash, false)
			if err != nil {
				log.
					WithField("tx_hash", txHashHex).
					WithField("tm_endpoint", tmEndPoint).
					WithError(err).
					Panic("Cannot get transaction from Tendermint endpoint")
			}
			rawTx = txResult.Tx
		} else {
			log.
				WithField("args", args[0]).
				Debug("Treated args as rawTx")
			rawTxBase64 := args[0]
			var err error
			rawTx, err = base64.StdEncoding.DecodeString(rawTxBase64)
			if err != nil {
				log.
					WithField("raw_tx_base64", rawTxBase64).
					WithError(err).
					Panic("Cannot decode rawTx as base64")
			}
		}
		log.
			WithField("raw_tx", rawTx).
			Debug("Got rawTx")
		var tx txs.Transaction
		err := types.AminoCodec().UnmarshalBinaryLengthPrefixed(rawTx, &tx)
		if err != nil {
			log.
				WithField("raw_tx", rawTx).
				WithError(err).
				Panic("Cannot deserialize rawTx")
		}
		printValue(reflect.ValueOf(tx), 0)
	},
}
