FROM likechain/base as builder

WORKDIR /likechain
COPY . .
RUN go build -o /go/bin/liked cmd/liked/main.go
RUN go build -o /go/bin/likecli cmd/likecli/main.go

FROM alpine:latest

WORKDIR /bin
COPY --from=builder /go/bin/liked .
COPY --from=builder /go/bin/likecli .
