#!/bin/bash

set -e

LIKE_HOME="$(dirname "$0")/.."
pushd "$LIKE_HOME" > /dev/null
LIKE_HOME=$(pwd)
popd > /dev/null

MONIKER="$1"
GENESIS_PATH="$2"
SEED_NODES="$3"

if [ -z $MONIKER ] || [ -z $GENESIS_PATH ]; then
    echo "Usage: $0 YOUR_NODE_NAME path_to_genesis.json SEED_NODES"
    echo "Example: $0 chung ~/Downloads/genesis-likechain-cosmos-testnet-1.json '7c93876c5ffce59b5bc07a4b4b7891dd0bfe4cea@35.226.174.222:26656'"
    exit 1
fi

CHAIN_ID=`grep chain_id "$GENESIS_PATH" | sed 's/ *"chain_id": *"\(.*\)"/\1/g' | sed 's/,$//g'`

if [ ! -f "$LIKE_HOME/.liked/config/genesis.json" ]; then
    mkdir -p "$LIKE_HOME/.liked"
    docker run --rm --volume "$LIKE_HOME/.liked:/likechain/.liked" likechain/likechain:sheungwan liked --home /likechain/.liked init --chain-id "$CHAIN_ID" "$MONIKER"
    cp "$GENESIS_PATH" "$LIKE_HOME/.liked/config/genesis.json"
else
    echo "Warning: .liked already exists, not re-initializing. You may need to replace .liked/config/genesis.json manually."
fi

if [ ! -f "$LIKE_HOME/docker-compose.yml" ]; then
    sed "s/__SEED_NODES__/$SEED_NODES/g" "$LIKE_HOME/docker-compose.template.yml" > "$LIKE_HOME/docker-compose.yml"
else
    echo "Warning: docker-compose.yml already exists, not modifying. You may need to modify docker-compose.yml manually to add the seed nodes parameter."
fi

mkdir -p "$LIKE_HOME/.likecli"

docker run --rm -it --volume "$LIKE_HOME/.likecli:/likechain/.likecli" likechain/likechain:sheungwan likecli --home /likechain/.likecli keys add validator
ADDRESS=$(docker run --rm -it --volume "$LIKE_HOME/.likecli:/likechain/.likecli" likechain/likechain:sheungwan likecli --home /likechain/.likecli keys show validator -a)
echo ""
echo "--------------------------------------------------------------------------------"
echo "Key initialized, your address is $ADDRESS"
