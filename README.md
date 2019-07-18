LikeChain is a blockchain built on the Cosmos SDK.

## Requirements

 - At least 20 GB disk space
 - Docker
 - Docker Compose

## Building

Run `./scripts/build.sh`.

## Running testnet node as a validator

1. Initialize keys by running `./scripts/init.sh [moniker] [path-to-genesis.json] [persistent-node]`, where `moniker` is the custom identifier of your node, `path-to-genesis.json` is the path to `genesis.json`, and `persistent-node` is the node ID and IP address of the test node.

Example: `./scripts/init.sh chung ~/Downloads/genesis-likechain-cosmos-testnet-1.json '7c93876c5ffce59b5bc07a4b4b7891dd0bfe4cea@35.226.174.222:26656'`

2. After step 1, a Cosmos address for the validator will be initialized. Send the address to us and we will send some token into the account for staking.

3. Start the node by running `docker-compose up -d`. Note that the node is still not a validator, you need to stake the token after receiving it from us for becoming a validator.

4. After receiving tokens, you can stake them by running `./scripts/staking.sh`.

## Development

 - Setup or reset the one node local testnet by running `./dev/testnet-local.sh`.
 - Use the `docker-compose.yml` in `dev` to run a local server with light client.
 - When code is updated and `go.mod` and `go.sum` are not updated, you can use `./docker/app/build.sh` to quickly rebuild the image.