## Introduction

This submission for the hackathon is composed with 2 parts:

 - a Cosmos SDK module for user to record content metadata according to the [ISCN specification](https://github.com/likecoin/iscn-specs/issues)
 - IPFS plugins for querying and parsing these content metadata from the chain

## Setup

If you use Docker, you may use the build script: `scripts/build.sh`

Otherwise, compile the `liked` and `likecli` from `cmd/liked/main.go` and `cmd/likecli/main.go` respectively, similar to gaia.

Some configs in `genesis.json`:

 - app_state.iscn.params.feePerByte: the Coin object representing the fee for uploading records

## Transactions

There are mainly 2 kinds of transactions:

### Add Entity

In the ISCN module, there is a concept of [entity](https://github.com/likecoin/iscn-specs/issues/3#entity).

An entity represents person or organization, e.g. copyright owner, content publisher, etc.

An `AddEntity` message in the transaction looks like this:

```JSON
{
  "type": "likechain/MsgAddEntity",
  "value": {
    "from": "cosmos13sh5wnu006fgmsrsn4uycyx487jfjykvum4rr5",
    "entity": {
      "description": "Blockchain developer",
      "id": "cosmos13sh5wnu006fgmsrsn4uycyx487jfjykvum4rr5",
      "name": "Chung"
    }
  }
}
```

The `id` field could be any string for identifying the entity. `name` and `description` fields are for information.

The `entity` field will be encoded by the lite client using CBOR encoding.

If succeeded, the transaction will emit events with `add_entity` type, and the `entity_cid` field is the CID of the entity.

### Create ISCN

An ISCN record represents a content metadata record. It records the stakeholders (e.g. sharings of each stakeholder) of the content, rights information (e.g. license) of the content, and the information about the content itself (e.g. title, fingerprint).

A `CreateISCN` message in the transaction looks like this:

```JSON
{
  "type": "likechain/MsgCreateISCN",
  "value": {
    "from": "cosmos13sh5wnu006fgmsrsn4uycyx487jfjykvum4rr5",
    "iscnKernel": {
      "content": {
        "fingerprint": "hash://sha256/0baa914a7ad24ba17e9ee470d0b732bf04d9e92afdc76f6a0fd6f9cd2e29e95a",
        "tags": [
          "blog",
          "blockchain"
        ],
        "title": "My Blog",
        "type": "article",
        "version": 1
      },
      "parent": null,
      "rights": {
        "rights": [
          {
            "holder": {
              "/": "z4hviFYYMHNnCnVNTMfazjVFkxTUKB3Qp3XGQf3DtwYpdMfp5AY"
            },
            "period": {
              "from": "2020-01-23T12:34:56Z"
            },
            "terms": {
              "/": "QmZhRNkZaSnhDr6gBC22zwhTjsGyUx39tm8gjFYnTr2SjN"
            },
            "type": "License"
          }
        ]
      },
      "stakeholders": {
        "stakeholders": [
          {
            "sharing": 123,
            "stakeholder": {
              "/": "z4hviFYYMHNnCnVNTMfazjVFkxTUKB3Qp3XGQf3DtwYpdMfp5AY"
            },
            "type": "Creator"
          },
          {
            "sharing": 123,
            "stakeholder": {
              "description": "GitHub's static page host",
              "id": "https://pages.github.com",
              "name": "GitHub Page"
            },
            "type": "Publisher"
          }
        ]
      },
      "timestamp": "2020-01-23T12:34:56Z",
      "version": 1
    }
  }
}
```

Note that for entities, previously added entities could be reused by providing a CID, which is in the following format:

```JSON
{
  "/": "z4hviFYYMHNnCnVNTMfazjVFkxTUKB3Qp3XGQf3DtwYpdMfp5AY"
}
```

The CIDs could be in different encodings (e.g. base58, base32), all are accepted.

If succeeded, the transaction will emit events with `add_iscn_kernel`, `add_iscn_content` types, and the `iscn_kernel_cid` and `iscn_content_cid` fields are the CIDs for the kernels (root record including stakeholders and rights) and the content (dedicated for content metadata, e.g. title) respectively. If entities are nested, `add_entity` events with `entity_cid` field will also be emitted.

There is also an event with type `create_iscn`, the field `iscn_id` is the ID for the ISCN record. Unlike CID, this is not created by hashing the data, but assigned by the chain. This is used to represent the unchanged identifier, and could be useful when the metadata change, since one would want to use an unchanged identifier for the content.

## Queries

In lite client, 2 queries are implemented:

 - `/iscn/kernels/{ISCN_ID}`, which is used for querying the kernel object CID for the given ISCN ID.
   - note that for convenience, the `1/` prefix should be removed, i.e. for ISCN ID `1/zHAEhhiUttvxzmficYUKn3SiZBBYC61gPPKQ7ds6B7BcA`, one should query `/iscn/kernels/zHAEhhiUttvxzmficYUKn3SiZBBYC61gPPKQ7ds6B7BcA`.
 - `/iscn/cids/{CID}`, which is used for querying anything with CID.
   - For some CIDs begin with `/ipfs/`, please remove that prefix. Examples:
     - `/iscn/cids/z4gAY85mbwTANbfr4n5NDjV2Q6isFt28LqD3yMYmYidCcDpQraz` (in base58)
     - `/iscn/cids/bahuaierav3bfvm4ytx7gvn4yqeu4piiocuvtvdpyyb5f6moxniwemae4tjyq` (in base32)

In the module, 3 more queries are implemented: `cid_get`, `cid_get_size`, `cid_has`.

These are for IPFS plugins to query the infomation about the raw bytes of the given CID, and are used in the [cosmosds](https://github.com/likecoin/likecoin-ipfs-cosmosds) repository, so are not exposed from the lite client.

## IPFS plugin: Cosmos Datastore

The main repo of the datastore IPFS plugin (cosmosds) is at https://github.com/likecoin/likecoin-ipfs-cosmosds.

This acts as a datastore of the IPFS software, which connects IPFS with the chain.

If a chain node is running an IPFS software with this datastore, it can provide chain data to other IPFS nodes through the IPFS network.

## IPFS plugin: ISCN-IPLD

The main repo of the ISCN-IPLD IPFS plugin (cosmosds) is at https://github.com/likecoin/iscn-ipld.

It provides decoders to IPFS software for decoding the ISCN data.

## Example testing flow

1. build the Docker image by `scripts/build.sh`
2. `cd testing-data`
  - the `.likecli` folder contains a key `validator`, which the password is `password`, and is used for both validator gentx and the testing transactions
3. `docker-compose up`, which will run both the node and the lite client
4. send the `AddEntity` transaction by POSTing the content of `post-author.json` to the `/txs` endpoint of the lite client. For instance, when using HTTPie: `http POST localhost:1317/txs < post-author.json`
5. record down the `entity_cid` value in `add_entity` event
6. send the `CreateISCN` transaction by POSTing the content of `post-iscn.json` to the `/txs` endpoint of the lite client. For instance, when using HTTPie: `http POST localhost:1317/txs < post-iscn.json`
7. record down the CID related events, and also the `iscn_id` event value
8. query the kernel CID from `localhost:1317/kernels/{ISCN_ID}`, note that the `1/` prefix needs to be removed
9. query the record of the kernel from `localhost:1317/cids/{KERNEL_CID}`
10. query the records of the CID fields from `localhost:1317/cids/{CID}`, note that `/ipfs/` prefix needs to be removed
11. clone the cosmosds repo from https://github.com/likecoin/likecoin-ipfs-cosmosds
12. run the `main.go` example program by `go run main.go`
  - if you have changed the liked port from 26657 to others, please supply the endpoint by `go run main.go 'tcp://localhost:PORT'`
13. enter the CIDs to query the records stored in the chain from the IPFS plugin