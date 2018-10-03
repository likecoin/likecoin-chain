#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

./build.base.sh
./abci/build.sh

cd "${PWD}"
