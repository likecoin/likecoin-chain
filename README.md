
# LikeChain

[![CircleCI](https://circleci.com/gh/likecoin/likechain.svg?style=svg)](https://circleci.com/gh/likecoin/likechain) [![codecov](https://codecov.io/gh/likecoin/likechain/branch/master/graph/badge.svg)](https://codecov.io/gh/likecoin/likechain)

## Requirement
Linux or macOS, with Docker installed and working.

Tested OS: Ubuntu 16.04, Ububtu 18.04, macOS

Recommended Docker version: 18.09.2 or later

Recommeneded Docker Compose version: 1.23.2 or later

Recommended hardware: 2-core CPU, 1 GB RAM, 20 GB Harddisk

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

# Initialize 4 nodes
./init.sh 4

# Start development
docker-compose up
```

### Scale nodes
```sh
# Change number of nodes
./init.sh 5
docker-compose up --force-recreate
```

### Reset nodes
```sh
./reset.sh

# Start development with fresh containers
docker-compose up --force-recreate
```

## Manage dependencies
We use Go 1.11 modules for dependencies, so after installing Go 1.11, the general commands (`go run`, `go build`, `go test`, etc) should install the dependencies automatically.

Since module support is an experimental feature in Go 1.11, if you are placing the project in `GOPATH`, you need to `export GO111MODULE=on`.

If you are having network problem on packages with golang.org, you may use the alternative `go.mod.replace` to replace the original `go.mod`.

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
