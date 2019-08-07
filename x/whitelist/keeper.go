package whitelist

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramstore params.Subspace
	codespace  sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
		codespace:  codespace,
	}
}

func (keeper Keeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

func (keeper Keeper) GetWhitelist(ctx sdk.Context) (whitelist Whitelist) {
	bz := ctx.KVStore(keeper.storeKey).Get(WhitelistKey)
	if bz == nil {
		return nil
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &whitelist)
	return whitelist
}

func (keeper Keeper) SetWhitelist(ctx sdk.Context, whitelist Whitelist) {
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(whitelist)
	ctx.KVStore(keeper.storeKey).Set(WhitelistKey, bz)
}

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (k Keeper) Approver(ctx sdk.Context) (res sdk.AccAddress) {
	k.paramstore.Get(ctx, KeyApprover, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) Params {
	return Params{
		Approver: k.Approver(ctx),
	}
}

func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
