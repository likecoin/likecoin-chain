
# LikeChain

[![CircleCI](https://circleci.com/gh/likecoin/likechain.svg?style=svg)](https://circleci.com/gh/likecoin/likechain)

## Development Setup
```sh
# Build the docker images, run it for the first time or you have dependency updates
./build.sh

export NODE_COUNT=3

# Setup init script
cd tendermint/cli
dep ensure --vendor-only

# Initialize nodes
cd ../..
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

## Manage dependancies
We use `dep` for package manager
```sh
# After import package(s) in code / add constrain(s) in Gopkg.toml
cd abci
dep ensure

# Build image
./build.sh
```

## Usage
```sh
# Query the block info
curl 'http://localhost:3000/v1/block'

# Check peers connection
curl 'http://localhost:26657/net_info'
```

## Run in Production
```sh
# Build in production
./build.production.sh

# Using Docker Compose
docker-compose -f docker-compose.yml -f docker-compose.production.yml up
```
