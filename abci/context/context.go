package context

import (
	logger "github.com/likecoin/likechain/abci/log"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

var log = logger.L

// TODO: Configurable
var cacheSize = 1048576

// ApplicationContext stores context of application
type ApplicationContext struct {
	state *MutableState
}

// New creates an ApplicationContext using GoLevelDB
func New(dbPath string) *ApplicationContext {
	stateDb := newDb("state", dbPath)
	stateTree := newTree(stateDb)
	if _, err := stateTree.Load(); err != nil {
		log.WithError(err).Panic("Unable to load state tree from disk")
	}

	withdrawDb := newDb("withdraw", dbPath)
	withdrawTree := newTree(withdrawDb)
	if _, err := withdrawTree.Load(); err != nil {
		log.WithError(err).Panic("Unable to load withdraw tree from disk")
	}

	return &ApplicationContext{
		state: &MutableState{
			appDb:      newDb("app", dbPath),
			stateDb:    stateDb,
			withdrawDb: withdrawDb,

			stateTree:    stateTree,
			withdrawTree: withdrawTree,
		},
	}
}

func newDb(name, path string) *db.GoLevelDB {
	db, err := db.NewGoLevelDB(name, path)
	if err != nil {
		log.WithError(err).Panic("Unable to create GoLevelDB")
	}
	return db
}

func newTree(db db.DB) *iavl.MutableTree {
	return iavl.NewMutableTree(db, cacheSize)
}

// GetImmutableState returns an immutable context
func (ctx *ApplicationContext) GetImmutableState() *ImmutableState {
	appDb := ctx.state.appDb
	if ctx.state.GetHeight() == 0 {
		return &ImmutableState{
			appDb:        appDb,
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
		appDb:        appDb,
		stateTree:    stateTree,
		withdrawTree: withdrawTree,
	}
}

// GetMutableState returns a mutable context
func (ctx *ApplicationContext) GetMutableState() *MutableState {
	return ctx.state
}

// CloseDb closes dbs
func (ctx *ApplicationContext) CloseDb() {
	ctx.state.appDb.Close()
	ctx.state.stateDb.Close()
	ctx.state.withdrawDb.Close()
}
