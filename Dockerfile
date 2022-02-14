FROM likecoin/rbuilder:cf0d1a9f3731e30540bbfa36a36d13e4dcccf5eb as builder

ARG LIKED_VERSION=unknown
ARG LIKED_COMMIT=unknown

WORKDIR /cosmovisor
RUN wget https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv1.1.0/cosmovisor-v1.1.0-linux-amd64.tar.gz
RUN tar -xzvf cosmovisor-v1.1.0-linux-amd64.tar.gz

COPY . /sources
WORKDIR /sources

ENV TARGET_PLATFORMS='linux/amd64'
ENV APP='liked'
ENV LEDGER_ENABLED=true
ENV VERSION=$LIKED_VERSION
ENV COMMIT=$LIKED_COMMIT
RUN rm -rf /sources/artifacts
RUN /bin/bash -c /sources/.build.sh

FROM debian:buster

WORKDIR /usr/bin
ENV DAEMON_NAME liked
ENV DAEMON_HOME /likechain/.liked
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES true
ENV DAEMON_RESTART_AFTER_UPGRADE true
RUN mkdir -p /likechain/.liked/cosmovisor/genesis/bin /likechain/.liked/data /likechain/.liked/config
RUN apt-get update && apt-get install -y curl
COPY --from=builder /cosmovisor/cosmovisor /usr/bin/cosmovisor
COPY --from=builder ./home/builder/artifacts/liked-*-linux-amd64 /usr/bin/liked
RUN cp /usr/bin/liked /likechain/.liked/cosmovisor/genesis/bin/liked
