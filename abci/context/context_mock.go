package context

import (
	"math/big"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// MockApplicationContext is a struct mocking ApplicationContext for testing
type MockApplicationContext struct {
	*ApplicationContext
}

// NewMock creates an MockApplicationContext using MemDB
func NewMock() *MockApplicationContext {
	return &MockApplicationContext{
		ApplicationContext: &ApplicationContext{
			state: &MutableState{
				appDb:        db.NewMemDB(),
				stateTree:    iavl.NewMutableTree(db.NewMemDB(), 0),
				withdrawTree: iavl.NewMutableTree(db.NewMemDB(), 0),
			},
		},
	}
}

// Reset resets state tree and withdraw tree to last saved version
func (appCtx *MockApplicationContext) Reset() {
	itr := appCtx.state.appDb.Iterator(nil, nil)
	for ; itr.Valid(); itr.Next() {
		appCtx.state.appDb.Delete(itr.Key())
	}
	appCtx.state.stateTree.Rollback()
	appCtx.state.withdrawTree.Rollback()
}

// SetInitialBalance set the initial balance for new account
func (appCtx *MockApplicationContext) SetInitialBalance(i *big.Int) {
	appCtx.state.initialBalance = i
}
