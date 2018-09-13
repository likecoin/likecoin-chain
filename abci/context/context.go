package context

import (
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

var log = logger.L

// IImmutableState is an interface for accessing mutable context
type IImmutableState interface {
	ImmutableStateTree() *iavl.ImmutableTree
	ImmutableWithdrawTree() *iavl.ImmutableTree
	GetBlockHash() []byte
	GetHeight() int64
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
			stateTree:    newTree("state", dbPath),
			withdrawTree: newTree("withdraw", dbPath),
		},
	}
}

func newTree(name, path string) *iavl.MutableTree {
	db, err := db.NewGoLevelDB(name, path)
	if err != nil {
		log.WithError(err).Panic("Unable to create GoLevelDB")
	}
	return iavl.NewMutableTree(db, cacheSize)
}

// GetImmutableState returns an immutable context
func (ctx *ApplicationContext) GetImmutableState() *ImmutableState {
	if ctx.state.GetHeight() == 0 {
		return &ImmutableState{
			stateTree:    ctx.state.stateTree.ImmutableTree,
			withdrawTree: ctx.state.withdrawTree.ImmutableTree,
		}
	}

	stateTreeVersion := ctx.state.stateTree.Version64()
	stateTree, err := ctx.state.stateTree.GetImmutable(stateTreeVersion)
	if err != nil {
		log.
			WithError(err).
			WithField("version", stateTreeVersion).
			Panic("Unable to get versioned state tree")
	}

	withdrawTreeVersion := ctx.state.withdrawTree.Version64()
	withdrawTree, err := ctx.state.withdrawTree.GetImmutable(withdrawTreeVersion)
	if err != nil {
		log.
			WithError(err).
			WithField("version", withdrawTreeVersion).
			Panic("Unable to get versioned withdraw tree")
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
