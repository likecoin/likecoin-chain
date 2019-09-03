#!/bin/bash

set -e

SYSTEM="$(uname)"
if [ $SYSTEM == "Darwin" ]; then
    BUILD_UID=1000
    BUILD_GID=1000
else
    BUILD_UID=$(id -u $USER)
    BUILD_GID=$(id -g $USER)
fi

pushd "$(dirname "$0")/../.." > /dev/null
docker build -f docker/app/Dockerfile -t likechain/likechain \
    --build-arg UID=$BUILD_UID --build-arg GID=$BUILD_GID .
popd > /dev/null
