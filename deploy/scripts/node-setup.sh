#!/bin/bash

set -e

MONIKER="$1"
GENESIS_URL="$2"
LIKED_WORKDIR="$3"
LIKED_HOME="$4"
LIKED_USER="$5"

if [ -z $MONIKER ]; then
	echo "Usage: $0 NODE_NAME <url to genesis.json>"
	echo "Example: $0 likecoin-test"
	echo "Example: $0 likecoin-test https://example.com/genesis.json"
	exit 1
fi

if [ -z $LIKED_VERSION ]; then
	LIKED_VERSION="2.0.0-alpha"
fi

if [ -z $COSMOVISOR_VERSION ]; then
	COSMOVISOR_VERSION="1.1.0"
fi

if [ -z $LIKED_USER ]; then
	LIKED_USER="${USER}"
fi

if [ -z $LIKED_WORKDIR ]; then
	LIKED_WORKDIR="$(HOME)"
fi

if [ -z $LIKED_HOME ]; then
	LIKED_HOME="${HOME}/.liked"
fi 

if [ ! -f "$LIKED_WORKDIR/cosmovisor" ]; then
	echo "Downloading the latest cosmovisor binary..."
	mkdir -p cosmovisor_temp
	cd cosmovisor_temp
	curl -sL "https://github.com/cosmos/cosmos-sdk/releases/download/cosmovisor%2Fv$COSMOVISOR_VERSION/cosmovisor-v$COSMOVISOR_VERSION-$(uname -s)-$(uname -m | sed "s|x86_64|amd64|").tar.gz" | tar zx
	cp cosmovisor $LIKED_WORKDIR/cosmovisor
	cd ..
	rm -r cosmovisor_temp
fi 

if [ ! -f "$LIKED_WORKDIR/liked" ]; then
	echo "Downloading the latest liked binary..."
	mkdir -p liked_temp
	cd liked_temp
	curl -sL "https://github.com/likecoin/likecoin-chain/releases/download/v${LIKED_VERSION}/likecoin-chain_${LIKED_VERSION}_$(uname -s)_$(uname -m).tar.gz" | tar xz
	cp bin/liked $LIKED_WORKDIR
	cd ..
	rm -r liked_temp
fi 

if [ ! -f "$LIKED_HOME/config/genesis.json" ]; then
	$LIKED_WORKDIR/liked --home "$LIKED_HOME" init "$MONIKER"
else
	echo "Like instance already initialized, skipping..."
fi

if [ ! -z $GENESIS_URL ]; then
	mkdir -p "$LIKED_HOME/config/"
	curl -OL "$GENESIS_URL"
	mv -f "genesis.json" "$LIKED_HOME/config/genesis.json"
	CHAIN_ID=`grep chain_id "$LIKED_HOME/config/genesis.json" | sed 's/ *"chain_id": *"\(.*\)"/\1/g' | sed 's/,$//g'`
else
	CHAIN_ID="likecoin-devnet-1"
fi

if [ ! -f "$LIKED_HOME/cosmovisor/genesis/bin/liked" ]; then
	echo "Copying binary to cosmovisor genesis"
	mkdir -p "$LIKED_HOME/cosmovisor/genesis/bin"
	cp "$LIKED_WORKDIR/liked" "$LIKED_HOME/cosmovisor/genesis/bin/liked"
fi

chown -R $LIKED_USER $LIKED_HOME

sed "s|<USER>|$LIKED_USER|g; s|<WORKDIR>|$LIKED_WORKDIR|g;" ./liked.service.template > ./liked.service
echo "Setup complete, Please setup DAEMON_NAME and DAEMON_HOME environment variable and run 'cosmovisor run start' to start a node locally." 