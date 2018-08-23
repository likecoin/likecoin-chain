package context

import (
	"github.com/sirupsen/logrus"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// MockApplicationContext is a struct mocking ApplicationContext for testing
type MockApplicationContext struct {
	*ApplicationContext
}

// NewMock creates an MockApplicationContext using MemDB
func NewMock() *MockApplicationContext {
	// Hide all logs below error level for cleaner testing display
	log.SetLevel(logrus.ErrorLevel)

	return &MockApplicationContext{
		ApplicationContext: &ApplicationContext{
			state: &MutableState{
				stateTree:    iavl.NewMutableTree(db.NewMemDB(), 0),
				withdrawTree: iavl.NewMutableTree(db.NewMemDB(), 0),
			},
		},
	}
}

// Reset resets state tree and withdraw tree to last saved version
func (appCtx *MockApplicationContext) Reset() {
	appCtx.GetMutableState().MutableStateTree().Rollback()
	appCtx.GetMutableState().MutableWithdrawTree().Rollback()
}
