FROM golang:1.10.3-alpine3.8

ENV DATA_ROOT /tendermint
ENV TMHOME $DATA_ROOT

# Set user right away for determinism
RUN addgroup tmuser && \
    adduser -S -G tmuser tmuser

# Create directory for persistence and give our user ownership
RUN mkdir -p $DATA_ROOT && \
    chown -R tmuser:tmuser $DATA_ROOT

# Install Git
RUN apk update && apk upgrade && apk add --no-cache build-base git

# Download latest source code of Tendermint on master branch
RUN mkdir -p $GOPATH/src/github.com/tendermint && \
    cd $GOPATH/src/github.com/tendermint && \
    git clone https://github.com/tendermint/tendermint.git && \
    cd tendermint

WORKDIR $GOPATH/src/github.com/tendermint/tendermint

# Get Tools & Dependencies
RUN make get_tools && \
    make get_vendor_deps && \
    make install && \
    make build

RUN mkdir -p $DATA_ROOT

WORKDIR $DATA_ROOT

# rpc port
EXPOSE 26657

ENTRYPOINT ["tendermint"]
CMD ["node", "--consensus.create_empty_blocks=false"]