#!/bin/bash

set -e

pushd "$(dirname "$0")/.." > /dev/null
docker build -t likechain/likechain .
popd > /dev/null
