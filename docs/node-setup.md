# Node Setup

## Quickstart (Recommended)

To setup a node using our setup script, please run the following command in the project root folder

```
MONIKER=<your moniker> make -C deploy setup-node
```

After the initialization, you may run the following commands to start the node as a service.

```
make -C deploy initialize-systemctl

make -C deploy start-node
```

**NOTE: If any existing liked instance is running, Please terminate them before starting a new node to prevent double signing**

## Docker

**NOTE: The docker image has updated to be integrated with cosmovisor and is currently unsupported for Apple M1 ARM Macs.**

### Setting up a full node

1. Get the URL of the genesis file and other parameters (e.g. seed node) of the network.
2. Copy `.env.template` to `.env`, and also `docker-compose.yml.template` to `docker-compose.yml`.
3. Edit `.env` for config on `LIKECOIN_CHAIN_ID`, `LIKECOIN_MONIKER`, `LIKECOIN_GENESIS_URL` and `LIKECOIN_SEED_NODES`. See comments in the file.
4. Run `docker-compose run --rm init` to setup node data in `.liked` folder.
5. Run `docker-compose up -d` to start up the node and wait for synchronization.
6. Then you may check the logs by `docker-compose logs --tail 1000 -f`.

### Setting up a validator node

1. Setup a full node by following the section above.
2. Make sure the node is synchronized, by checking `localhost:26657/status` and see if `result.sync_info.catching_up` is `false`.
3. Setup validator key by `docker-compose run --rm liked-command keys add validator` and follow the instructions. This will generate a key named `validator` in the keystore.
4. Get the address and mnemonic words from the output of the command above. Jot down the address (`cosmos1...`) and backup the mnemonic words.
5. Get some LIKE in the address above. The LIKE tokens are needed for creating validator.
6. Run `docker-compose run --rm create-validator --amount <AMOUNT> --details <DETAILS> --commission-rate <COMMISSION_RATE>` to create and activate validator. `<AMOUNT>` is the amount for self-delegation (e.g. `100000000000nanolike` for 100 LIKE), `<DETAILS>` is the introduction of the validator, `<COMMISSION_RATE>` is the commission you receive from delegators (e.g. `0.1` for 10%).
7. After sending the create validator transaction, your node should become a validator.

## Starting node locally (Not recommended)

If you wish to start the node locally, run the following commands

```
MONIKER=<your moniker> make -C deploy setup-node

export DAEMON_NAME=liked
export DAEMON_HOME="$HOME/.liked"

$(LIKED_WORKDIR)/cosmovisor run start
```

**NOTE: This method does not restart the node automatically if the device is restarted under normal circumstances and is not recommended to be used as a validator node.**

## Advanced Setup

By default, the [node-setup.sh](../deploy/scripts/node-setup.sh) script would download the latest cosmovisor and like binary to the working folder to setup the node automatically. Please refer to the table below for environment variables you can override.

| Env                | Description                                                                         | Default                              |
| ------------------ | ----------------------------------------------------------------------------------- | ------------------------------------ |
| MONIKER            | Moniker identifier for the node, this is mandatory for the node setup               | Empty                                |
| LIKED_VERSION      | Release version of the latest like binary                                           | Latest tag using git describe --tags |
| COSMOVISOR_VERSION | Release version of the latest cosmovisor binary                                     | 1.1.0                                |
| LIKED_WORKDIR      | Working directory, binaries will be downloaded here                                 | $HOME                                |
| LIKED_HOME         | Home directory for the like node, chain data and configurations will be stored here | $HOME/.liked                         |
| LIKED_USER         | User used for 'liked.service' to run on behalf of                                   | $USER                                |

# Upgrades

As a process manager, Cosmovisor automatically download and upgrade binary as upgrade proposals are approved on chain.

There are two methods to upgrading the chain version using cosmovisor.

## Automatic Upgrade

An automatic upgrade will be executed if an upgrade proposal is submitted and approved with the upgrade info attached. With this method, no extra action is required and cosmovisor will be able to download the suitable binary and replace the current executable in which the process will restart itself automatically when the upgrade block is reached. However, one should monitor the upgrade process to ensure the upgrade is executed successfully.

## Manual Upgrade

If a manual upgrade is expected, the new binary should be placed inside

```
$HOME/.liked/cosmovisor/upgrades/<upgrade_name>/bin
```

Once the upgrade block is reached, cosmovisor should be able to link the `current` folder to the destinated upgrade folder and automatically restart itself to continue the block syncing process with the latest binary. An extra `upgrade-info.json` file will be generated to indicate the metadata for the upgrade.

# Dependencies

[Cosmovisor](https://docs.cosmos.network/master/run-node/cosmovisor.html)
