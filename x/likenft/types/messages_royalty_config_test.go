package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likecoin-chain/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateRoyaltyConfig_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateRoyaltyConfig
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateRoyaltyConfig{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateRoyaltyConfig{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateRoyaltyConfig_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateRoyaltyConfig
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateRoyaltyConfig{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateRoyaltyConfig{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteRoyaltyConfig_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteRoyaltyConfig
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteRoyaltyConfig{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteRoyaltyConfig{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
