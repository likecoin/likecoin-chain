#!/bin/bash

pushd "$(dirname $0)"

docker build -t likechain/key-cli -f ./Dockerfile ../..

popd
