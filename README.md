
# LikeChain

[![CircleCI](https://circleci.com/gh/likecoin/likechain.svg?style=svg)](https://circleci.com/gh/likecoin/likechain) [![codecov](https://codecov.io/gh/likecoin/likechain/branch/master/graph/badge.svg)](https://codecov.io/gh/likecoin/likechain)

## Development Setup
```sh
# Build all
./build.sh

# OR build the following separately
# - Tendermint
# - cli (for key generation)
# - base image (for faster building speed)
./build.base.sh

# Build ABCI
./abci/build.sh

export NODE_COUNT=4

# Initialize nodes
./init.sh

# Start development
docker-compose up
```

### Scale nodes
```sh
# Change number of nodes
export NODE_COUNT=5
./init.sh
docker-compose up --force-recreate
```

### Reset nodes
```sh
./reset.sh

# Start development with fresh containers
docker-compose up --force-recreate
```

## Manage dependancies
We use `dep` for package manager, for example:
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

## Configuration

### Using Config File
ABCI server and API server shares a common configuration. You may copy the default config file `abci/config.example.toml` and rename it to `abci/config.toml`

### Using Environment Variable
You may also override the config by setting environment variable, for example, if you wish to set `env=production`, you can set
```sh
LIKECHAIN_ENV=production
```

## Run in Production
```sh
# Build in production
./build.production.sh

# OR if you have run the following before
./build.base.sh
./abci/build.sh
# You can directly run this
./abci/build.production.sh

# Using Docker Compose
docker-compose -f docker-compose.yml -f docker-compose.production.yml up
```
