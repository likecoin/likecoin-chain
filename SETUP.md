# Node Setup

This document describe how to setup cosmovisor and respective likecoin node. The intended audience is developer would like to test out locally. For a more validator production setup document, one should found at https://docs.like.co

The developer team maintains 3 way on running the node software, they are:

1. `systemd` based.
2. `docker` and `docker-compose` based.
3. Run at shell.

It depends on the situation which one is better. Since the development resource is limiting, we are planning to drop maintaining `docker-compose` way. So if you are new to the ecosystem, I recommend to check out the `systemd` way.

## systemd based, recommended

Following is assuming your OS is Linux based with systemd installed.

Following command will download a pre-compiled binary from Github. For compiling locally, please checkout `RELEASE.md` for details.

To setup a node using our setup script, please run the following command in the project root folder

```
make -C deploy setup-node MONIKER=$MONIKER
```

After the initialization, you may run the following commands to start the node as a service.

```
make -C deploy initialize-systemctl

make -C deploy start-node
```

Above command is warp around `systemctl`. For checking logs, you can run

```
journalctl -u liked.service -f
```

For more advance usage, please checkout the doc of `systemctl` and `journalctl`.

**NOTE: If any existing liked instance is running, Please terminate them before starting a new node to prevent double signing**

For converting the full node into validator, you should able to interact with `liked` directly.

### Advanced Setup

By default, the [node-setup.sh](../deploy/scripts/node-setup.sh) script would download the latest cosmovisor and like binary to the working folder to setup the node automatically. Please refer to the table below for environment variables you can override.

| Env                | Description                                                                         | Default                              |
| ------------------ | ----------------------------------------------------------------------------------- | ------------------------------------ |
| MONIKER            | Moniker identifier for the node, this is mandatory for the node setup               | Empty                                |
| LIKED_VERSION      | Release version of the latest like binary                                           | Latest tag using git describe --tags |
| COSMOVISOR_VERSION | Release version of the latest cosmovisor binary                                     | 1.1.0                                |
| LIKED_WORKDIR      | Working directory, binaries will be downloaded here                                 | $HOME                                |
| LIKED_HOME         | Home directory for the like node, chain data and configurations will be stored here | $HOME/.liked                         |
| LIKED_USER         | User used for 'liked.service' to run on behalf of                                   | $USER                                |

## Docker based

**NOTE: The docker image is only for amd64, arm(i.e. Apple M1) is currently unsupported.**

### Setting up a full node

1. Get the URL of the genesis file and other parameters (e.g. seed node) of the network.
2. Copy `.env.template` to `.env`, and also `docker-compose.yml.template` to `docker-compose.yml`.
3. Edit `.env` for config on `LIKECOIN_CHAIN_ID`, `LIKECOIN_MONIKER`, `LIKECOIN_GENESIS_URL` and `LIKECOIN_SEED_NODES`. See comments in the file.
4. Run `docker-compose run --rm init` to setup node data in `.liked` folder.
5. Run `docker-compose up -d` to start up the node and wait for synchronization.
6. Then you may check the logs by `docker-compose logs --tail 1000 -f`.

### Converting a full not into validator node

1. Setup a full node by following the section above.
2. Make sure the node is synchronized, by checking `localhost:26657/status` and see if `result.sync_info.catching_up` is `false`.
3. Setup validator key by `docker-compose run --rm liked-command keys add validator` and follow the instructions. This will generate a key named `validator` in the keystore.
4. Get the address and mnemonic words from the output of the command above. Jot down the address (`cosmos1...`) and backup the mnemonic words.
5. Deposit some LIKE in the above address. LIKE tokens are needed for creating validator.
6. Run `docker-compose run --rm create-validator --amount <AMOUNT> --details <DETAILS> --commission-rate <COMMISSION_RATE>` to create and activate validator. `<AMOUNT>` is the amount for self-delegation (e.g. `100000000000nanolike` for 100 LIKE), `<DETAILS>` is the introduction of the validator, `<COMMISSION_RATE>` is the commission you receive from delegators (e.g. `0.1` for 10%).
7. After sending the create validator transaction, your node should become a validator.

## Run node at shell locally for development or testing

If you wish to start the node locally, run the following commands

```
MONIKER=<your moniker> make -C deploy setup-node

export DAEMON_NAME=liked
export DAEMON_HOME="$HOME/.liked"

$(LIKED_WORKDIR)/cosmovisor run start
```

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

## Upgrade Proposal

To submit an auto upgrade proposal, execute the following command

```
./liked tx gov submit-proposal software-upgrade $UPGRADED_VERSION \
    --title "$TITLE" \
    --description "$DESCRIPTION" \
    --from $ACCOUNT \
    --upgrade-height $UPGRADE_HEIGHT \
    --upgrade-info '{"binaries":{"linux/amd64":"$BINARY_URL","darwin/amd64":"$BINARY_URL"}}' \
    --deposit 10000000nanolike \
    --chain-id $CHAIN_ID \
    -y
```

Query the proposal to ensure its existance

```
./liked query gov proposal $PROPOSAL_ID
```

Deposite a certain amount of `LIKE` to the proposal

```
./liked tx gov deposit $PROPOSAL_ID 10000000nanolike --from $ACCOUNT --yes
```

Proceed to submit a `Yes` vote for the proposal

```
./liked tx gov vote $PROPOSAL_ID yes --from $ACCOUNT --chain-id $CHAIN_ID -y
```

To submit a manual uprade proposal, remove `--upgrade-info` from the proposal submission command above.
Note that this will require the new binary to be placed in `/.liked/cosmovisor/upgrades/$UPGRADED_VERSION/bin` for the
upgrade to be executed successfully.

# Dependencies

[Cosmovisor](https://docs.cosmos.network/master/run-node/cosmovisor.html)
