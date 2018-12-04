#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

./docker/golang/build.sh
./tendermint/cli/build.sh
./tendermint/build.sh

cd "${PWD}"
