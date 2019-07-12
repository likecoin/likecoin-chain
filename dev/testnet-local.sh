#!/bin/bash

set -e

CHAIN_ID="likechain-local-testnet"
MONIKER="local-dev"
PASSWORD="password"

LIKE_HOME=$(dirname "$0")

mkdir -p "$LIKE_HOME"
pushd "$LIKE_HOME" > /dev/null
LIKE_HOME="$(pwd)"
popd > /dev/null

if [ ! -d "$LIKE_HOME/.likecli" ]; then
    mkdir -p "$LIKE_HOME/.likecli"
    printf "$PASSWORD\n$PASSWORD\n" | docker run --rm -i --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys add validator
    printf "$PASSWORD\n$PASSWORD\n" | docker run --rm -i --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys add faucet
fi

if [ ! -d "$LIKE_HOME/.liked" ]; then
    mkdir -p "$LIKE_HOME/.liked"
    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain liked init --chain-id "$CHAIN_ID" "$MONIKER" > /dev/null 2>&1
    cp "$LIKE_HOME/genesis.json" "$LIKE_HOME/.liked/config/genesis.json"
    # not using sed -i since different behaviour on Linux and Mac
    sed "s/persistent_peers *=.*/persistent_peers=\"$PEERS\"/g" "$LIKE_HOME/.liked/config/config.toml" \
    | sed "s/create_empty_blocks_interval *=.*/create_empty_blocks_interval=\"60s\"/g" \
    | sed "s/^timeout_commit *=.*/timeout_commit=\"60s\"/g" \
    > "$LIKE_HOME/.liked/config/config.toml.new"
    mv "$LIKE_HOME/.liked/config/config.toml.new" "$LIKE_HOME/.liked/config/config.toml"

    VALIDATOR_ADDRESS=`docker run --rm --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys show validator -a`
    FAUCET_ADDRESS=`docker run --rm --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys show faucet -a`
    TM_PUBKEY=`docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain liked tendermint show-validator`

    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked add-genesis-account "$VALIDATOR_ADDRESS" 1000000000000000nanolike

    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked add-genesis-account "$FAUCET_ADDRESS" 849000000000000000nanolike

    printf "$PASSWORD\n" | \
    docker run --rm -i --volume "$LIKE_HOME/.liked:/root/.liked" --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain \
        liked gentx \
            --name validator \
            --amount 1000000000000000nanolike \
            --details "Only for local development" \
            --pubkey "$TM_PUBKEY"
    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked collect-gentxs
else
    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked unsafe-reset-all
fi