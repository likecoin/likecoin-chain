#!/bin/bash

set -e

pushd "$(dirname "$0")" > /dev/null

VERSION="fotan-1"
COMMIT=$(git rev-parse HEAD)

echo "Building image for $VERSION using commit $COMMIT"
docker build --build-arg "VERSION=$VERSION" --build-arg "COMMIT=$COMMIT" -t "likecoin/likecoin-chain:$VERSION" .

popd > /dev/null
