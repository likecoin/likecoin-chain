FROM golang:1.12-alpine

RUN apk update && apk add --no-cache build-base git bash curl linux-headers ca-certificates
WORKDIR /
RUN mkdir -p ./likechain
COPY ./go.mod ./likechain/go.mod
COPY ./go.sum ./likechain/go.sum
WORKDIR /likechain
RUN go mod download
