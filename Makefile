#!/usr/bin/make -f

NAME := likecoin-chain
APP := liked
VERSION := $(shell git describe --tags)
COMMIT := $(shell git rev-parse HEAD)
LEDGER_ENABLED ?= true
DOCKER := $(shell which docker)
IMAGE_TAG = likecoin/likecoin-chain:$(VERSION)
BUILDDIR ?= $(CURDIR)/build
SWAGGER_OUT := swagger-gen
COSMOS_SDK_VERSION := $(shell grep "github.com/cosmos/cosmos-sdk" go.mod | head -n 1 | sed 's/.*github.com\/cosmos\/cosmos-sdk \(.*\)/\1/g')

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

go-mod-cache: go.sum
	echo "--> Download go modules to local cache"
	go mod download

go.sum: go.mod
	echo "--> Ensure dependencies have not been modified"
	go mod verify

gen-proto: x/iscn
	mkdir -p ${SWAGGER_OUT}
	protoc \
		-I "${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${COSMOS_SDK_VERSION}/proto" \
		-I "${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${COSMOS_SDK_VERSION}/third_party/proto" \
		--gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
		--grpc-gateway_out=logtostderr=true:. \
		--proto_path proto \
		./proto/iscn/* \
		--swagger_out ${SWAGGER_OUT} \
		--swagger_opt logtostderr=true --swagger_opt fqn_for_swagger_name=true --swagger_opt simple_operation_ids=true
	mv github.com/likecoin/likechain/x/iscn/types/* x/iscn/types/

build-reproducible: go.sum
	$(DOCKER) rm latest-build || true
	$(DOCKER) run --volume=$(CURDIR):/sources:ro \
        --env TARGET_PLATFORMS='linux/amd64 darwin/amd64 linux/arm64 windows/amd64' \
        --env APP=$(APP) \
        --env VERSION=$(VERSION) \
        --env COMMIT=$(COMMIT) \
        --env LEDGER_ENABLED=$(LEDGER_ENABLED) \
        --name latest-build likecoin/rbuilder:latest
	$(DOCKER) cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/

build-docker: go.sum
	echo "Building image for $(VERSION) using commit $(COMMIT)"
	$(DOCKER) build \
        --build-arg VERSION=$(VERSION) \
        --build-arg COMMIT=$(COMMIT) \
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

.PHONY: go-mod-cache gen-proto build-reproducible build-docker build install test clean
