package simulation

import (
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/keeper"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func SimulateMsgCreateOffer(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		i := r.Int()
		msg := &types.MsgCreateOffer{
			Creator: simAccount.Address.String(),
			ClassId: strconv.Itoa(i),
			NftId:   strconv.Itoa(i),
		}

		_, found := k.GetOffer(ctx, msg.ClassId, msg.NftId, simAccount.Address)
		if found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Offer already exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgUpdateOffer(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount = simtypes.Account{}
			offer      = types.Offer{}
			msg        = &types.MsgUpdateOffer{}
			allOffer   = types.MapOffersToPublicRecords(k.GetAllOffer(ctx))
			found      = false
		)
		for _, obj := range allOffer {
			simAccount, found = FindAccount(accs, obj.Buyer)
			if found {
				offer = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "offer creator not found"), nil, nil
		}
		msg.Creator = simAccount.Address.String()

		msg.ClassId = offer.ClassId
		msg.NftId = offer.NftId

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgDeleteOffer(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var (
			simAccount = simtypes.Account{}
			offer      = types.Offer{}
			msg        = &types.MsgUpdateOffer{}
			allOffer   = types.MapOffersToPublicRecords(k.GetAllOffer(ctx))
			found      = false
		)
		for _, obj := range allOffer {
			simAccount, found = FindAccount(accs, obj.Buyer)
			if found {
				offer = obj
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "offer creator not found"), nil, nil
		}
		msg.Creator = simAccount.Address.String()

		msg.ClassId = offer.ClassId
		msg.NftId = offer.NftId

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
