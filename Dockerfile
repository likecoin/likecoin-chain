FROM likecoin/rbuilder:cf0d1a9f3731e30540bbfa36a36d13e4dcccf5eb as builder

COPY . /sources
ARG VERSION
ARG COMMIT
WORKDIR /sources

ENV TARGET_PLATFORMS='linux/amd64'
ENV APP='liked'
ENV LEDGER_ENABLED=true
RUN /bin/bash -c /sources/.build.sh

FROM debian:buster

WORKDIR /usr/bin
RUN apt-get update && apt-get install -y curl
COPY --from=builder /sources/artifacts/liked-*-linux-amd64 /usr/bin/liked
