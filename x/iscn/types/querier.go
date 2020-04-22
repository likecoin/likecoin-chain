package types

const (
	QueryIscnRecord = "records"
	QueryAuthor     = "author"
	QueryParams     = "params"
)

// TODO: query author and terms by name?

type QueryAuthorParams struct {
	Cid []byte
}

type QueryRightTermParams struct {
	Cid []byte
}

type QueryRecordParams struct {
	Id []byte // TODO: string?
}
