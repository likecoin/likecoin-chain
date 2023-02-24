package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

var (
	addr1 = "like1ukmjl5s6pnw2txkvz2hd2n0f6dulw34h9rw5zn"

	fingerprint1 = "hash://sha256/9564b85669d5e96ac969dd0161b8475bbced9e5999c6ec598da718a3045d6f2e"

	stakeholder1 = IscnInput(`
{
	"entity": {
		"@id": "did:cosmos:5sy29r37gfxvxz21rh4r0ktpuc46pzjrmz29g45",
		"name": "Chung Wu"
	},
	"rewardProportion": 95,
	"contributionType": "http://schema.org/author"
}`)

	stakeholder2 = IscnInput(`
{
	"rewardProportion": 5,
	"contributionType": "http://schema.org/citation",
	"footprint": "https://en.wikipedia.org/wiki/Fibonacci_number",
	"description": "The blog post referred the matrix form of computing Fibonacci numbers."
}`)

	contentMetadata1 = IscnInput(`
{
	"@context": "http://schema.org/",
	"@type": "CreativeWorks",
	"title": "使用矩陣計算遞歸關係式",
	"description": "An article on computing recursive function with matrix multiplication.",
	"datePublished": "2019-04-19",
	"version": 1,
	"url": "https://nnkken.github.io/post/recursive-relation/",
	"author": "https://github.com/nnkken",
	"usageInfo": "https://creativecommons.org/licenses/by/4.0",
	"keywords": "matrix,recursion"
}`)

	createIscnRecord1 = IscnRecord{
		ContentFingerprints: []string{fingerprint1},
		Stakeholders:        []IscnInput{stakeholder1, stakeholder2},
		ContentMetadata:     contentMetadata1,
	}

	createIscnRecord2 = IscnRecord{
		ContentMetadata: IscnInput(`{}`),
	}

	createIscnRecord3 = IscnRecord{
		RecordNotes:         "",
		ContentFingerprints: []string{},
		Stakeholders:        []IscnInput{},
		ContentMetadata:     IscnInput(`{}`),
	}

	createIscnRecord4 = IscnRecord{
		ContentMetadata: IscnInput(`{
			"emptyField": "",
			"emptyArray": [],
			"emptyObject": {}
		}`),
	}

	msgCreateIscnRecord1NoNonce1 = MsgCreateIscnRecord{
		From:   addr1,
		Record: createIscnRecord2,
	}
	msgCreateIscnRecord1NoNonce2 = MsgCreateIscnRecord{
		From:   addr1,
		Record: createIscnRecord3,
	}
	msgCreateIscnRecord1NoNonce3 = MsgCreateIscnRecord{
		From:   addr1,
		Record: createIscnRecord4,
	}
	msgCreateIscnRecord1Nonce1 = MsgCreateIscnRecord{
		From:   addr1,
		Record: createIscnRecord2,
		Nonce:  1,
	}

	msgCreateIscnRecordBytesNoNonce1 = []byte(`{"type":"likecoin-chain/MsgCreateIscnRecord","value":{"from":"like1ukmjl5s6pnw2txkvz2hd2n0f6dulw34h9rw5zn","record":{"contentMetadata":{}}}}`)
	msgCreateIscnRecordBytesNoNonce2 = []byte(`{"type":"likecoin-chain/MsgCreateIscnRecord","value":{"from":"like1ukmjl5s6pnw2txkvz2hd2n0f6dulw34h9rw5zn","record":{"contentMetadata":{"emptyArray":[],"emptyField":"","emptyObject":{}}}}}`)
	msgCreateIscnRecordBytesNonce1   = []byte(`{"type":"likecoin-chain/MsgCreateIscnRecord","value":{"from":"like1ukmjl5s6pnw2txkvz2hd2n0f6dulw34h9rw5zn","nonce":"1","record":{"contentMetadata":{}}}}`)

	registryName = "likecoin-chain"

	iscnIdNoNonce = IscnId{
		Prefix: IscnIdPrefix{
			RegistryName: registryName,
			ContentId:    "rv5ahVKmxSu93jZlO6-X0oHP5NoIk0uQJj0zg84qKqs",
		},
		Version: 1,
	}
	iscnIdNonce1 = IscnId{
		Prefix: IscnIdPrefix{
			RegistryName: registryName,
			ContentId:    "5e1QGL5xM8GUFNU87poFMOQHcMATyqvPiGKIUQduuKw",
		},
		Version: 1,
	}
)

func SetAddressPrefixes() {
	bech32PrefixesAccAddr := []string{"like", "cosmos"}
	bech32PrefixesAccPub := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesValAddr := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesValPub := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesConsAddr := make([]string, 0, len(bech32PrefixesAccAddr))
	bech32PrefixesConsPub := make([]string, 0, len(bech32PrefixesAccAddr))

	for _, prefix := range bech32PrefixesAccAddr {
		bech32PrefixesAccPub = append(bech32PrefixesAccPub, prefix+"pub")
		bech32PrefixesValAddr = append(bech32PrefixesValAddr, prefix+"valoper")
		bech32PrefixesValPub = append(bech32PrefixesValPub, prefix+"valoperpub")
		bech32PrefixesConsAddr = append(bech32PrefixesConsAddr, prefix+"valcons")
		bech32PrefixesConsPub = append(bech32PrefixesConsPub, prefix+"valconspub")
	}
	config := sdk.GetConfig()
	config.SetBech32PrefixesForAccount(bech32PrefixesAccAddr, bech32PrefixesAccPub)
	config.SetBech32PrefixesForValidator(bech32PrefixesValAddr, bech32PrefixesValPub)
	config.SetBech32PrefixesForConsensusNode(bech32PrefixesConsAddr, bech32PrefixesConsPub)
}

func TestMsgCreateIscnRecord_ValidateBasic(t *testing.T) {
	SetAddressPrefixes()
	tests := []struct {
		name string
		msg  MsgCreateIscnRecord
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateIscnRecord{
				From:   "invalid_address",
				Record: createIscnRecord1,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address and no nonce",
			msg: MsgCreateIscnRecord{
				From:   addr1,
				Record: createIscnRecord1,
			},
		}, {
			name: "valid address and has nonce = 0",
			msg: MsgCreateIscnRecord{
				From:   addr1,
				Record: createIscnRecord1,
				Nonce:  0,
			},
		}, {
			name: "valid address and has nonce = 1",
			msg: MsgCreateIscnRecord{
				From:   addr1,
				Record: createIscnRecord1,
				Nonce:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgCreateIscnRecord_GetSignBytes(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateIscnRecord
		want []byte
	}{
		{
			name: "valid with no nonce",
			msg:  msgCreateIscnRecord1NoNonce1,
			want: msgCreateIscnRecordBytesNoNonce1,
		},
		{
			name: "valid with no nonce and ignore empty optional fields in IscnRecord",
			msg:  msgCreateIscnRecord1NoNonce2,
			want: msgCreateIscnRecordBytesNoNonce1,
		},
		{
			name: "valid with no nonce and keep empty fields in ContentMetadata",
			msg:  msgCreateIscnRecord1NoNonce3,
			want: msgCreateIscnRecordBytesNoNonce2,
		},
		{
			name: "valid with assigned nonce",
			msg:  msgCreateIscnRecord1Nonce1,
			want: msgCreateIscnRecordBytesNonce1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.GetSignBytes(); string(got) != string(tt.want) {
				t.Errorf("MsgCreateIscnRecord.GetSignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgCreateIscnRecord_GenerateNewIscnIdWithSeed(t *testing.T) {
	tests := []struct {
		name string
		seed []byte
		want IscnId
	}{
		{
			name: "valid with no nonce",
			seed: msgCreateIscnRecordBytesNoNonce1,
			want: iscnIdNoNonce,
		},
		{
			name: "valid with assigned nonce",
			seed: msgCreateIscnRecordBytesNonce1,
			want: iscnIdNonce1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateNewIscnIdWithSeed(registryName, tt.seed); got != tt.want {
				t.Errorf("MsgCreateIscnRecord.GenerateNewIscnIdWithSeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
