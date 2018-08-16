package context

type LikeContextMock struct {
	*LikeContext
}

func NewMock() *LikeContextMock {
	return &LikeContextMock{
		LikeContext: NewWithMemDB(),
	}
}

func (ctx *LikeContext) Reset() {
	ctx.MutableStateTree().Rollback()
	ctx.MutableWithdrawTree().Rollback()
}
