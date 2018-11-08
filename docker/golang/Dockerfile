FROM golang:1.11.2-alpine3.8

ENV GO111MODULE on

# Install Git
RUN apk update && apk upgrade && apk add --no-cache build-base git

# Install dependencies
ARG SRC=/go/src/github.com/likecoin/likechain
RUN mkdir -p ${SRC}
WORKDIR ${SRC}
COPY go.mod go.sum ./
RUN go mod download
