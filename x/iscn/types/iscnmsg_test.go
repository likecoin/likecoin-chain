package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

var (
	// priv1 = secp256k1.GenPrivKey()
	// addr1 = sdk.AccAddress(priv1.PubKey().Address())
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

	msgCreateIscnRecord1 = MsgCreateIscnRecord{
		From: addr1,
		Record: IscnRecord{
			ContentFingerprints: []string{fingerprint1},
			Stakeholders:        []IscnInput{stakeholder1, stakeholder2},
			ContentMetadata:     contentMetadata1,
		},
	}

	msgCreateIscnRecordBytes1 = []byte{
		123, 34, 116, 121, 112, 101, 34, 58, 34, 108, 105, 107, 101, 99, 111, 105, 110, 45, 99, 104, 97, 105, 110, 47, 77, 115, 103, 67, 114, 101, 97, 116, 101, 73, 115, 99, 110, 82, 101, 99, 111, 114, 100, 34, 44, 34, 118, 97, 108, 117, 101, 34, 58, 123, 34, 102, 114, 111, 109, 34, 58, 34, 99, 111, 115, 109, 111, 115, 49, 121, 53, 52, 101, 120, 109, 120, 56, 52, 99, 113, 116, 97, 115, 118, 106, 110, 115, 107, 102, 57, 102, 54, 51, 100, 106, 117, 117, 106, 54, 56, 112, 55, 104, 113, 102, 52, 55, 34, 44, 34, 114, 101, 99, 111, 114, 100, 34, 58, 123, 34, 99, 111, 110, 116, 101, 110, 116, 70, 105, 110, 103, 101, 114, 112, 114, 105, 110, 116, 115, 34, 58, 91, 34, 104, 97, 115, 104, 58, 47, 47, 115, 104, 97, 50, 53, 54, 47, 57, 53, 54, 52, 98, 56, 53, 54, 54, 57, 100, 53, 101, 57, 54, 97, 99, 57, 54, 57, 100, 100, 48, 49, 54, 49, 98, 56, 52, 55, 53, 98, 98, 99, 101, 100, 57, 101, 53, 57, 57, 57, 99, 54, 101, 99, 53, 57, 56, 100, 97, 55, 49, 56, 97, 51, 48, 52, 53, 100, 54, 102, 50, 101, 34, 93, 44, 34, 99, 111, 110, 116, 101, 110, 116, 77, 101, 116, 97, 100, 97, 116, 97, 34, 58, 123, 34, 64, 99, 111, 110, 116, 101, 120, 116, 34, 58, 34, 104, 116, 116, 112, 58, 47, 47, 115, 99, 104, 101, 109, 97, 46, 111, 114, 103, 47, 34, 44, 34, 64, 116, 121, 112, 101, 34, 58, 34, 67, 114, 101, 97, 116, 105, 118, 101, 87, 111, 114, 107, 115, 34, 44, 34, 97, 117, 116, 104, 111, 114, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 110, 110, 107, 107, 101, 110, 34, 44, 34, 100, 97, 116, 101, 80, 117, 98, 108, 105, 115, 104, 101, 100, 34, 58, 34, 50, 48, 49, 57, 45, 48, 52, 45, 49, 57, 34, 44, 34, 100, 101, 115, 99, 114, 105, 112, 116, 105, 111, 110, 34, 58, 34, 65, 110, 32, 97, 114, 116, 105, 99, 108, 101, 32, 111, 110, 32, 99, 111, 109, 112, 117, 116, 105, 110, 103, 32, 114, 101, 99, 117, 114, 115, 105, 118, 101, 32, 102, 117, 110, 99, 116, 105, 111, 110, 32, 119, 105, 116, 104, 32, 109, 97, 116, 114, 105, 120, 32, 109, 117, 108, 116, 105, 112, 108, 105, 99, 97, 116, 105, 111, 110, 46, 34, 44, 34, 107, 101, 121, 119, 111, 114, 100, 115, 34, 58, 34, 109, 97, 116, 114, 105, 120, 44, 114, 101, 99, 117, 114, 115, 105, 111, 110, 34, 44, 34, 116, 105, 116, 108, 101, 34, 58, 34, 228, 189, 191, 231, 148, 168, 231, 159, 169, 233, 153, 163, 232, 168, 136, 231, 174, 151, 233, 129, 158, 230, 173, 184, 233, 151, 156, 228, 191, 130, 229, 188, 143, 34, 44, 34, 117, 114, 108, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 110, 110, 107, 107, 101, 110, 46, 103, 105, 116, 104, 117, 98, 46, 105, 111, 47, 112, 111, 115, 116, 47, 114, 101, 99, 117, 114, 115, 105, 118, 101, 45, 114, 101, 108, 97, 116, 105, 111, 110, 47, 34, 44, 34, 117, 115, 97, 103, 101, 73, 110, 102, 111, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 99, 114, 101, 97, 116, 105, 118, 101, 99, 111, 109, 109, 111, 110, 115, 46, 111, 114, 103, 47, 108, 105, 99, 101, 110, 115, 101, 115, 47, 98, 121, 47, 52, 46, 48, 34, 44, 34, 118, 101, 114, 115, 105, 111, 110, 34, 58, 49, 125, 44, 34, 115, 116, 97, 107, 101, 104, 111, 108, 100, 101, 114, 115, 34, 58, 91, 123, 34, 99, 111, 110, 116, 114, 105, 98, 117, 116, 105, 111, 110, 84, 121, 112, 101, 34, 58, 34, 104, 116, 116, 112, 58, 47, 47, 115, 99, 104, 101, 109, 97, 46, 111, 114, 103, 47, 97, 117, 116, 104, 111, 114, 34, 44, 34, 101, 110, 116, 105, 116, 121, 34, 58, 123, 34, 64, 105, 100, 34, 58, 34, 100, 105, 100, 58, 99, 111, 115, 109, 111, 115, 58, 53, 115, 121, 50, 57, 114, 51, 55, 103, 102, 120, 118, 120, 122, 50, 49, 114, 104, 52, 114, 48, 107, 116, 112, 117, 99, 52, 54, 112, 122, 106, 114, 109, 122, 50, 57, 103, 52, 53, 34, 44, 34, 110, 97, 109, 101, 34, 58, 34, 67, 104, 117, 110, 103, 32, 87, 117, 34, 125, 44, 34, 114, 101, 119, 97, 114, 100, 80, 114, 111, 112, 111, 114, 116, 105, 111, 110, 34, 58, 57, 53, 125, 44, 123, 34, 99, 111, 110, 116, 114, 105, 98, 117, 116, 105, 111, 110, 84, 121, 112, 101, 34, 58, 34, 104, 116, 116, 112, 58, 47, 47, 115, 99, 104, 101, 109, 97, 46, 111, 114, 103, 47, 99, 105, 116, 97, 116, 105, 111, 110, 34, 44, 34, 100, 101, 115, 99, 114, 105, 112, 116, 105, 111, 110, 34, 58, 34, 84, 104, 101, 32, 98, 108, 111, 103, 32, 112, 111, 115, 116, 32, 114, 101, 102, 101, 114, 114, 101, 100, 32, 116, 104, 101, 32, 109, 97, 116, 114, 105, 120, 32, 102, 111, 114, 109, 32, 111, 102, 32, 99, 111, 109, 112, 117, 116, 105, 110, 103, 32, 70, 105, 98, 111, 110, 97, 99, 99, 105, 32, 110, 117, 109, 98, 101, 114, 115, 46, 34, 44, 34, 102, 111, 111, 116, 112, 114, 105, 110, 116, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 101, 110, 46, 119, 105, 107, 105, 112, 101, 100, 105, 97, 46, 111, 114, 103, 47, 119, 105, 107, 105, 47, 70, 105, 98, 111, 110, 97, 99, 99, 105, 95, 110, 117, 109, 98, 101, 114, 34, 44, 34, 114, 101, 119, 97, 114, 100, 80, 114, 111, 112, 111, 114, 116, 105, 111, 110, 34, 58, 53, 125, 93, 125, 125, 125,
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
				From: "invalid_address",
				Record: IscnRecord{
					ContentMetadata: contentMetadata1,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateIscnRecord{
				From: addr1,
				Record: IscnRecord{
					ContentMetadata: contentMetadata1,
				},
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
			name: "valid",
			msg: MsgCreateIscnRecord{
				From: addr1,
				Record: IscnRecord{
					ContentFingerprints: []string{fingerprint1},
					Stakeholders:        []IscnInput{stakeholder1, stakeholder2},
					ContentMetadata:     contentMetadata1,
				},
			},
			want: msgCreateIscnRecordBytes1,
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
