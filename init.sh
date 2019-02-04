#!/bin/bash

set -e

if [[ -z $1 ]]; then
    echo "Usage: $0 [NODE_COUNT]"
    exit 1
fi

NODE_COUNT="$1"
CLI_IMAGE="likechain/key-cli"
WD=$(cd $(dirname "$0"); pwd)

rm -rf $WD/tendermint/nodes

docker run --rm \
    -v $WD:/likechain \
    $CLI_IMAGE init --profile_dir /likechain/tendermint/nodes --node_count $NODE_COUNT
