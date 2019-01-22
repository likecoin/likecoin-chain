#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

./build.sh
./abci/build.production.sh
./services/build.production.sh

cd "${PWD}"
