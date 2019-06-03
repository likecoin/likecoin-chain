#!/bin/bash

set -e

pushd "$(dirname $0)"
docker build -t likechain/gaia .
popd
