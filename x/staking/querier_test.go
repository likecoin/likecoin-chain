package staking

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	codec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const StakingDenom = "nanolike"
const BalanceMax = 1000000000
const BlockMax = 300
const ValidatorCount = 10
const DelegatorCount = 100
const MaxActionPerBlock = 10

func sleepMs(ms int) {
	toSleepDuration := time.Duration(int64(ms) * int64(time.Millisecond))
	time.Sleep(toSleepDuration)
}

func coin(n int64) sdk.Coin {
	return sdk.NewCoin(StakingDenom, sdk.NewInt(int64(n)))
}

func createAccount(serial int) (sdk.AccAddress, authtypes.GenesisAccount, sdk.Coin) {
	bz := make([]byte, 20)
	for i := 0; i < 20; i++ {
		if serial == 0 {
			break
		}
		bz[i] = byte(serial % 256)
		serial /= 256
	}
	addr := sdk.AccAddress(bz)
	acc := authtypes.NewBaseAccount(addr, nil, uint64(serial), 0)
	c := coin(BalanceMax)
	return addr, acc, c
}

type TestSetup struct {
	Context       sdk.Context
	DB            db.DB
	CMS           sdk.CommitMultiStore
	AppCodec      codec.Codec
	AccountKeeper *authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	StakingKeeper *keeper.Keeper
	Validators    []sdk.ValAddress
	Delegators    []sdk.AccAddress
}

func setupTest(logger log.Logger) TestSetup {
	var err error
	storeDB := db.NewMemDB()
	multistore := store.NewCommitMultiStore(storeDB)

	legacyAmino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(interfaceRegistry)
	std.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterLegacyAminoCodec(legacyAmino)
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterLegacyAminoCodec(legacyAmino)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterLegacyAminoCodec(legacyAmino)
	stakingtypes.RegisterInterfaces(interfaceRegistry)

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	}
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		paramstypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	paramsKeeper := paramskeeper.NewKeeper(
		appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey],
	)
	authSubspace := paramsKeeper.Subspace(authtypes.ModuleName)
	bankSubspace := paramsKeeper.Subspace(banktypes.ModuleName)
	stakingSubspace := paramsKeeper.Subspace(stakingtypes.ModuleName)

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], authSubspace, authtypes.ProtoBaseAccount, maccPerms,
	)
	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], accountKeeper, bankSubspace, modAccAddrs,
	)
	stakingKeeper := keeper.NewKeeper(appCodec, keys[stakingtypes.StoreKey], accountKeeper, bankKeeper, stakingSubspace)

	for _, k := range keys {
		multistore.MountStoreWithDB(k, sdk.StoreTypeIAVL, nil)
	}

	for _, k := range tkeys {
		multistore.MountStoreWithDB(k, sdk.StoreTypeTransient, nil)
	}
	err = multistore.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	genesisTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	header := tmproto.Header{Height: 1, Time: genesisTime}
	accounts := authtypes.GenesisAccounts{}
	balances := []banktypes.Balance{}
	supply := sdk.NewCoin(StakingDenom, sdk.NewInt(0))

	validators := []sdk.ValAddress{}
	for i := 0; i < ValidatorCount; i++ {
		addr, acc, c := createAccount(i)
		v := sdk.ValAddress(addr)
		validators = append(validators, v)
		accounts = append(accounts, acc)
		balances = append(balances, banktypes.Balance{addr.String(), sdk.NewCoins(c)})
		supply = supply.Add(c)
	}
	delegators := []sdk.AccAddress{}
	for i := 0; i < DelegatorCount; i++ {
		addr, acc, c := createAccount(ValidatorCount + i)
		delegators = append(delegators, addr)
		accounts = append(accounts, acc)
		balances = append(balances, banktypes.Balance{addr.String(), sdk.NewCoins(c)})
		supply = supply.Add(c)
	}

	ctx := sdk.NewContext(multistore, header, false, logger)
	auth.InitGenesis(ctx, accountKeeper, *authtypes.NewGenesisState(authtypes.DefaultParams(), accounts))
	bankKeeper.InitGenesis(ctx, banktypes.NewGenesisState(banktypes.DefaultParams(), balances, sdk.NewCoins(supply), nil))

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = StakingDenom
	stakingParams.MaxValidators = ValidatorCount * 2
	stakingParams.MaxEntries = BlockMax
	stakingParams.UnbondingTime = 1500 * time.Millisecond
	staking.InitGenesis(ctx, stakingKeeper, accountKeeper, bankKeeper, &stakingtypes.GenesisState{
		Params: stakingParams,
	})

	msgServer := keeper.NewMsgServerImpl(stakingKeeper)
	for i, v := range validators {
		addr := sdk.AccAddress(v)
		seed := v
		consPubKey := ed25519.GenPrivKeyFromSecret(seed).PubKey()
		pubKeyAny, err := codectypes.NewAnyWithValue(consPubKey)
		_, err = msgServer.CreateValidator(sdk.WrapSDKContext(ctx), &stakingtypes.MsgCreateValidator{
			ValidatorAddress:  v.String(),
			Description:       stakingtypes.NewDescription(fmt.Sprintf("validator_%04d", i), "", "", "", ""),
			Commission:        stakingtypes.NewCommissionRates(sdk.NewDec(1), sdk.NewDec(1), sdk.NewDec(1)),
			MinSelfDelegation: sdk.NewInt(0),
			DelegatorAddress:  addr.String(),
			Value:             coin(1),
			Pubkey:            pubKeyAny,
		})
		if err != nil {
			panic(err)
		}
	}

	return TestSetup{
		Context:       ctx,
		DB:            storeDB,
		CMS:           multistore,
		AppCodec:      appCodec,
		AccountKeeper: &accountKeeper,
		BankKeeper:    bankKeeper,
		StakingKeeper: &stakingKeeper,
		Validators:    validators,
		Delegators:    delegators,
	}
}

func TestFuzz(t *testing.T) {
	logger := log.NewTMLogger(os.Stdout)

	setup := setupTest(logger)
	t.Cleanup(func() {
		setup.DB.Close()
	})

	indexingBackendDB := db.NewMemDB()
	t.Cleanup(func() {
		indexingBackendDB.Close()
	})
	// indexingDB := &DBDebugWrap{DB: indexingBackendDB, Logger: logger}
	indexingDB := indexingBackendDB
	indexedQuerier := NewQuerier(setup.StakingKeeper, setup.AppCodec, indexingDB)

	setup.StakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(DebugHooks{}, NewHooks(indexedQuerier)))

	headerChan := make(chan tmproto.Header, 1)
	queryDoneChan := make(chan stakingtypes.DelegationResponses, 1)

	t.Run("DeliverTx thread", func(t *testing.T) {
		t.Parallel()
		l := logger.With("thread", "DeliverTx")
		genesisTime := setup.Context.BlockHeader().Time
		msgServer := keeper.NewMsgServerImpl(*setup.StakingKeeper)
		defer close(headerChan)
		delegations := []stakingtypes.DelegationResponse{}
		r := rand.New(rand.NewSource(1337))

		type ActionEntry struct {
			Ratio int
			Run   func(ctx context.Context)
		}

		actionDelegate := func(ctx context.Context) {
			// delegate
			delegator := setup.Delegators[r.Intn(DelegatorCount)].String()
			validator := setup.Validators[r.Intn(ValidatorCount)].String()
			amount := r.Int63n(BalanceMax/BlockMax) + 1
			l.Debug("Testcase doing Delegate", "delegator", delegator, "validator", validator, "amount", amount)
			_, err := msgServer.Delegate(ctx, &stakingtypes.MsgDelegate{
				DelegatorAddress: delegator,
				ValidatorAddress: validator,
				Amount:           coin(amount),
			})
			if err != nil {
				l.Info("Delegate failed", "err", err)
			}
		}

		pickDelegationAndAmount := func() (stakingtypes.DelegationResponse, bool, int64) {
			index := r.Intn(len(delegations))
			delegationResponse := delegations[index]

			delAddr, err := sdk.AccAddressFromBech32(delegationResponse.Delegation.DelegatorAddress)
			if err != nil {
				panic(err)
			}
			valAddr, err := sdk.ValAddressFromBech32(delegationResponse.Delegation.ValidatorAddress)
			if err != nil {
				panic(err)
			}
			isSelfDelegation := delAddr.Equals(valAddr)

			delegationAmount := delegationResponse.Balance.Amount.Int64()
			operationAmount := delegationAmount
			if delegationAmount > 1 && r.Intn(2) == 0 {
				operationAmount = r.Int63n(delegationAmount) + 1
			}

			return delegationResponse, isSelfDelegation, operationAmount
		}

		actionUndelegate := func(ctx context.Context) {
			// undelegate, possibly from the validator itself
			delegationResponse, isSelfDelegation, operationAmount := pickDelegationAndAmount()
			if isSelfDelegation {
				l.Info("Validator undelegating", "validator", delegationResponse.Delegation.ValidatorAddress)
			}
			l.Debug("Testcase doing Undelegate", "delegation", delegationResponse.Delegation, "amount", operationAmount)
			_, err := msgServer.Undelegate(ctx, &stakingtypes.MsgUndelegate{
				DelegatorAddress: delegationResponse.Delegation.DelegatorAddress,
				ValidatorAddress: delegationResponse.Delegation.ValidatorAddress,
				Amount:           coin(operationAmount),
			})
			if err != nil {
				l.Info("Undelegate failed", "err", err)
			}
		}

		actionRedelegate := func(ctx context.Context) {
			delegationResponse, isSelfDelegation, operationAmount := pickDelegationAndAmount()
			if isSelfDelegation {
				l.Info("Validator redelegating", "validator", delegationResponse.Delegation.ValidatorAddress)
			}
			dstValidator := setup.Validators[r.Intn(ValidatorCount)].String()
			for dstValidator == delegationResponse.Delegation.ValidatorAddress {
				dstValidator = setup.Validators[r.Intn(ValidatorCount)].String()
			}
			l.Debug("Testcase doing Redelegate", "delegation", delegationResponse.Delegation, "dst_validator", dstValidator, "amount", operationAmount)
			_, err := msgServer.BeginRedelegate(ctx, &stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    delegationResponse.Delegation.DelegatorAddress,
				ValidatorSrcAddress: delegationResponse.Delegation.ValidatorAddress,
				ValidatorDstAddress: dstValidator,
				Amount:              coin(operationAmount),
			})
			if err != nil {
				l.Info("BeginRedelegate failed", "err", err)
			}
		}

		actions := []ActionEntry{}

		actions = append(actions, ActionEntry{
			1, func(ctx context.Context) {
				actionDelegate(ctx)
				actions = []ActionEntry{
					{900, actionDelegate},
					{70, actionUndelegate},
					{30, actionRedelegate},
				}
			},
		})

		for height := int64(1); height <= BlockMax; height++ {
			blockTime := genesisTime.Add(time.Duration((height - 1) * int64(time.Second)))
			header := tmproto.Header{Height: height, Time: blockTime}
			l.Debug("Generated header", "height", header.Height)
			store := setup.CMS.CacheMultiStore()
			ctx := sdk.NewContext(store, header, false, l)
			// simulate BeginBlock
			// TODO: use BeginBlocker
			indexedQuerier.BeginWriteIndex(ctx)
			wrappedCtx := sdk.WrapSDKContext(ctx)
			actionCount := r.Intn(MaxActionPerBlock)
			for i := 0; i < actionCount; i++ {
				actionSum := 0
				for _, action := range actions {
					actionSum += action.Ratio
				}
				actionIndex := r.Intn(actionSum)
				action := actions[0]
				for i := 0; i < len(actions); i++ {
					if actionIndex >= action.Ratio {
						actionIndex -= action.Ratio
						action = actions[i+1]
					}
				}
				action.Run(wrappedCtx)
			}
			// simulate EndBlock
			// TODO: use EndBlocker
			indexedQuerier.CommitWriteIndex(ctx)
			store.Write()
			setup.CMS.Commit()
			logger.Debug("After committing writes", "height", height, "delegation_count", len(delegations))
			headerChan <- header
			delegations = <-queryDoneChan
		}
	})

	t.Run("Query thread", func(t *testing.T) {
		t.Parallel()
		defer close(queryDoneChan)
		l := logger.With("thread", "Query")
		originalQuerier := keeper.Querier{Keeper: *setup.StakingKeeper}

		doQuery := func(ctx sdk.Context, valAddrStr string, limit uint64) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
			blockHeight := ctx.BlockHeight()
			l.Debug("Testing ValidatorDelegations query", "block_height", blockHeight, "validator", valAddrStr)
			req := &stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: valAddrStr,
				Pagination: &query.PageRequest{
					Limit: limit,
				},
			}
			isSlowPath := blockHeight != int64(indexedQuerier.GetHeight())
			wrappedCtx := sdk.WrapSDKContext(ctx)
			t0 := time.Now()
			indexedRes, indexedErr := indexedQuerier.ValidatorDelegations(wrappedCtx, req)
			t1 := time.Now()
			originalRes, originalErr := originalQuerier.ValidatorDelegations(wrappedCtx, req)
			t2 := time.Now()
			// Can't require indexedErr is the same as originalErr since original querier didn't check validator address format...
			require.Condition(t,
				func() bool {
					return (originalErr == nil) == (indexedErr == nil)
				},
				"Expected indexed querier to return same error with original querier, height = %d, validator = %s",
				blockHeight, valAddrStr,
			)
			if originalErr != nil {
				l.Debug(
					"ValidatorDelegations done with error",
					"is_slow_path", isSlowPath,
					"indexed_querier_time", t1.Sub(t0),
					"original_querier_time", t2.Sub(t1),
					"height", blockHeight,
					"validator", valAddrStr,
					"original_err", originalErr,
					"indexed_err", indexedErr,
				)
				return nil, indexedErr
			}
			// Pagination keys could be different, so compare desired fields one by one
			require.EqualValues(
				t, originalRes.DelegationResponses, indexedRes.DelegationResponses,
				"Expected indexed querier to return same delegation responses with original querier, height = %d, validator = %s",
				blockHeight, valAddrStr,
			)
			require.EqualValues(
				t, originalRes.Pagination.Total, indexedRes.Pagination.Total,
				"Expected indexed querier to return same delegation total count with original querier, height = %d, validator = %s",
				blockHeight, valAddrStr,
			)
			l.Debug(
				"ValidatorDelegations benchmark",
				"is_slow_path", isSlowPath,
				"indexed_querier_time", t1.Sub(t0),
				"original_querier_time", t2.Sub(t1),
				"height", blockHeight,
				"validator", valAddrStr,
				"res_len", len(originalRes.DelegationResponses),
			)
			return originalRes, nil
		}

		testQueries := func(header tmproto.Header) stakingtypes.DelegationResponses {
			store, err := setup.CMS.CacheMultiStoreWithVersion(header.Height)
			if err != nil {
				panic("Test on ValidatorDelegations query skipped due to store version not found")
			}
			ctx := sdk.NewContext(store, header, true, l)
			delegations := stakingtypes.DelegationResponses{}
			for _, v := range setup.Validators {
				res, err := doQuery(ctx, v.String(), 1e10)
				if err != nil {
					continue
				}
				delegations = append(delegations, res.DelegationResponses...)
			}
			// non-existing validator
			accAddr, _, _ := createAccount(123456789)
			valAddrStr := sdk.ValAddress(accAddr).String()
			doQuery(ctx, valAddrStr, 1e10)

			// invalid checksum
			doQuery(ctx, "cosmosvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqkh53tw", 1e10)

			// unsupported prefix
			doQuery(ctx, "dosmosvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqql5ax3d", 1e10)

			// invalid address strings
			doQuery(ctx, "asdf", 1e10)
			doQuery(ctx, "", 1e10)

			return delegations
		}
		queryHeader := <-headerChan
		l.Debug("Received header", "height", queryHeader.Height)
		delegations := testQueries(queryHeader)
		oldQueryHeader := queryHeader
		queryDoneChan <- delegations
		for queryHeader = range headerChan {
			l.Debug("Received header", "height", queryHeader.Height)
			testQueries(oldQueryHeader)
			delegations := testQueries(queryHeader)
			oldQueryHeader = queryHeader
			queryDoneChan <- delegations
			// trying to trigger write-during-query case
			testQueries(oldQueryHeader)
			sleepMs(1)
			testQueries(oldQueryHeader)
		}
	})
}
