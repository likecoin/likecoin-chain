FROM likechain/golang AS builder

ARG SRC=/go/src/github.com/likecoin/likechain/tendermint/cli

RUN mkdir -p ${SRC}
WORKDIR ${SRC}

COPY main.go ./

# Build the excutable
RUN go build -a -o /bin/likechain/cli main.go

FROM alpine:latest
WORKDIR /bin/likechain/
COPY --from=builder /bin/likechain/cli .

ENTRYPOINT ["./cli"]
