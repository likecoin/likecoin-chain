package context

import (
	"encoding/binary"
	"fmt"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type ImmutableContext interface {
	StateTree() *iavl.ImmutableTree
	WithdrawTree() *iavl.ImmutableTree
	GetBlockHash() []byte
}

type MutableContext interface {
	ImmutableContext

	MutableStateTree() *iavl.MutableTree
	MutableWithdrawTree() *iavl.MutableTree
}

type LikeContext struct {
	stateTree           *iavl.MutableTree
	withdrawTree        *iavl.MutableTree
	stateTreeVersion    int64
	withdrawTreeVersion int64
	blockHash           []byte
}

var cacheSize = 1048576

func newTree(dbPath, dir string) (tree *iavl.MutableTree, version int64) {
	db, err := db.NewGoLevelDB(dbPath, "like-state")
	if err != nil {
		panic(err)
	}
	tree = iavl.NewMutableTree(db, cacheSize)
	version = 0
	versionBytes := db.Get([]byte("META_VERSION"))
	if versionBytes != nil && len(versionBytes) == 8 {
		version = int64(binary.BigEndian.Uint64(versionBytes))
	}
	return tree, version
}

func NewWithMemDB() *LikeContext {
	return &LikeContext{
		stateTree:           iavl.NewMutableTree(db.NewMemDB(), 0),
		withdrawTree:        iavl.NewMutableTree(db.NewMemDB(), 0),
		stateTreeVersion:    0,
		withdrawTreeVersion: 0,
		blockHash:           []byte{1, 3, 3, 7},
	}
}

func New(dbPath string) *LikeContext {
	stateTree, stateTreeVersion := newTree(dbPath, "state")
	withdrawTree, withdrawTreeVersion := newTree(dbPath, "withdraw")
	return &LikeContext{
		stateTree:           stateTree,
		withdrawTree:        withdrawTree,
		stateTreeVersion:    stateTreeVersion,
		withdrawTreeVersion: withdrawTreeVersion,
		blockHash:           nil,
	}
}

func (ctx *LikeContext) StateTree() *iavl.ImmutableTree {
	tree, err := ctx.stateTree.GetImmutable(ctx.stateTreeVersion)
	if err != nil {
		panic(fmt.Sprintf("Cannot get state tree with version %d", ctx.stateTreeVersion))
	}
	return tree
}

func (ctx *LikeContext) WithdrawTree() *iavl.ImmutableTree {
	tree, err := ctx.withdrawTree.GetImmutable(ctx.withdrawTreeVersion)
	if err != nil {
		panic(fmt.Sprintf("Cannot get withdraw tree with version %d", ctx.withdrawTreeVersion))
	}
	return tree
}

func (ctx *LikeContext) MutableStateTree() *iavl.MutableTree {
	return ctx.stateTree
}

func (ctx *LikeContext) MutableWithdrawTree() *iavl.MutableTree {
	return ctx.withdrawTree
}

func (ctx *LikeContext) GetBlockHash() []byte {
	return ctx.blockHash
}

func (ctx *LikeContext) SetBlockHash(blockHash []byte) {
	ctx.blockHash = blockHash
}

func (ctx *LikeContext) Save() {
	// TODO
}
