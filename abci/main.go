package main

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/db"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

type NumberAdderApplication struct {
	types.BaseApplication

	db db.DB
}

func (app *NumberAdderApplication) rawHeight() []byte {
	return app.db.Get([]byte("__BLOCK_HEIGHT__"))
}

func (app *NumberAdderApplication) rawState() []byte {
	return app.db.Get([]byte("__STATE__"))
}

func (app *NumberAdderApplication) AppHash() []byte {
	if app.Height() == 0 {
		return nil
	}
	rawState := app.rawState()
	return ethCrypto.Keccak256(rawState)
}

func (app *NumberAdderApplication) Height() uint64 {
	rawHeight := app.rawHeight()
	if len(rawHeight) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(rawHeight)
}

func (app *NumberAdderApplication) State() *big.Int {
	rawState := app.rawState()
	state := big.NewInt(0)
	if len(rawState) == 0 {
		return state
	}
	return state.SetBytes(rawState)
}

func (app *NumberAdderApplication) SetHeight(newHeight uint64) {
	rawHeight := make([]byte, 8)
	binary.BigEndian.PutUint64(rawHeight, newHeight)
	app.db.Set([]byte("__BLOCK_HEIGHT__"), rawHeight)
}

func (app *NumberAdderApplication) SetState(newState *big.Int) {
	app.db.Set([]byte("__STATE__"), newState.Bytes())
}

func (app *NumberAdderApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	inputInt := big.NewInt(0)
	inputInt, succ := inputInt.SetString(string(tx), 10)
	if inputInt == nil || !succ {
		return types.ResponseCheckTx{Code: 1, Log: "Invalid input number"}
	}
	return types.ResponseCheckTx{Code: 0}
}

func (app *NumberAdderApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	inputInt := big.NewInt(0)
	inputInt, succ := inputInt.SetString(string(tx), 10)
	if inputInt == nil || !succ {
		return types.ResponseDeliverTx{Code: 1, Log: "Invalid input number"}
	}
	newState := app.State()
	newState.Add(newState, inputInt)
	app.SetState(newState)
	return types.ResponseDeliverTx{Code: 0}
}

func (app *NumberAdderApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	fmt.Printf("EndBlock, Height = %v\n", req.Height)
	return types.ResponseEndBlock{}
}

func (app *NumberAdderApplication) Commit() (resp types.ResponseCommit) {
	app.SetHeight(app.Height() + 1)
	appHash := app.AppHash()
	fmt.Printf("Commit, AppHash = %v\n", common.Bytes2Hex(appHash))
	return types.ResponseCommit{Data: appHash}
}

func (app *NumberAdderApplication) InitChain(params types.RequestInitChain) types.ResponseInitChain {
	fmt.Println("InitChain")
	app.SetState(big.NewInt(0))
	app.SetHeight(0)
	return types.ResponseInitChain{}
}

func (app *NumberAdderApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	fmt.Printf("Query, path = %s\n", reqQuery.Path)
	switch reqQuery.Path {
	case "hash":
		return types.ResponseQuery{Value: []byte(fmt.Sprintf("%v", app.AppHash()))}
	case "state":
		return types.ResponseQuery{Value: []byte(app.State().String())}
	default:
		return types.ResponseQuery{
			Log: fmt.Sprintf("Invalid query path: %s", reqQuery.Path),
		}
	}
}

func (app *NumberAdderApplication) Info(req types.RequestInfo) types.ResponseInfo {
	fmt.Println("Info")
	appHash := app.AppHash()
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"hash\":\"%v\"}", app.AppHash()),
		Version:          "1",
		LastBlockHeight:  int64(app.Height()),
		LastBlockAppHash: appHash,
	}
}

func NewNumberAdderApplication(db db.DB) *NumberAdderApplication {
	app := &NumberAdderApplication{db: db}
	fmt.Println("NewNumberAdderApplication")
	fmt.Printf("height: %v\n", app.Height())
	return app
}

func main() {
	db, err := db.NewGoLevelDB("num-adder", "/tmp")
	if err != nil {
		panic("Cannot open LevelDB file")
	}
	defer db.Close()

	app := NewNumberAdderApplication(db)
	svr, err := server.NewServer("tcp://0.0.0.0:26658", "socket", app)
	if err != nil {
		fmt.Println("Error when initializing server: ", err)
	}
	err = svr.Start()
	if err != nil {
		fmt.Println("Error when starting server: ", err)
	}
	cmn.TrapSignal(func() {
		svr.Stop()
	})
}
