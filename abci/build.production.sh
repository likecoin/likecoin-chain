#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

docker build -f Dockerfile.production -t likechain/abci-production .

cd "${PWD}"
