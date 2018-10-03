FROM likechain/golang-godep

# Hack dep issue, ref: 
# - https://github.com/golang/dep/issues/1847
# - https://github.com/ethereum/go-ethereum/issues/2738
RUN go get -v "github.com/ethereum/go-ethereum/crypto" "github.com/ethereum/go-ethereum/common"

ARG SRC=/go/src/github.com/likecoin/likechain/abci

RUN mkdir -p ${SRC}
WORKDIR ${SRC}

# Copy the Gopkg.toml and Gopkg.lock to WORKDIR
COPY Gopkg.lock Gopkg.toml ./

# Install the dependencies without checking for go code
RUN dep ensure -v -vendor-only

CMD ["go", "run", "cmd/abci/main.go"]
