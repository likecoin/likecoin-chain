FROM golang:1.12-alpine as builder

WORKDIR /
COPY . ./likechain
RUN apk update && apk add --no-cache build-base git bash curl linux-headers
WORKDIR /likechain
RUN go mod download
RUN go build -o /go/bin/liked cmd/liked/main.go
RUN go build -o /go/bin/likecli cmd/likecli/main.go

FROM alpine:latest
WORKDIR /bin

COPY --from=builder /go/bin/liked .
COPY --from=builder /go/bin/likecli .
