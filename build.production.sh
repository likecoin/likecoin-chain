#!/bin/bash

docker build -t likechain/tendermint tendermint
docker build -f abci/Dockerfile.production -t likechain/abci abci
