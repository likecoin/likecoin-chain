FROM likechain/golang as builder

# Copy files to WORKDIR
COPY ./services ./services
COPY ./abci ./abci

# Build executables
RUN go build -a -o /bin/likechain/like_service services/main.go

FROM alpine:latest

WORKDIR /bin/likechain/

COPY --from=builder /bin/likechain/like_service .
RUN apk add ca-certificates

CMD ["./like_service"]
