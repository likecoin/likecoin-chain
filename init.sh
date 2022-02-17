#!/bin/bash

set -e

LIKECOIN_ROOT="$(dirname "$0")"
LIKECOIN_LIKED_HOME="$LIKECOIN_ROOT/.liked"

MONIKER="$1"
GENESIS_PATH="$2"

if [ -z $MONIKER ] || [ -z $GENESIS_PATH ]; then
    echo "Usage: $0 YOUR_NODE_NAME GENESIS_PATH_OR_URL"
    echo
    echo "Example 1:"
    echo "  $0 'chung-node' ~/Downloads/genesis-likecoin-chain-sheungwan.json"
    echo "Example 2:"
    echo "  $0 my-awesome-likecoin-chain-node https://gist.githubusercontent.com/nnkken/1d1b9d4aae4acb3d835dd3150f546d44/raw/4d97fd471b4bf3be8c5475efbc0361f4926e65e5/genesis.json"
    exit 1
fi

if [ ! -f "$LIKECOIN_LIKED_HOME/config/genesis.json" ]; then
    mkdir -p "$LIKECOIN_LIKED_HOME/config" "$LIKECOIN_LIKED_HOME/data"
    liked --home "$LIKECOIN_LIKED_HOME" init "$MONIKER" > /dev/null 2>&1
    echo "Initialized ./.liked folder as node data folder."

    LIKECOIN_VALIDATOR_PUBKEY=$(liked --home "$LIKECOIN_LIKED_HOME" tendermint show-validator)
    echo "LIKECOIN_VALIDATOR_PUBKEY=\"$LIKECOIN_VALIDATOR_PUBKEY\"" >> "$LIKECOIN_ROOT/.env"

    echo "Validator public key ($LIKECOIN_VALIDATOR_PUBKEY) has been written into .env file."
else
    echo "$LIKECOIN_LIKED_HOME/config/genesis.json exists, skipping node initialization."
fi

echo ""

GENESIS_OUTPUT="$LIKECOIN_LIKED_HOME/config/genesis.json"
if [[ $GENESIS_PATH == "http"* ]]; then
    echo "The genesis path is a URL, downloading."
    WGET_PATH=$(which wget || true)
    if [ ! -z $WGET_PATH ]; then
        echo "wget found, downloading using wget."
        wget -O "$GENESIS_OUTPUT" "$GENESIS_PATH" > /dev/null 2>&1
    else
        CURL_PATH=$(which curl || true)
        if [ ! -z $CURL_PATH ]; then
            echo "curl found, downloading using curl."
            curl -Lo "$GENESIS_OUTPUT" "$GENESIS_PATH" > /dev/null 2>&1
        else
            echo "Cannot find curl or wget. Please download the genesis file on your own."
            exit 1
        fi
    fi
else
    echo "The genesis path is a local path, copying."
    cp "$LIKECOIN_ROOT/$GENESIS_PATH" "$GENESIS_OUTPUT"
fi

echo "Create cosmovisor directories"
mkdir "$LIKECOIN_LIKED_HOME/cosmovisor/upgrades"

echo ""
echo "Genesis file installed into .liked/config folder."
echo "Please verify with the SHA256 checksum:"
sha256sum "$GENESIS_OUTPUT"
