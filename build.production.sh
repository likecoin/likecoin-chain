#!/bin/bash

pushd "$(dirname $0)"

./build.sh
./abci/build.production.sh
./services/build.production.sh

popd
