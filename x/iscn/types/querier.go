package types

const (
	QueryRecord = "records"
	QueryParams = "params"
)

type QueryRecordParams struct {
	Id []byte // TODO: string?
}
