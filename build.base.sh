#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

./docker/golang-godep/build.sh
./tendermint/cli/build.sh
./tendermint/build.sh

cd "${PWD}"
