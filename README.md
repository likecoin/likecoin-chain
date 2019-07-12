LikeChain is a blockchain built on the Cosmos SDK.

## Requirements

 - At least 20 GB disk space
 - Docker
 - Docker Compose

## Building

Run `./build.sh`.

## Running testnet node as a validator

1. Initialize keys by running `./init.sh [moniker] [path-to-genesis.json] [persistent-node]`, where `moniker` is the custom identifier of your node, `path-to-genesis.json` is the path to `genesis.json`, and `persistent-node` is the node ID and IP address of the test node.

Example: `./init.sh chung ~/Downloads/genesis-likechain-cosmos-testnet-1.json '7c93876c5ffce59b5bc07a4b4b7891dd0bfe4cea@35.226.174.222:26656'`

2. After step 1, `docker-compose.yml` will be created, and a Cosmos address for the validator will be initialized. Send the address to us and we will send some token into the account for staking.

3. Start the node by running `docker-compose up -d`. Note that the node is still not a validator, you need to stake the token after receiving it from us for becoming a validator.

4. After receiving tokens, you can stake them by running `./staking.sh`.