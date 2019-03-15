#!/bin/bash

pushd "$(dirname $0)"

docker build -t likechain/tendermint .

popd
