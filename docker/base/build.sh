#!/bin/bash

set -e

pushd "$(dirname "$0")/../.." > /dev/null
docker build -f docker/base/Dockerfile -t likechain/base .
popd > /dev/null
