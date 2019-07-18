#!/bin/bash

set -e

pushd "$(dirname "$0")/../.." > /dev/null
docker build -f docker/app/Dockerfile -t likechain/likechain .
popd > /dev/null
