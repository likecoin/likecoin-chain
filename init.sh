#!/bin/bash

set -e

node_count=${NODE_COUNT:=1}
tendermint_image="likechain/tendermint"

rm -rf tendermint/nodes/*

is_osx () {
    [[ "$OSTYPE" =~ ^darwin ]] || return 1
}

init() {
    SED=sed
    if is_osx; then
        SED=gsed
        if ! which gsed &> /dev/zero ; then
            brew install gnu-sed
        fi

        if ! which jq &> /dev/zero; then
            brew install jq
        fi
    else
        if ! which jq &> /dev/zero; then
            sudo apt-get install jq -y
        fi
    fi

    if is_osx; then
        rm -rf *data
    else
        sudo rm -rf *data
    fi
}
init

default_genesis="./tendermint/nodes/1/config/genesis.json"
persistent_peers=""

# Reset configs
cat docker-compose.yml > docker-compose.override.yml

# Setup config for each node
for (( i = 1; i <= $node_count; i++ )); do
    echo "Initializing node $i..."

    mkdir -p tendermint/nodes/${i}
    chmod 777 tendermint/nodes/${i}

    go run tendermint/cli/main.go --output_dir tendermint/nodes/${i}/config --type secp256k1
    docker run --rm -v `pwd`/tendermint/nodes/${i}:/tendermint $tendermint_image init

    chmod -R 777 tendermint/nodes/${i}

    node_id=$(docker run --rm -v `pwd`/tendermint/nodes/${i}:/tendermint $tendermint_image show_node_id)
    persistent_peers="$persistent_peers$node_id@tendermint-$i:26656"
    if [[ $i < $node_count ]]; then
        persistent_peers="$persistent_peers,"
    fi

    if [[ $i != 1 ]]; then
        echo $(cat $default_genesis | jq ".validators |= .+ $(cat tendermint/nodes/${i}/config/genesis.json | jq '.validators')") > $default_genesis

        port=$((26660 + $i))

        echo "
  abci-${i}:
    <<: *abci-node
    container_name: likechain_abci-${i}
  tendermint-${i}:
    <<: *tendermint-node
    container_name: likechain_tendermint-${i}
    hostname: tendermint-${i}
    depends_on:
        - abci-${i}
    volumes:
        - ./tendermint/nodes/${i}:/tendermint
    ports:
        - ${port}:26657
    command:
        - --proxy_app=tcp://abci-${i}:26658" >> ./docker-compose.override.yml
    fi

    echo $(cat $default_genesis | jq ".validators[$i-1].name = \"tendermint-$i\" ") > $default_genesis
done

# Apply pesistent_peers option
$SED -i "s/persistent_peers=.*/persistent_peers=$persistent_peers/g" ./docker-compose.override.yml

for (( i = 1; i <= $node_count; i++ )); do
    # Sync all genesis.json of all nodes
    if [[ $i > 1 ]]; then
        cp -f $default_genesis ./tendermint/nodes/${i}/config/genesis.json
    fi

    # Update config.toml
    $SED -i "s/index_all_tags =.*/index_all_tags = true/g" ./tendermint/nodes/${i}/config/config.toml
done
