package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	captypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evtypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibcxfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	ibccoretypes "github.com/cosmos/cosmos-sdk/x/ibc/core/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	v038 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v038"
	v039 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v039"
	v040 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v040"

	"github.com/likecoin/likechain/cmd/liked/cmd/oldgenesis"

	iscntypes "github.com/likecoin/likechain/x/iscn/types"
)

const flagGenesisTime = "genesis-time"
const flagIscnRegistryId = "iscn-registry-id"
const flagIscnFeePerByte = "iscn-fee-per-byte"

func migrateState(initialState types.AppMap, ctx client.Context, iscnParams iscntypes.Params) types.AppMap {
	state := initialState
	state = v038.Migrate(state, ctx)
	state = v039.Migrate(state, ctx)
	state = v040.Migrate(state, ctx)
	delete(state, "whitelist")

	var stakingGenesis stakingtypes.GenesisState
	ctx.JSONMarshaler.MustUnmarshalJSON(state[stakingtypes.ModuleName], &stakingGenesis)
	stakingGenesis.Params.HistoricalEntries = 10000

	ibcTransferGenesis := ibcxfertypes.DefaultGenesisState()
	ibcCoreGenesis := ibccoretypes.DefaultGenesisState()
	capGenesis := captypes.DefaultGenesis()
	evGenesis := evtypes.DefaultGenesisState()

	ibcTransferGenesis.Params.ReceiveEnabled = false
	ibcTransferGenesis.Params.SendEnabled = false

	ibcCoreGenesis.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}

	iscnGenesis := iscntypes.NewGenesisState(iscnParams, nil, nil)

	state[ibcxfertypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(ibcTransferGenesis)
	state[host.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(ibcCoreGenesis)
	state[captypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(capGenesis)
	state[evtypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(evGenesis)
	state[stakingtypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(&stakingGenesis)
	state[iscntypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(iscnGenesis)
	state[paramstypes.ModuleName] = params.AppModuleBasic{}.DefaultGenesis(ctx.JSONMarshaler)
	state[upgradetypes.ModuleName] = upgrade.AppModuleBasic{}.DefaultGenesis(ctx.JSONMarshaler)
	return state
}

func MigrateGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [genesis-file-from-sheungwan]",
		Short: "Migrate genesis from SheungWan to FoTan",
		Long: (`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ liked migrate /path/to/genesis.json --chain-id=likecoin-chain-fotan --genesis-time=2021-12-31T04:00:00Z --iscn-fee-per-byte=1234.560000000000000000nanolike
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			importGenesis := args[0]

			oldGenDoc, err := oldgenesis.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState types.AppMap
			if err := json.Unmarshal(oldGenDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			newGenDoc := tmtypes.GenesisDoc{}
			newGenDoc.AppHash = oldGenDoc.AppHash
			newGenDoc.ConsensusParams = tmtypes.DefaultConsensusParams()
			newGenDoc.ConsensusParams.Block = oldGenDoc.ConsensusParams.Block
			newGenDoc.ConsensusParams.Validator = oldGenDoc.ConsensusParams.Validator
			newGenDoc.Validators = oldGenDoc.Validators

			// TODO: custom  block height?
			iscnRegistryId, _ := cmd.Flags().GetString(flagIscnRegistryId)
			iscnFeePerByteStr, _ := cmd.Flags().GetString(flagIscnFeePerByte)
			iscnFeePerByte, err := sdk.ParseDecCoin(iscnFeePerByteStr)
			if err != nil {
				return errors.Wrap(err, "failed to parse ISCN fee per byte parameter")
			}
			iscnParam := iscntypes.Params{
				RegistryId: iscnRegistryId,
				FeePerByte: iscnFeePerByte,
			}
			newGenState := migrateState(initialState, clientCtx, iscnParam)

			newGenDoc.AppState, err = json.Marshal(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime, _ := cmd.Flags().GetString(flagGenesisTime)
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				newGenDoc.GenesisTime = t
			} else {
				newGenDoc.GenesisTime = oldGenDoc.GenesisTime
			}

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID != "" {
				newGenDoc.ChainID = chainID
			} else {
				newGenDoc.ChainID = oldGenDoc.ChainID
			}

			bz, err := tmjson.Marshal(newGenDoc)
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flags.FlagChainID, "", "override chain_id with this flag")
	cmd.Flags().String(flagIscnRegistryId, iscntypes.DefaultRegistryId, "ISCN registry ID parameter in the migrated genesis state")
	cmd.Flags().String(flagIscnFeePerByte, iscntypes.DefaultFeePerByte.String(), "ISCN fee per byte parameter in the migrated genesis state")

	return cmd
}
