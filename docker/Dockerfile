FROM golang:1.16-alpine AS base

RUN apk update && apk add --no-cache build-base git bash curl linux-headers ca-certificates
WORKDIR /
RUN mkdir -p ./likechain
COPY ./go.mod ./likechain/go.mod
COPY ./go.sum ./likechain/go.sum
WORKDIR /likechain
RUN go mod download

FROM base AS builder

WORKDIR /likechain
COPY . .
RUN go build -o /go/bin/liked cmd/liked/main.go
RUN go build -o /go/bin/likecli cmd/likecli/main.go

FROM alpine:latest

ARG UID=1000
ARG GID=1000

RUN apk add ca-certificates
RUN addgroup -g $GID likechain && \
    adduser --uid $UID --disabled-password --gecos "" --ingroup likechain --home /likechain likechain
USER $UID
WORKDIR /bin
RUN mkdir -p /likechain/.liked/data
RUN mkdir -p /likechain/.liked/config
RUN mkdir -p /likechain/.likecli
COPY --from=builder /go/bin/liked .
COPY --from=builder /go/bin/likecli .
