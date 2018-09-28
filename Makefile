PACKAGES=$(shell go list ./...)

default: build

get_vendor_deps:
	@go get -v "github.com/ethereum/go-ethereum/crypto" "github.com/ethereum/go-ethereum/common"
	@cd abci;dep ensure -v -vendor-only
	@cd tendermint/cli;dep ensure -v -vendor-only

test:
	@GOCACHE=off go test -p 1 $(PACKAGES)
