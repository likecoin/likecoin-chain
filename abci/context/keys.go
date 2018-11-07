package context

var (
	initKey = []byte("$init")

	appBlockHashKey     = []byte("$app_blockHash")
	appBlockTimeKey     = []byte("$app_blockTime")
	appHeightKey        = []byte("$app_height")
	appMetadataAtHeight = []byte("$app_hToMeta")
)
