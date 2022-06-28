package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/likecoin/likecoin-chain/v3/testutil"

	bech32migrationtestutil "github.com/likecoin/likecoin-chain/v3/bech32-migration/testutil"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	bech32authmigration "github.com/likecoin/likecoin-chain/v3/bech32-migration/auth"
	bech32govmigration "github.com/likecoin/likecoin-chain/v3/bech32-migration/gov"
	bech32slashingmigration "github.com/likecoin/likecoin-chain/v3/bech32-migration/slashing"
	bech32stakingmigration "github.com/likecoin/likecoin-chain/v3/bech32-migration/staking"
)

type MTAppOptions struct{}

// Get implements AppOptions
func (ao MTAppOptions) Get(o string) interface{} {
	if o == crisis.FlagSkipGenesisInvariants {
		return true
	}
	return nil
}

func main() {
	// Skip this test if genesis.json is not found
	if len(os.Args) < 3 || os.Args[1] == "" || os.Args[2] == "" {
		panic(fmt.Errorf("Usage: `go run migrate in_genesis.json out_genesis.json`"))
	}
	inFilePath := os.Args[1]
	outFilePath := os.Args[2]

	if _, err := os.Stat(inFilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic(fmt.Errorf("input file does not exist: %s", err.Error()))
		}
		panic(fmt.Errorf("input genesis json found but failed to stat: %s", err.Error()))
	}

	if _, err := os.Stat(outFilePath); err == nil || !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("output genesis json already exists"))
	}

	jsonFile, err := os.Open(inFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to open json file: %s", err.Error()))
	}
	defer jsonFile.Close()

	exportedBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(fmt.Errorf("failed to read json file: %s", err.Error()))
	}

	var exportedState map[string]json.RawMessage
	if err := json.Unmarshal(exportedBytes, &exportedState); err != nil {
		panic(fmt.Errorf("failed to unmarshal json file: %s", err.Error()))
	}

	// dealloc unused large var
	exportedBytes = nil

	// Init test app and inject genesis
	fmt.Printf("> setup test app\n")
	app := testutil.SetupTestAppWithState(exportedState["app_state"], MTAppOptions{})
	ctx := app.Context
	keys := app.GetKeys()
	appCodec := app.AppCodec()

	// dealloc unused large var
	exportedState = nil

	// Apply upgrade
	fmt.Printf("> apply upgrade\n")
	ctx = ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())
	bech32stakingmigration.MigrateAddressBech32(ctx, keys[stakingtypes.StoreKey], appCodec)
	bech32slashingmigration.MigrateAddressBech32(ctx, keys[slashingtypes.StoreKey], appCodec)
	bech32govmigration.MigrateAddressBech32(ctx, keys[govtypes.StoreKey], appCodec)
	bech32authmigration.MigrateAddressBech32(ctx, keys[authtypes.StoreKey], appCodec)
	app.Commit()

	// Assert
	fmt.Printf("> assert address prefix in stores\n")
	if ok := bech32migrationtestutil.AssertAuthAddressBech32(ctx, keys[authtypes.StoreKey], appCodec); !ok {
		fmt.Printf("!! Assert auth addresses failed\n")
	}
	if ok := bech32migrationtestutil.AssertGovAddressBech32(ctx, keys[govtypes.StoreKey], appCodec); !ok {
		fmt.Printf("!! Assert gov addresses failed\n")
	}
	if ok := bech32migrationtestutil.AssertSlashingAddressBech32(ctx, keys[slashingtypes.StoreKey], appCodec); !ok {
		fmt.Printf("!! Assert slashing addresses failed\n")
	}
	if ok := bech32migrationtestutil.AssertStakingAddressBech32(ctx, keys[stakingtypes.StoreKey], appCodec); !ok {
		fmt.Printf("!! Assert staking addresses failed\n")
	}

	// Export genesis
	fmt.Printf("> export upgraded genesis\n")
	exportedApp, err := app.ExportAppStateAndValidators(false, []string{})
	if err != nil {
		panic(fmt.Errorf("failed to export app state: %s", err.Error()))
	}

	if err := os.WriteFile(outFilePath, exportedApp.AppState, 0644); err != nil {
		panic(fmt.Errorf("failed to write result file: %s", err.Error()))
	}

	// dealloc unused large var
	exportedApp = types.ExportedApp{}
}
