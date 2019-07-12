#!/bin/bash

set -e

LIKE_HOME="$(dirname "$0")/.."
pushd "$LIKE_HOME" > /dev/null
LIKE_HOME=$(pwd)
popd > /dev/null

MONIKER="$1"
GENESIS_PATH="$2"
PEERS="$3"

if [ -z $MONIKER ] || [ -z $GENESIS_PATH ]; then
    echo "Usage: $0 YOUR_NODE_NAME path_to_genesis.json PEER"
    echo "Example: $0 chung ~/Downloads/genesis-likechain-cosmos-testnet-1.json '7c93876c5ffce59b5bc07a4b4b7891dd0bfe4cea@35.226.174.222:26656'"
    exit 1
fi

CHAIN_ID=`grep chain_id "$GENESIS_PATH" | sed 's/ *"chain_id": *"\(.*\)"/\1/g' | sed 's/,$//g'`

mkdir -p "$LIKE_HOME/.liked"
mkdir -p "$LIKE_HOME/.likecli"

docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain liked init --chain-id "$CHAIN_ID" "$MONIKER" > /dev/null 2>&1

# not using sed -i since different behaviour on Linux and Mac
sed "s/persistent_peers *=.*/persistent_peers=\"$PEERS\"/g" "$LIKE_HOME/.liked/config/config.toml" \
| sed "s/create_empty_blocks_interval *=.*/create_empty_blocks_interval=\"60s\"/g" \
| sed "s/^timeout_commit *=.*/timeout_commit=\"60s\"/g" \
> "$LIKE_HOME/.liked/config/config.toml.new"
mv "$LIKE_HOME/.liked/config/config.toml.new" "$LIKE_HOME/.liked/config/config.toml"
cp "$GENESIS_PATH" "$LIKE_HOME/.liked/config/genesis.json"

docker run --rm -it --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys add validator
ADDRESS=$(docker run --rm -it --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys show validator -a)
echo ""
echo "--------------------------------------------------------------------------------"
echo "Key initialized, your address is $ADDRESS"
echo "Send us this address to get tokens from faucet for staking"