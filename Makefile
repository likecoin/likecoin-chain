default: build

get_vendor_deps:
	@cd abci;dep ensure -v -vendor-only
	@go get -v "github.com/ethereum/go-ethereum/crypto/secp256k1"
	@cp -r /go/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 abci/vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/

test: test-app

test-app:
	@go test -v github.com/likecoin/likechain/abci/app
