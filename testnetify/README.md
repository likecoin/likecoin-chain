# LikeCoin chain testnetify script

The testnetify script enables users to convert a mainnet genesis export into a testnet genesis, giving them enough voting power to pass any submitted proposals, providing a way to test out upgrade migrations.

## Prerequsite

- Python 3.10
- Genesis export of a mainnet
- Daemon binary

## Setting up

1. Prepare a genesis export of a mainnet `genesis-testnet.json`

2. Adjust variables in the script `testnetify.py` if necessary

3. Run the script, note that the script will use the existing validator key from the daemon folder

```
python ./testnetify.py genesis-testnet.json
```

4. A backup file will be created for the genesis file in the same folder

5. After around 15 minutes, the genesis file will be modified with testnetified data

6. Setup a node by following the [guide](https://docs.like.co/validator/likecoin-chain-node/setup-a-node)

7. Copy the modified `genesis.json` (and `priv_validator_key.json` if the script is executed by a different machine) from the current machine to the target machine's daemon home `HOME/.liked/config`

8. Disable pex by setting `pex=false` in the `HOME/.liked/config/config.toml` file

9. Start the node, p2p errors are expected to be seen after initialization and the node should start generating block after around an hour

## Customization

| Variable Name       | Description                                                                                                                 | Default                                      |
| ------------------- | --------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------- |
| chain_id            | chain id after modification                                                                                                 | mainnet-upgrade-test                         |
| daemon_name         | the name of the binary, this is expected to be in the PATH environment variable                                             | liked                                        |
| minimal_denom       | the denom of the coin the target network is using                                                                           | nanolike                                     |
| operator_prefix     | bech32 prefix of the validator operator address                                                                             | likevaloper                                  |
| consensus_prefix    | bech32 prefix of the validator consensus address                                                                            | likevalcons                                  |
| voting_period       | voting period after modification                                                                                            | 180s                                         |
| delegation_increase | delegation amount increase for the validator                                                                                | 60000000000000000000                         |
| power_increase      | power increase for the validator, this should correspond to the delegation increase                                         | 60000000000000                               |
| balance_increase    | number of native token added to the selected operator                                                                       | 100000000000000000000000                     |
| op_address          | The address of the operator that will be controlling the selected validator, balance_increase will be added to this account | like1ukmjl5s6pnw2txkvz2hd2n0f6dulw34h9rw5zn  |
| op_pubkey           | The address of the operator that will be controlling the selected validator                                                 | AykpD45ZUhhL7tpcNtOdm4+7fPLQcx4u+9OUkfuzN7KT |
