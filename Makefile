default: build

get_vendor_deps:
	@go mod download

test:
	@./test.sh
