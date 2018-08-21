package context

import (
	"fmt"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// IImmutableState is an interface for accessing mutable context
type IImmutableState interface {
	ImmutableStateTree() *iavl.ImmutableTree
	ImmutableWithdrawTree() *iavl.ImmutableTree
	GetBlockHash() []byte
}

// IMutableState is an interface for accessing immutable context
type IMutableState interface {
	IImmutableState
	MutableStateTree() *iavl.MutableTree
	MutableWithdrawTree() *iavl.MutableTree
}

// TODO: Configurable
var cacheSize = 1048576

// ApplicationContext stores context of application
type ApplicationContext struct {
	state *MutableState
}

// New creates an ApplicationContext using GoLevelDB
func New(dbPath string) *ApplicationContext {
	return &ApplicationContext{
		state: &MutableState{
			stateTree:    newTree(dbPath, "state"),
			withdrawTree: newTree(dbPath, "withdraw"),
		},
	}
}

func newTree(dbPath, dir string) *iavl.MutableTree {
	db, err := db.NewGoLevelDB(dbPath, dir)
	if err != nil {
		panic(fmt.Sprintf("Unable to create GoLevelDB: %v", err))
	}
	return iavl.NewMutableTree(db, cacheSize)
}

// GetImmutableState returns an immutable context
func (ctx *ApplicationContext) GetImmutableState() *ImmutableState {
	stateTree, err := ctx.state.stateTree.GetImmutable(ctx.state.stateTree.Version64())
	if err != nil {
		panic(fmt.Sprintf("Unable to get versioned state tree: %v", err))
	}
	withdrawTree, err := ctx.state.withdrawTree.GetImmutable(ctx.state.withdrawTree.Version64())
	if err != nil {
		panic(fmt.Sprintf("Unable to get versioned withdraw tree: %v", err))
	}
	return &ImmutableState{
		stateTree:    stateTree,
		withdrawTree: withdrawTree,
	}
}

// GetMutableState returns a mutable context
func (ctx *ApplicationContext) GetMutableState() *MutableState {
	return ctx.state
}
