# Changelog

## [unreleased]

## [v4.1.1](https://github.com/likecoin/likecoin-chain/releases/v4.1.1)
- Fix missing upgrade handler for 4.1.x upgrade
- Remove unneeded module handlers
- Fix reproducible build not working in macos
- Fix unclear ignite version instruction for protobuf generation and dev env

## [v4.1.0](https://github.com/likecoin/likecoin-chain/releases/v4.1.0)
- Upgrade cosmovisor included in docker image to 1.5.0
- Upgrade ibc-go to 6.2.1

## [v4.0.2](https://github.com/likecoin/likecoin-chain/releases/v4.0.2)
- Upgrade cosmos-sdk to 0.46.15

## [v4.0.1](https://github.com/likecoin/likecoin-chain/releases/v4.0.1)

- Upgrade cosmos-sdk to 0.46.13 which includes barberry security fix:
- Add commands for snapshots management

## [v4.0.0](https://github.com/likecoin/likecoin-chain/releases/v4.0.0)

- Upgrade golang to 1.19.5
- Upgrade cosmos-sdk to 0.46.12
- Upgrade ibc-go to 5.3.1
- Remove x/nft backport
- Add `full_pay_to_royalty` flag to sell NFT event
- Add deterministic ISCN ID when creating ISCN
- Add custom authz message in iscn and likenft module
- Add feegrant in iscn and likenft module messages

## [v3.1.1](https://github.com/likecoin/likecoin-chain/releases/v3.1.1)

- Upgrade cosmos-sdk to 0.45.11
- Upgrade ibc-go to 2.4.2
- Upgrade cosmosvisor in Dockerfile to 1.3.0
- Implement ApplicationQueryService for min gas prices query

## [v3.1.0](https://github.com/likecoin/likecoin-chain/releases/v3.1.0)

- Upgrade cosmos-sdk to 0.45.9 which patches ibc security vulnerability dragonberry

## [v3.0.0](https://github.com/likecoin/likecoin-chain/releases/v3.0.0)

- Upgrade cosmos-sdk to 0.45.6, ibc-go to 2.3.0
- Backport x/nft module from cosmos-sdk 0.46.0-rc1
- Add x/likenft module, support minting NFT via ISCN, blindbox NFT and NFT marketplace

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
