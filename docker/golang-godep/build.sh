#!/bin/bash

PWD=`pwd`
WD=`cd $(dirname "$0") && pwd -P`

cd "${WD}"

cp ../../Gopkg.lock ../../Gopkg.toml ./
docker build -t likechain/golang-godep .
rm Gopkg.lock Gopkg.toml

cd "${PWD}"
