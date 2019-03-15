FROM likechain/golang AS builder

COPY ./tendermint/cli ./tendermint/cli

# Build the excutable
RUN go build -a -o /bin/cli ./tendermint/cli/main.go

FROM alpine:latest
WORKDIR /likechain
COPY --from=builder /bin/cli /bin/cli

ENTRYPOINT ["/bin/cli"]
