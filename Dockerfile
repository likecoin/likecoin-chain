FROM golang:1.12-alpine as builder

WORKDIR /
RUN apk update && apk add --no-cache build-base git bash curl linux-headers
RUN git clone https://github.com/cosmos/cosmos-sdk
WORKDIR cosmos-sdk
RUN git checkout v0.34.6
RUN make install

FROM alpine:latest
WORKDIR /bin

COPY --from=builder /go/bin/gaiad .
COPY --from=builder /go/bin/gaiacli .
