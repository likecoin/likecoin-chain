
# LikeChain

## Development Setup
```sh
# Build the docker images, run it for the first time or you have dependency updates
./build.sh

export NODE_COUNT=3
# Initialize nodes
./init.sh

# Start development
docker-compose up
```

### Scale nodes
```sh
# Change number of nodes
export NODE_COUNT=4
./init.sh
docker-compose up
```

### Reset nodes
```sh
./reset.sh
./init.sh

# Start development with fresh containers
docker-compose up --force-recreate
```

## Update dependancies
We use `dep` for package manager
```sh
# Make sure you have build the likechain/abci image
docker run --rm -v `pwd`/abci:/go/src/github.com/likecoin/likechain/ likechain/abci dep ensure -add github.com/ethereum/go-ethereum

# Or if you have installed `dep` locally
cd abci
dep ensure -add github.com/ethereum/go-ethereum

# Build image
./build.sh
```

## Usage
```sh
# Query the state
curl 'localhost:26657/abci_query?path="state"'

# Add number to the state
curl 'http://localhost:26657/broadcast_tx_commit?tx="100"'

# Check peers connection
curl 'http://localhost:26657/net_info'
```
