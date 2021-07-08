LikeCoin chain is a blockchain built on the Cosmos SDK. Project page: https://like.co/

## Requirements

 - At least 40 GB disk space (SSD preferred)
 - Docker
 - Docker Compose with version >= 1.28

## Building Docker image

Normally you don't need to build the image by yourself, as the image is already hosted on Docker Hub.

For building the image, run `./build.sh`. This will build and tag the iamge.

## Setting up a full node

1. Get the URL of the genesis file and other parameters (e.g. seed node) of the network.
1. Copy `.env.template` to `.env`, and also `docker-compose.yml.template` to `docker-compose.yml`.
1. Edit `.env` for config on `LIKECOIN_CHAIN_ID`, `LIKECOIN_MONIKER`, `LIKECOIN_GENESIS_URL` and `LIKECOIN_SEED_NODES`. See comments in the file.
1. Run `docker-compose run --rm init` to setup node data in `.liked` folder.
1. Run `docker-compose up -d` to start up the node and wait for synchronization.
1. Then you may check the logs by `docker-compose logs --tail 1000 -f`.

## Setting up a validator node

1. Setup a full node by following the section above.
1. Make sure the node is synchronized, by checking `localhost:26657/status` and see if `result.sync_info.catching_up` is `false`.
1. Setup validator key by `docker-compose run --rm liked-command keys add validator` and follow the instructions. This will generate a key named `validator` in the keystore.
1. Get the address and mnemonic words from the output of the command above. Jot down the address (`cosmos1...`) and backup the mnemonic words.
1. Get some LIKE in the address above. The LIKE tokens are needed for creating validator.
1. Run `docker-compose run --rm create-validator --amount <AMOUNT> --details <DETAILS> --commission-rate <COMMISSION_RATE>` to create and activate validator. `<AMOUNT>` is the amount for self-delegation (e.g. `100000000000nanolike` for 100 LIKE), `<DETAILS>` is the introduction of the validator, `<COMMISSION_RATE>` is the commission you receive from delegators (e.g. `0.1` for 10%).
1. After sending the create validator transaction, your node should become a validator.
