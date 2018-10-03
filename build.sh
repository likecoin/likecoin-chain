#!/bin/bash

docker build -t likechain/key-cli tendermint/cli
docker build -t likechain/tendermint tendermint
docker build -t likechain/abci abci
