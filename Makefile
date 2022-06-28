#!/usr/bin/make -f

NAME := likecoin-chain
APP := liked
VERSION := $(shell git describe --tags)
COMMIT := $(shell git rev-parse HEAD)
LEDGER_ENABLED ?= true
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
LIKE_HOME := '$(HOME)/.liked'
IMAGE_TAG = likecoin/likecoin-chain:$(VERSION)
RBUILDER_IMAGE_TAG = likecoin/rbuilder:go1.18
BUILDDIR ?= $(CURDIR)/build
GOPATH ?= '$(HOME)/go'
GOLANG_VERSION ?= 1.18.3
GOLANG_CROSS_VERSION := v$(GOLANG_VERSION)

###############################################################################
###                            Development                                  ###
###############################################################################

all: install test

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

$(BUILDDIR)/liked: build

vendor: go.sum
	@echo "--> Download go modules to work directory"
	go mod vendor

download: go.sum
	@echo "--> Download go modules to local cache"
	go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

install: go.sum $(BUILDDIR)/
	go install -mod=readonly $(BUILD_FLAGS) ./...

test:
	go test -v ./...

clean:
	rm -rf $(BUILDDIR)/ artifacts/

init: $(BUILDDIR)/liked
	$(BUILDDIR)/liked --home $(LIKE_HOME) init ${MONIKER} --chain-id "${CHAIN_ID}"

lint:
	golangci-lint run --disable-all -E errcheck --timeout 10m
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/cosmos/cosmos-sdk

.PHONY: all vendor download install test clean lint format

###############################################################################
###                               Build                                     ###
###############################################################################

build-reproducible: go.sum
	$(DOCKER) rm latest-build || true
	$(DOCKER) run --volume=$(CURDIR):/sources:ro \
		--env TARGET_PLATFORMS='linux/amd64 darwin/amd64 darwin/arm64 linux/arm64 windows/amd64' \
		--env APP=$(APP) \
		--env VERSION=$(VERSION) \
		--env COMMIT=$(COMMIT) \
		--env LEDGER_ENABLED=$(LEDGER_ENABLED) \
		--name latest-build $(RBUILDER_IMAGE_TAG)
	$(DOCKER) cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/

docker-build: go.sum
	@echo "Building image for $(VERSION) using commit $(COMMIT)"
	$(DOCKER) build \
		--build-arg LIKED_VERSION=$(VERSION) \
		--build-arg LIKED_COMMIT=$(COMMIT) \
		--tag $(IMAGE_TAG) \
		.

docker-push:
	@echo "Pushing image $(IMAGE_TAG) to registry"
	$(DOCKER) push $(IMAGE_TAG)

build: go.sum $(BUILDDIR)/
	go build -mod=readonly \
    -ldflags "\
			-w -s \
			-X \"github.com/cosmos/cosmos-sdk/version.Name=$(NAME)\" \
			-X \"github.com/cosmos/cosmos-sdk/version.AppName=$(APP)\" \
			-X \"github.com/cosmos/cosmos-sdk/version.BuildTags=netgo,ledger\" \
			-X \"github.com/cosmos/cosmos-sdk/version.Version=${VERSION}\" \
			-X \"github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT}\" \
		" \
	-tags "netgo,ledger" \
	-o $(BUILDDIR)/ ./...

.PHONY: build-reproducible docker-login docker-build docker-push build

###############################################################################
###                               Release                                   ###
###############################################################################


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
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist --skip-validate

.PHONY: release

###############################################################################
###                              Protobuf                                   ###
###############################################################################

proto-all: proto-update-deps proto-format proto-lint gen-proto

gen-proto: x/
	ignite generate proto-go

proto-format:
	@$(DOCKER_BUF) format -w

proto-format-check:
	@$(DOCKER_BUF) format --diff --exit-code

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-update-deps:
	@$(DOCKER_BUF) mod update proto

.PHONY: proto-all gen-proto proto-format proto-format-check proto-lint proto-update-deps
