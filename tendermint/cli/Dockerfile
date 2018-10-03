FROM likechain/golang-godep AS builder

ARG SRC=/go/src/github.com/likecoin/likechain/tendermint/cli

RUN mkdir -p ${SRC}
WORKDIR ${SRC}

# Copy the Gopkg.toml and Gopkg.lock to WORKDIR
COPY Gopkg.lock Gopkg.toml main.go ./

# Install the dependencies without checking for go code
RUN dep ensure -v -vendor-only

# Build the excutable
RUN go build -a -o /bin/likechain/cli main.go

FROM alpine:latest
WORKDIR /bin/likechain/
COPY --from=builder /bin/likechain/cli .

ENTRYPOINT ["./cli"]
