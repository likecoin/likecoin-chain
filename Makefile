#!/usr/bin/make -f

NAME := likecoin-chain
APP := liked
VERSION := $(shell git describe --tags)
COMMIT := $(shell git rev-parse HEAD)
LEDGER_ENABLED ?= true
DOCKER := $(shell which docker)
IMAGE_TAG = likecoin/likecoin-chain:$(VERSION)
RBUILDER_IMAGE_TAG = cf0d1a9f3731e30540bbfa36a36d13e4dcccf5eb
BUILDDIR ?= $(CURDIR)/build
GOPATH ?= '$(HOME)/go'
GOLANG_VERSION        ?= 1.17.2
GOLANG_CROSS_VERSION  := v$(GOLANG_VERSION)

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(LIKE_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(NAME) \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=$(APP) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq (cleveldb,$(findstring cleveldb,$(LIKE_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (,$(findstring nostrip,$(LIKE_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(LIKE_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

all: install test

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

vendor: go.sum
	@echo "--> Download go modules to local cache"
	go mod download

go-mod-cache: vendor


go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

gen-proto: x/
	./gen_proto.sh

build-reproducible: go.sum
	$(DOCKER) rm latest-build || true
	$(DOCKER) run --volume=$(CURDIR):/sources:ro \
        --env TARGET_PLATFORMS='linux/amd64 darwin/amd64 linux/arm64 windows/amd64' \
        --env APP=$(APP) \
        --env VERSION=$(VERSION) \
        --env COMMIT=$(COMMIT) \
        --env LEDGER_ENABLED=$(LEDGER_ENABLED) \
        --name latest-build likecoin/rbuilder:$(RBUILDER_IMAGE_TAG)
	$(DOCKER) cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/

build-docker: go.sum
	@echo "Building image for $(VERSION) using commit $(COMMIT)"
	$(DOCKER) build \
        --build-arg LIKED_VERSION=$(VERSION) \
        --build-arg LIKED_COMMIT=$(COMMIT) \
        --tag $(IMAGE_TAG) \
		.

build: go.sum $(BUILDDIR)/
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/ ./...

install: go.sum $(BUILDDIR)/
	go install -mod=readonly $(BUILD_FLAGS) ./...

test:
	go test -v ./...

clean:
	rm -rf $(BUILDDIR)/ artifacts/

lint:
	golangci-lint run --disable-all -E errcheck --timeout 10m
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/cosmos/cosmos-sdk

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(NAME) \
		-w /go/src/$(NAME) \
		ghcr.io/troian/golang-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist --skip-validate

.PHONY: go-mod-cache gen-proto build-reproducible build-docker build install test clean lint format vendor release-dry-run release
