#!/bin/bash

pushd "$(dirname $0)"

docker build -f Dockerfile.production -t likechain/abci-production .

popd
