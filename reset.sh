#!/bin/bash

for (( i = 1; i <= $NODE_COUNT; i++ )); do
    echo "Resetting node $i..."

    docker-compose run --rm abci-${i} rm -rf /like/app.db /like/state.db /like/withdraw.db 
    docker-compose run --rm --entrypoint "tendermint unsafe_reset_all" --no-deps tendermint-${i} 
done
