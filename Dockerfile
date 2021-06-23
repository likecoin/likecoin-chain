FROM golang:1.16-buster as base

WORKDIR /
RUN mkdir -p ./likechain
COPY ./go.mod ./likechain/go.mod
COPY ./go.sum ./likechain/go.sum
WORKDIR /likechain
RUN go mod download


FROM base as builder

WORKDIR /likechain
COPY . .
ARG VERSION="sheungwan-2"
ARG COMMIT=""
RUN go build \
    -ldflags "\
    -X \"github.com/cosmos/cosmos-sdk/version.Name=likecoin-chain\" \
    -X \"github.com/cosmos/cosmos-sdk/version.AppName=liked\" \
    -X \"github.com/cosmos/cosmos-sdk/version.BuildTags=netgo ledger\" \
    -X \"github.com/cosmos/cosmos-sdk/version.Version=${VERSION}\" \
    -X \"github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT}\" \
    " \
    -tags "netgo ledger" \
    -o /go/bin/liked cmd/liked/main.go
RUN go build \
    -ldflags "\
    -X \"github.com/cosmos/cosmos-sdk/version.Name=likecoin-chain\" \
    -X \"github.com/cosmos/cosmos-sdk/version.AppName=liked\" \
    -X \"github.com/cosmos/cosmos-sdk/version.BuildTags=netgo ledger\" \
    -X \"github.com/cosmos/cosmos-sdk/version.Version=${VERSION}\" \
    -X \"github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT}\" \
    " \
    -tags "netgo ledger" \
    -o /go/bin/likecli cmd/likecli/main.go


FROM debian:buster

WORKDIR /usr/bin
RUN apt-get update && apt-get install -y curl
COPY --from=builder /go/bin/liked .
COPY --from=builder /go/bin/likecli .
