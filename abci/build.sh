#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

docker build -t likechain/abci .

cd "${PWD}"
