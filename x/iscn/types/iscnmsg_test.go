package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

var (
	addr1 = "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"

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

	msgCreateIscnRecord1NoNonce = MsgCreateIscnRecord{
		From:   addr1,
		Record: createIscnRecord2,
	}
	msgCreateIscnRecord1Nonce1 = MsgCreateIscnRecord{
		From:   addr1,
		Record: createIscnRecord2,
		Nonce:  1,
	}

	msgCreateIscnRecordBytesNoNonce = []byte(`{"type":"likecoin-chain/MsgCreateIscnRecord","value":{"from":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","record":{"contentMetadata":{}}}}`)
	msgCreateIscnRecordBytesNonce1  = []byte(`{"type":"likecoin-chain/MsgCreateIscnRecord","value":{"from":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","nonce":"1","record":{"contentMetadata":{}}}}`)

	registryName = "likecoin-chain"

	iscnIdNoNonce = IscnId{
		Prefix: IscnIdPrefix{
			RegistryName: registryName,
			ContentId:    "L0pDvwnj_9yt1ZajXpV_lsZf8niv-UQWADKoancbfAw",
		},
		Version: 1,
	}
	iscnIdNonce1 = IscnId{
		Prefix: IscnIdPrefix{
			RegistryName: registryName,
			ContentId:    "tKcMZedw5ktnw74K-1kyg5Tw8u4N3ZbE5mkRa1sOewo",
		},
		Version: 1,
	}
)

func TestMsgCreateIscnRecord_ValidateBasic(t *testing.T) {
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
			msg:  msgCreateIscnRecord1NoNonce,
			want: msgCreateIscnRecordBytesNoNonce,
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
			seed: msgCreateIscnRecordBytesNoNonce,
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
