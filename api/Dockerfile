FROM golang:1.10.3-alpine3.8

ENV DEP_VERSION 0.5.0

# Install Git
RUN apk update && apk upgrade && apk add --no-cache build-base git

# Install dep
ADD https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 /usr/local/bin/dep
RUN chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/likecoin/likechain/api
WORKDIR /go/src/github.com/likecoin/likechain/api

# Copy the Gopkg.toml and Gopkg.lock to WORKDIR
COPY Gopkg.lock Gopkg.toml ./

# Install the dependencies without checking for go code
RUN dep ensure -v -vendor-only

# Hack dep issue, ref: 
# - https://github.com/golang/dep/issues/1847
# - https://github.com/ethereum/go-ethereum/issues/2738
RUN go get -v "github.com/ethereum/go-ethereum/crypto/secp256k1"
RUN cp -r \
  /go/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
  vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/

EXPOSE 3000

CMD ["go", "run", "main.go"]