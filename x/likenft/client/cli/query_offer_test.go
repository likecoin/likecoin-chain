package cli_test

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/likecoin/likechain/testutil/network"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/client/cli"
	"github.com/likecoin/likechain/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithOfferObjects(t *testing.T, n int) (*network.Network, []types.Offer) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		// Create random address
		pubBz := make([]byte, ed25519.PubKeySize)
		rand.Read(pubBz)
		pub := &ed25519.PubKey{Key: pubBz}
		address, _ := sdk.Bech32ifyAddressBytes("like", pub.Address())
		offer := types.Offer{
			ClassId: strconv.Itoa(i),
			NftId:   strconv.Itoa(i),
			Buyer:   address,
		}
		nullify.Fill(&offer)
		state.OfferList = append(state.OfferList, offer)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.OfferList
}

func TestShowOffer(t *testing.T) {
	net, objs := networkWithOfferObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc      string
		idClassId string
		idNftId   string
		idBuyer   string

		args []string
		err  error
		obj  types.Offer
	}{
		{
			desc:      "found",
			idClassId: objs[0].ClassId,
			idNftId:   objs[0].NftId,
			idBuyer:   objs[0].Buyer,

			args: common,
			obj:  objs[0],
		},
		{
			desc:      "not found",
			idClassId: strconv.Itoa(100000),
			idNftId:   strconv.Itoa(100000),
			idBuyer:   strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idClassId,
				tc.idNftId,
				tc.idBuyer,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowOffer(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryOfferResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Offer)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Offer),
				)
			}
		})
	}
}

func TestListOffer(t *testing.T) {
	net, objs := networkWithOfferObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOffer(), args)
			require.NoError(t, err)
			var resp types.QueryOfferIndexResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Offers), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Offers),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOffer(), args)
			require.NoError(t, err)
			var resp types.QueryOfferIndexResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Offers), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Offers),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListOffer(), args)
		require.NoError(t, err)
		var resp types.QueryOfferIndexResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Offers),
		)
	})
}
