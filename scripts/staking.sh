#!/bin/bash

COMMISSION_RATE="0.5"
COMMISSION_RATE_MAX="0.8"
COMMISSION_RATE_CHANGE="0.10"

set -e

LIKE_HOME="$(dirname "$0")/.."
pushd "$LIKE_HOME" > /dev/null
LIKE_HOME=$(pwd)
popd > /dev/null

CHAIN_ID=$(grep chain_id "$LIKE_HOME/.liked/config/genesis.json" | sed 's/ *"chain_id": *"\(.*\)"/\1/g' | sed 's/,$//g')
echo "Chain ID: $CHAIN_ID"

ADDRESS=$(docker exec likechain_liked likecli keys show validator -a)
echo "Address: $ADDRESS"

VAL_PUBKEY=$(docker exec likechain_liked liked tendermint show-validator)
echo "Validator public key: $VAL_PUBKEY"

CONF_MONIKER=$(grep moniker "$LIKE_HOME/.liked/config/config.toml" | sed 's/ *moniker *= *"\(.*\)"/\1/g')
read -p "Enter your moniker (the name others used to identify your node), default='$CONF_MONIKER':" MONIKER
if [ -z $MONIKER ]; then
    MONIKER=$CONF_MONIKER
fi
echo "Moniker: '$MONIKER'"

BALANCE=$(docker exec likechain_liked likecli query account "$ADDRESS" --chain-id $CHAIN_ID | grep "^\s*Coins" | sed 's/ *Coins: *\([0-9][^ ]*\)/\1/g')
echo "Your balance: $BALANCE"

read -p "Enter the amount you want to stake (including the coin name, example: '1000000000000000nanolike'), or just leave it empty and press Enter to use all balance:" AMOUNT
if [ -z $AMOUNT ]; then
    AMOUNT="$BALANCE"
fi
echo "Staking amount: $AMOUNT"

read -p "Enter some description of your node: " DETAILS

read -p "(Optional) Enter identity of your node (e.g. UPort, Keybase): " IDENTITY
read -p "(Optional) Enter the website of your site (optional): " WEBSITE

echo ""
echo "Now the script will generate and send the stake transaction, please confirm and enter your passphrase."

docker exec -it likechain_liked \
    likecli tx staking create-validator \
        --amount "$AMOUNT" \
        --moniker "$MONIKER" \
        --pubkey "$VAL_PUBKEY" \
        --commission-rate "$COMMISSION_RATE" \
        --commission-max-rate "$COMMISSION_RATE_MAX" \
        --commission-max-change-rate "$COMMISSION_RATE_CHANGE" \
        --identity "$IDENTITY" \
        --details "$DETAILS" \
        --website "$WEBSITE" \
        --min-self-delegation 1 \
        --from validator \
        --chain-id "$CHAIN_ID" \
        --node tcp://liked:26657

echo "Staking transaction sent."
echo "You can use the transaction hash to query the transaction result on the block explorer."