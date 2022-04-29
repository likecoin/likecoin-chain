package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	nfttypes "github.com/likecoin/likechain/backport/cosmos-sdk/v0.46.0-alpha2/x/nft"
	"github.com/likecoin/likechain/testutil/network"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/likecoin/likechain/x/likenft/client/cli"
	"github.com/likecoin/likechain/x/likenft/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = strconv.IntSize

type ISCNByClass struct {
	ClassId      string
	IscnIdPrefix string
}

func networkWithISCNByClass(t *testing.T, n int) (*network.Network, []ISCNByClass) {
	t.Helper()
	cfg := network.DefaultConfig()
	likenftState := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &likenftState))
	nftState := nfttypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[nfttypes.ModuleName], &nftState))
	iscnState := iscntypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[iscntypes.ModuleName], &iscnState))
	iscnByClassList := []ISCNByClass{}

	for i := 0; i < n; i++ {
		iscnIdPrefix := fmt.Sprintf("iscn://likecoin-chain/%s", strconv.Itoa(i))
		classId := fmt.Sprintf("likenft1%s", strconv.Itoa(i))

		classData := types.ClassData{
			Metadata: types.JsonInput(`{"aaaa": "bbbb"}`),
			Parent: types.ClassParent{
				Type:         types.ClassParentType_ISCN,
				IscnIdPrefix: iscnIdPrefix,
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

		ownerAddressBytes := []byte{0, 1, 0, 1, 0, 1, 0, 1}
		ownerAddress, _ := sdk.Bech32ifyAddressBytes("like", ownerAddressBytes)
		iscnState.ContentIdRecords = append(iscnState.ContentIdRecords, iscntypes.GenesisState_ContentIdRecord{
			IscnId:        iscnIdPrefix,
			Owner:         ownerAddress,
			LatestVersion: 1,
		})

		iscnRecord := fmt.Sprintf(`
		{
			"@id": "%s/1",
			"recordVersion": 1,
			"contentFingerprints": [
				"ipfs://1"
			]
		}
		`, iscnIdPrefix)

		iscnState.IscnRecords = append(iscnState.IscnRecords, iscntypes.IscnInput(iscnRecord))

		classesByIscn := types.ClassesByISCN{
			IscnIdPrefix: iscnIdPrefix,
			ClassIds:     []string{classId},
		}
		likenftState.ClassesByIscnList = append(likenftState.ClassesByIscnList, classesByIscn)

		iscnByClass := ISCNByClass{
			ClassId:      classId,
			IscnIdPrefix: iscnIdPrefix,
		}
		iscnByClassList = append(iscnByClassList, iscnByClass)
	}

	nftBuf, err := cfg.Codec.MarshalJSON(&nftState)
	require.NoError(t, err)
	cfg.GenesisState[nfttypes.ModuleName] = nftBuf

	iscnBuf, err := cfg.Codec.MarshalJSON(&iscnState)
	require.NoError(t, err)
	cfg.GenesisState[iscntypes.ModuleName] = iscnBuf

	likenftBuf, err := cfg.Codec.MarshalJSON(&likenftState)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = likenftBuf

	return network.New(t, cfg), iscnByClassList
}

func TestShowISCNByClass(t *testing.T) {
	net, objs := networkWithISCNByClass(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}

	for _, tc := range []struct {
		desc      string
		idClassId string

		args []string
		err  error
		obj  ISCNByClass
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdISCNByClass(), args)

			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryISCNByClassResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.obj.IscnIdPrefix, resp.IscnIdPrefix)
			}
		})
	}
}
