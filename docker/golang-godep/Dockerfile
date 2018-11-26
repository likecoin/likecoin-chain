FROM golang:1.10.3-alpine3.8

ENV DEP_VERSION 0.5.0

# Install Git
RUN apk update && apk upgrade && apk add --no-cache build-base git

# Install dep
ADD https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 /usr/local/bin/dep
RUN chmod +x /usr/local/bin/dep

# Install dependencies using dep
RUN go get -v "github.com/ethereum/go-ethereum/crypto" "github.com/ethereum/go-ethereum/common"
ARG SRC=/go/src/github.com/likecoin/likechain

RUN mkdir -p ${SRC}
WORKDIR ${SRC}

# Copy the Gopkg.toml and Gopkg.lock to WORKDIR
COPY Gopkg.lock Gopkg.toml ./

# Install the dependencies without checking for go code
RUN dep ensure -v -vendor-only
