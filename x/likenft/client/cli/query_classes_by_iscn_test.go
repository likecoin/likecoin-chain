package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/likecoin/likechain/testutil/network"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/client/cli"
	"github.com/likecoin/likechain/x/likenft/testutil"
	"github.com/likecoin/likechain/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithClassesByISCNObjects(t *testing.T, n int) (*network.Network, []types.ClassesByISCN) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		classesByISCN := types.ClassesByISCN{
			IscnIdPrefix: strconv.Itoa(i),
		}
		nullify.Fill(&classesByISCN)
		state.ClassesByISCNList = append(state.ClassesByISCNList, classesByISCN)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.ClassesByISCNList
}

func TestShowClassesByISCN(t *testing.T) {
	net, _objs := networkWithClassesByISCNObjects(t, 2)
	objs := testutil.BatchDummyConcretizeClassesByISCN(_objs)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc           string
		idIscnIdPrefix string

		args []string
		err  error
		obj  types.ConcreteClassesByISCN
	}{
		{
			desc:           "found",
			idIscnIdPrefix: objs[0].IscnIdPrefix,

			args: common,
			obj:  objs[0],
		},
		{
			desc:           "not found",
			idIscnIdPrefix: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.InvalidArgument, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idIscnIdPrefix,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowClassesByISCN(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryClassesByISCNResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.ClassesByISCN)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.ClassesByISCN),
				)
			}
		})
	}
}

func TestListClassesByISCN(t *testing.T) {
	net, objs := networkWithClassesByISCNObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClassesByISCN(), args)
			require.NoError(t, err)
			var resp types.QueryClassesByISCNIndexResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.ClassesByISCN), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.ClassesByISCN),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClassesByISCN(), args)
			require.NoError(t, err)
			var resp types.QueryClassesByISCNIndexResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.ClassesByISCN), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.ClassesByISCN),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClassesByISCN(), args)
		require.NoError(t, err)
		var resp types.QueryClassesByISCNIndexResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.ClassesByISCN),
		)
	})
}
