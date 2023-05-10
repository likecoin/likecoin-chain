package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likecoin-chain/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateListing_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateListing
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateListing{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateListing{
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

func TestMsgUpdateListing_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateListing
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateListing{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateListing{
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

func TestMsgDeleteListing_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteListing
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteListing{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteListing{
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
