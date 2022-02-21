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
