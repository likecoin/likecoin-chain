package cli_test

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/likecoin/likecoin-chain/v4/testutil/network"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/client/cli"
	"github.com/likecoin/likecoin-chain/v4/x/likenft/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = strconv.IntSize

type AccountByClass struct {
	Address string
	ClassId string
}

func networkWithAccountByClass(t *testing.T, n int) (*network.Network, []AccountByClass) {
	t.Helper()
	cfg := network.DefaultConfig()
	likenftState := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &likenftState))
	nftState := nfttypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[nfttypes.ModuleName], &nftState))

	accountByClassList := []AccountByClass{}

	for i := 0; i < n; i++ {
		// Create random address
		pubBz := make([]byte, ed25519.PubKeySize)
		rand.Read(pubBz)
		pub := &ed25519.PubKey{Key: pubBz}
		address, _ := sdk.Bech32ifyAddressBytes("cosmos", pub.Address())

		classId := fmt.Sprintf("likenft1%s", strconv.Itoa(i))

		classData := types.ClassData{
			Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
			Parent: types.ClassParent{
				Type:    types.ClassParentType_ACCOUNT,
				Account: address,
			},
			Config: types.ClassConfig{
				Burnable: false,
			},
		}
		classDataInAny, _ := cdctypes.NewAnyWithValue(&classData)
		nftState.Classes = append(nftState.Classes, &nfttypes.Class{
			Id:          classId,
			Name:        "Class Name",
			Symbol:      "ABC",
			Description: "Testing Class 123",
			Uri:         "ipfs://abcdef",
			UriHash:     "abcdef",
			Data:        classDataInAny,
		})

		classesByAccount := types.ClassesByAccount{
			Account:  address,
			ClassIds: []string{classId},
		}
		likenftState.ClassesByAccountList = append(likenftState.ClassesByAccountList, classesByAccount)

		accountByClass := AccountByClass{
			ClassId: classId,
			Address: address,
		}
		accountByClassList = append(accountByClassList, accountByClass)
	}

	for _, asd := range likenftState.ClassesByAccountList {
		fmt.Println(asd.Account, asd.ClassIds)
	}

	nftBuf, err := cfg.Codec.MarshalJSON(&nftState)
	require.NoError(t, err)
	cfg.GenesisState[nfttypes.ModuleName] = nftBuf

	likenftBuf, err := cfg.Codec.MarshalJSON(&likenftState)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = likenftBuf

	return network.New(t, cfg), accountByClassList
}

func TestShowAccountByClass(t *testing.T) {
	net, objs := networkWithAccountByClass(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}

	for _, tc := range []struct {
		desc      string
		idClassId string

		args []string
		err  error
		obj  AccountByClass
	}{
		{
			desc:      "found",
			idClassId: objs[0].ClassId,
			args:      common,
			obj:       objs[0],
		},
		{
			desc:      "found",
			idClassId: objs[1].ClassId,
			args:      common,
			obj:       objs[1],
		},
		{
			desc:      "not found",
			idClassId: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.InvalidArgument, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idClassId,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdAccountByClass(), args)

			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryAccountByClassResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.obj.Address, resp.Address)
			}
		})
	}
}
