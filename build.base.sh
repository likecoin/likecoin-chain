#!/bin/bash

pushd "$(dirname $0)"

./docker/golang/build.sh
./tendermint/cli/build.sh
./tendermint/build.sh

popd
