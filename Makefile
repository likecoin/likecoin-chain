default: build

get_vendor_deps:
	@go get -v "github.com/ethereum/go-ethereum/crypto" "github.com/ethereum/go-ethereum/common"
	@dep ensure -v -vendor-only

test:
	@./test.sh
