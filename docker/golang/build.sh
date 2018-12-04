#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

cp ../../go.mod ../../go.sum ./
docker build -t likechain/golang .
rm go.mod go.sum

cd "${PWD}"
