FROM likechain/base as builder

ARG NAME="likecoin-chain"
ARG APP_NAME="liked"
ARG BUILD_TAGS="netgo ledger"
ARG VERSION="fotan-1"
ARG COMMIT=""

WORKDIR /likechain

COPY . .

RUN go build \
    -ldflags "\
    -X \"github.com/cosmos/cosmos-sdk/version.Name=${NAME}\" \
    -X \"github.com/cosmos/cosmos-sdk/version.AppName=${APP_NAME}\" \
    -X \"github.com/cosmos/cosmos-sdk/version.Version=${VERSION}\" \
    -X \"github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT}\" \
    -X \"github.com/cosmos/cosmos-sdk/version.BuildTags=${BUILD_TAGS}\" \
    " \
    -tags "${BUILD_TAGS}" \
    -o /go/bin/liked cmd/liked/main.go

FROM alpine:latest

RUN apk add ca-certificates
WORKDIR /bin
COPY --from=builder /go/bin/liked .
