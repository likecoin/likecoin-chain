#!/bin/bash

pushd "$(dirname $0)"

cp ../../go.mod ../../go.sum ./
docker build -t likechain/golang .
rm go.mod go.sum

popd
