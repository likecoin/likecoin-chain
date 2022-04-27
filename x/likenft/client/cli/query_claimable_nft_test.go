package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/testutil/network"
	"github.com/likecoin/likechain/testutil/nullify"
	"github.com/likecoin/likechain/x/likenft/client/cli"
	"github.com/likecoin/likechain/x/likenft/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithMintableNFTObjects(t *testing.T, n int) (*network.Network, []types.MintableNFT) {
	t.Helper()
	cfg := network.DefaultConfig()
	// seed nft class
	nftState := nft.GenesisState{}
	// seed mintable
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		// seed nft class
		classData := types.ClassData{
			MintableCount: 0,
		}
		classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
		class := nft.Class{
			Id:   strconv.Itoa(i),
			Data: classDataInAny,
		}
		nftState.Classes = append(nftState.Classes, &class)
		// seed mintable
		mintableNFT := types.MintableNFT{
			ClassId: strconv.Itoa(i),
			Id:      strconv.Itoa(i),
		}
		state.MintableNFTList = append(state.MintableNFTList, mintableNFT)
	}
	// seed nft class
	nftBuf, err := cfg.Codec.MarshalJSON(&nftState)
	require.NoError(t, err)
	cfg.GenesisState[nft.ModuleName] = nftBuf
	// seed mintable
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.MintableNFTList
}

func TestShowMintableNFT(t *testing.T) {
	net, objs := networkWithMintableNFTObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc      string
		idClassId string
		idId      string

		args []string
		err  error
		obj  types.MintableNFT
	}{
		{
			desc:      "found",
			idClassId: objs[0].ClassId,
			idId:      objs[0].Id,

			args: common,
			obj:  objs[0],
		},
		{
			desc:      "not found",
			idClassId: strconv.Itoa(100000),
			idId:      strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.InvalidArgument, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idClassId,
				tc.idId,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowMintableNFT(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryMintableNFTResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.MintableNFT)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.MintableNFT),
				)
			}
		})
	}
}

func TestListMintableNFT(t *testing.T) {
	net, objs := networkWithMintableNFTObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListMintableNFT(), args)
			require.NoError(t, err)
			var resp types.QueryMintableNFTIndexResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.MintableNFT), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.MintableNFT),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListMintableNFT(), args)
			require.NoError(t, err)
			var resp types.QueryMintableNFTIndexResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.MintableNFT), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.MintableNFT),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListMintableNFT(), args)
		require.NoError(t, err)
		var resp types.QueryMintableNFTIndexResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.MintableNFT),
		)
	})
}
