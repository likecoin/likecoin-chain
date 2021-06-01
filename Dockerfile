FROM golang:1.16-alpine as base

RUN apk update && apk add --no-cache build-base git bash curl linux-headers ca-certificates
WORKDIR /
RUN mkdir -p ./likechain
COPY ./go.mod ./likechain/go.mod
COPY ./go.sum ./likechain/go.sum
WORKDIR /likechain
RUN go mod download


FROM base as builder

WORKDIR /likechain
COPY . .
ARG VERSION="fotan-1"
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


FROM alpine:latest

RUN apk add ca-certificates
WORKDIR /bin
COPY --from=builder /go/bin/liked .
