# v2.0 migration test

This script is for testing the v2.0 migration.

It imports a genesis state into test app, runs the migration code, then runs test cases and export the resultant state.

## Usage

1. Export state from mainnet / testnet node.

The file structure is: (default structure of `liked export`)

```json
{
  "app_state": {...}
}
```

2. Run the script:

```shell
go run ./migrate.go /path/to/input_state.json /path/to/output_state.json
```

## System requirement

- 8 core cpu
- 32GB RAM + 32GB Swap

Note 64GB of memory will be used in total. Using swap file on lower end machine can avoid OOM crashes but will lengthen the processing time significantly.
