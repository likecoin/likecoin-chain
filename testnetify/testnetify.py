# Reference https://github.com/osmosis-labs/osmosis/blob/8043e6571b809aecc8a4b458e0ee6f2ad132e57e/tests/localosmosis/mainnet_state/testnetify.py

import json
import subprocess
import re, shutil, tempfile
import argparse
import collections

## Modifiable parameters 
# --------------------------------------------------
chain_id = "mainnet-upgrade-test"
daemon_name = "liked"
minimal_denom = "nanolike"
operator_prefix = "likevaloper"
consensus_prefix = "likevalcons"
voting_period = "180s" # 3mins

delegation_increase = 60000000000000000000
power_increase = 60000000000000
balance_increase = 100000000000000000000000

# service people choice joy absurd around pony harsh outdoor forget leaf brown mobile notice ozone shed position van flavor lift organ apart assist muffin
op_address = "like1ukmjl5s6pnw2txkvz2hd2n0f6dulw34h9rw5zn"
op_pubkey = "AykpD45ZUhhL7tpcNtOdm4+7fPLQcx4u+9OUkfuzN7KT"
# --------------------------------------------------


def sed_inplace(filename, pattern, repl):
    '''
    Perform the pure-Python equivalent of in-place `sed` substitution: e.g.,
    `sed -i -e 's/'${pattern}'/'${repl}' "${filename}"`.
    '''
    # For efficiency, precompile the passed regular expression.
    pattern_compiled = re.compile(pattern)

    # For portability, NamedTemporaryFile() defaults to mode "w+b" (i.e., binary
    # writing with updating)
    with tempfile.NamedTemporaryFile(mode='w', delete=False) as tmp_file:
        with open(filename) as src_file:
            for line in src_file:
                tmp_file.write(pattern_compiled.sub(repl, line))

    # Overwrite the original file with the munged temporary file in a
    # manner preserving file attributes (e.g., permissions).
    shutil.copystat(filename, tmp_file.name)
    shutil.move(tmp_file.name, filename)


def main():
	parser = argparse.ArgumentParser(description="Likecoin testnetify")
	parser.add_argument("genesis_path")
	args = parser.parse_args()

	testnetify(args.genesis_path)

def testnetify(genesis_path):
	# Load genesis file 
	genesis = read_genesis(genesis_path)
	shutil.copy2(genesis_path, "%s.bak" % genesis_path)

	# Get current validator
	current_validator = get_current_validator()

	# Get eligible validator and delegator
	validator_and_delegator_to_replace = get_valid_validator_and_delegator(genesis, [])

	if validator_and_delegator_to_replace is None:
		print("Found no suitable validator and delegator")
		return

	# Global replace validator information
	validator_hex = validator_and_delegator_to_replace['validator_info']['address']
	validator_operator_address = convert_hex_to_address(validator_hex)
	validator_consensus_address = convert_prefix(validator_operator_address, consensus_prefix)
	replaced_validator_consensus_address = current_validator["bech32_consensus_address"]
	print("Replacing selected validator consensus address {} with {}".format(validator_consensus_address, replaced_validator_consensus_address))
	sed_inplace(genesis_path, validator_consensus_address, replaced_validator_consensus_address)

	validator_hex_address = validator_and_delegator_to_replace['validator_info']['address']
	replaced_validator_hex_address = current_validator["address_hex"]
	print("Replacing selected validator hex address {} with {}".format(validator_hex_address, replaced_validator_hex_address))
	sed_inplace(genesis_path, validator_hex_address, replaced_validator_hex_address)

	delegator_address = validator_and_delegator_to_replace['delegator']['address']
	print("Replacing selected delegator address {} with {}".format(delegator_address, op_address))
	sed_inplace(genesis_path, delegator_address, op_address)

	delegator_pub_key = validator_and_delegator_to_replace['delegator']['pub_key']['key']
	if delegator_pub_key is not None:
		print("Replacing selected delegator public key {} with {}".format(delegator_pub_key, op_pubkey))
		sed_inplace(genesis_path, delegator_pub_key, op_pubkey)

	# Reload genesis
	genesis = read_genesis(genesis_path)

	# Update chain id
	current_chain_id = genesis['chain_id']
	genesis["chain_id"] = chain_id
	print("Replacing chain-id {} with {}".format(current_chain_id, chain_id))


	# Add extra amount of delegation and power to the selected validator and delegator, in which full control of the network will be granted to the selected account
	update_validator_information(genesis, validator_and_delegator_to_replace['validator'], current_validator)
	update_operator_information(genesis)

	# Patch chain information to reflect the changes above
	patch_staking_power(genesis)
	patch_supply(genesis)
	patch_bonded_pool(genesis)

	update_params(genesis)

	# Overwrite genesis file
	genesis_file = open(genesis_path, 'w')
	json.dump(genesis, genesis_file)
	genesis_file.truncate()
	genesis_file.close()

def read_genesis(genesis_path):
	genesis = open(genesis_path, "r+")
	genesis_json = json.loads(genesis.read())
	genesis.close()
	return genesis_json


def get_current_validator():
	validator_pubkey_output_raw = subprocess.run([daemon_name,"tendermint","show-validator"], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
	validator_pubkey_output = validator_pubkey_output_raw.stdout.strip()
	##validator_pubkey_output = '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"3QVAkiUIkKR3B6kkbd+QqzWDdcExoggbZV5fwH4jKDs="}'
	validator_pubkey = validator_pubkey_output[validator_pubkey_output.find('key":') +6 :-2]


	debug_pubkey = subprocess.run([daemon_name,"debug", "pubkey", validator_pubkey_output], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
	address_hex = debug_pubkey.stderr[9: debug_pubkey.stderr.find("\n")]
	##address = "214D831D6F49A75F9104BDC3F2E12A6CC1FC5669"

	# Convert hex to address
	bech32_address = convert_hex_to_address(address_hex)

	# Convert to consensus address
	bech32_consensus_address = convert_prefix(bech32_address, consensus_prefix)
	## bech32_consensus_address = "likevalcons1y9xcx8t0fxn4lygyhhpl9cf2dnqlc4nfkdxg9f"

	return {
		"validator_pubkey": validator_pubkey,
		"address_hex": address_hex,
		"bech32_consensus_address": bech32_consensus_address
	}

def convert_prefix(address, prefix):
	address_output = subprocess.run([daemon_name,"debug", "convert-prefix", address, "-p", prefix], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
	converted_address = address_output.stderr[:address_output.stderr.find("\n")]
	return converted_address

def convert_hex_to_address(hex):
	bech32_address_output = subprocess.run([daemon_name,"debug", "addr", hex], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
	bech32_address = bech32_address_output.stderr[bech32_address_output.stderr.find("Acc: ") + 5: bech32_address_output.stderr.find("Bech32 Val")-1]
	return bech32_address

def get_module_account_by_name(genesis_json, name):
	accounts = genesis_json['app_state']['auth']['accounts']
	account = next(a for a in accounts if "cosmos.auth.v1beta1.ModuleAccount" in a['@type'] and a['name'] == name)
	return account

def get_valid_validator_and_delegator(genesis_json, invalid_validators = []):
	validators = genesis_json['app_state']['staking']['validators']
	valid_validators = [v for v in validators if v['jailed'] == False and v['operator_address'].startswith(operator_prefix) and v['operator_address'] not in invalid_validators]
	if len(valid_validators) == 0:
		return None

	valid_validator = valid_validators[0]
	# Select delegator who only delegated once
	delegation_count = collections.defaultdict(int)
	delegations = genesis_json['app_state']['staking']['delegations']
	for d in delegations: delegation_count[d['delegator_address']] += 1
	valid_delegators = [d for d in delegations if d['validator_address'] == valid_validator['operator_address'] and delegation_count[d["delegator_address"]] == 1]
	if len(valid_delegators) == 0:
		return get_valid_validator_and_delegator(genesis_json, invalid_validators + [valid_validator['operator_address']])

	valid_delegator = valid_delegators[0]
	validator_infos = genesis_json['validators']
	validator_info = next(v for v in validator_infos if v['name'] == valid_validator['description']['moniker'])

	accounts = genesis_json['app_state']['auth']['accounts']
	delegator_account = next(a for a in accounts if 'address' in a and a['address'] == valid_delegator['delegator_address'])

	return {
		"validator_info": validator_info,
		"validator": valid_validator,
		"delegator": delegator_account,
	}

def update_validator_information(genesis_json, target_validator, current_validator):
	# Update validator public key
	validators = genesis_json['app_state']['staking']['validators']
	validator_index = next(i for i, val in enumerate(validators) if val["description"]['moniker'] == target_validator['description']['moniker'])
	validators[validator_index]['consensus_pubkey']['key'] = current_validator['validator_pubkey']

	validator_infos = genesis_json['validators']
	validator_info_index = next(i for i, val in enumerate(validator_infos) if val["address"] == current_validator['address_hex'])
	validator_infos[validator_info_index]['pub_key']['value'] = current_validator['validator_pubkey']

	# Update validator power
	current_power = validator_infos[validator_info_index]['power']
	updated_power = str(int(current_power) + power_increase)
	print("Replacing current validator power from {} to {}".format(current_power, updated_power))
	validator_infos[validator_info_index]['power'] = updated_power

	# Update last power
	last_validator_powers = genesis_json['app_state']['staking']['last_validator_powers']
	last_validator_power_index = next(i for i, val in enumerate(last_validator_powers) if val["address"] == validators[validator_index]['operator_address'])
	last_power = last_validator_powers[last_validator_power_index]['power']
	updated_last_power = str(int(last_power) + power_increase)
	print("Replacing current last validator power from {} to {}".format(last_power, updated_last_power))
	last_validator_powers[last_validator_power_index]['power'] = updated_last_power

	# Update validator shares and tokens
	current_share = validators[validator_index]['delegator_shares']
	updated_share = "{}.000000000000000000".format(int(float(current_share)) + delegation_increase)
	print("Replacing current delegator share from {} to {}".format(current_share, updated_share))
	validators[validator_index]['delegator_shares'] = updated_share

	current_tokens = validators[validator_index]['tokens']
	updated_tokens =str(int(current_tokens) + delegation_increase)
	print("Replacing current tokens from {} to {}".format(current_tokens, updated_tokens))
	validators[validator_index]['tokens'] = updated_tokens

def update_operator_information(genesis_json):
	# Update delegation amount
	delegations = genesis_json['app_state']['staking']['delegations']
	delegation_index = next(i for i, val in enumerate(delegations) if op_address in val["delegator_address"])
	current_share = delegations[delegation_index]['shares']
	updated_share = "{}.000000000000000000".format(int(float(current_share)) + delegation_increase)
	print("Replacing current delegator share from {} to {}".format(current_share, updated_share))
	delegations[delegation_index]['shares'] = updated_share

	delegator_starting_infos = genesis_json['app_state']['distribution']['delegator_starting_infos']
	delegator_starting_info_index = next(i for i, val in enumerate(delegator_starting_infos) if op_address in val["delegator_address"])
	current_stake = delegator_starting_infos[delegator_starting_info_index]['starting_info']['stake']
	updated_stake = "{}.000000000000000000".format(int(float(current_stake)) + delegation_increase)
	print("Replacing current delegator stake from {} to {}".format(current_stake, updated_stake))
	delegator_starting_infos[delegator_starting_info_index]['starting_info']['stake'] = updated_stake

	# Update balance
	balances = genesis_json['app_state']['bank']['balances']
	balance_index = next(i for i, val in enumerate(balances) if op_address in val["address"])
	operator_wallet = balances[balance_index]['coins']
	coin_index = next(i for i, val in enumerate(operator_wallet) if val["denom"] == minimal_denom)
	coin_amount = operator_wallet[coin_index]['amount']
	updated_coin_amount = str(int(coin_amount) + balance_increase)
	print("Replacing current coin amount from {} to {}".format(coin_amount, updated_coin_amount))
	operator_wallet[coin_index]['amount'] = updated_coin_amount

def update_params(genesis_json):
	# Update voting period
	current_voting_period = genesis_json['app_state']['gov']['voting_params']['voting_period']
	print("Replacing current voting period from {} to {}".format(current_voting_period, voting_period))
	genesis_json['app_state']['gov']['voting_params']['voting_period'] = voting_period


def patch_staking_power(genesis_json):
	last_total_power = genesis_json['app_state']['staking']['last_total_power']
	updated_last_total_power = str(int(last_total_power) + power_increase)
	print("Replacing current last total power from {} to {}".format(last_total_power, updated_last_total_power))
	genesis_json['app_state']['staking']['last_total_power'] = updated_last_total_power

def patch_bonded_pool(genesis_json):
	bonded_pool_account = get_module_account_by_name(genesis_json, "bonded_tokens_pool")
	balances = genesis_json['app_state']['bank']['balances']
	module_balance_index = next(i for i, val in enumerate(balances) if bonded_pool_account['base_account']['address'] in val["address"])
	coins = balances[module_balance_index]['coins']
	coin_index = next(i for i, val in enumerate(coins) if val["denom"] == minimal_denom)
	current_coin_balance = coins[coin_index]['amount']
	updated_coin_balance = str(int(current_coin_balance) + delegation_increase)
	print("Patching bonded pool balance from {} to {}".format(current_coin_balance, updated_coin_balance))
	coins[coin_index]['amount'] = updated_coin_balance

def patch_supply(genesis_json):
	supply = genesis_json['app_state']['bank']['supply']
	coin_index = next(i for i, val in enumerate(supply) if val["denom"] == minimal_denom)
	current_coin_supply = supply[coin_index]['amount']
	updated_coin_supply = str(int(current_coin_supply) + delegation_increase + balance_increase)
	print("Patching supply from {} to {}".format(current_coin_supply, updated_coin_supply))
	supply[coin_index]['amount'] = updated_coin_supply
	
main()