# Changelog

## [v2.0.2](https://github.com/likecoin/likecoin-chain/releases/v2.0.2)

- Further hotfix on panic when unbonding unbonded validator

Please read https://github.com/likecoin/likecoin-chain/pull/75 for more info on the issue

## [v2.0.1](https://github.com/likecoin/likecoin-chain/releases/v2.0.1)

- Patch chain halt due to double unbounding validator
- Rename go package to likecoin/likecoin-chain/v2

## [v2.0.0](https://github.com/likecoin/likecoin-chain/releases/v2.0.0)

- Upgrade cosmos-sdk to 0.44.8, ibc-go to 2.1.0
- Add support for like account prefix (by LikerLand)
- Add migrations for bech32 stored in state (by LikerLand)
- Add liked debug convert-prefix command tool for bech32 conversion

## [v1.2.0](https://github.com/likecoin/likecoin-chain/releases/v1.2.0)

- Introduce the support of cosmovisor

## [fotan-1.2](https://github.com/likecoin/likecoin-chain/releases/fotan-1.2)
- Upgrade to Cosmos SDK v0.42.11
- Improve chain node performance thanks to [Enforcement of version consistency](https://github.com/likecoin/likecoin-chain/pull/39)
- Cherry pick upgrade module update in cosmos-sdk, support Cosmosvisor 1.x auto download

## [fotan-1.1](https://github.com/likecoin/likecoin-chain/releases/fotan-1.1)
- Upgrade to Cosmos SDK v0.42.9
- Fixes node crash due to concensus error caused by IBC module

## [fotan-1.0](https://github.com/likecoin/likecoin-chain/releases/fotan-1.0)
- Upgrade to Cosmos SDK v0.42.7
- Add ISCN module
- Add IBC support

## [sheungwan-2](https://github.com/likecoin/likecoin-chain/releases/sheungwan-2)
- Upgrade to Cosmos SDK v0.37.15
- Add halt-height and halt-time option support for future chain upgrade

## [sheungwan-1](https://github.com/likecoin/likecoin-chain/releases/sheungwan-1)
- Intial release of LikeCoin chain
- Cosmos SDK v0.37.4
