package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func TestParams_ValidateBasic(t *testing.T) {
	toDec := sdkmath.LegacyMustNewDecFromStr

	type fields struct {
		CommunityTax        sdkmath.LegacyDec
		BaseProposerReward  sdkmath.LegacyDec
		BonusProposerReward sdkmath.LegacyDec
		WithdrawAddrEnabled bool
		McaTax              sdkmath.LegacyDec
		McaAddress          string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				CommunityTax:        toDec("0.1"),
				BaseProposerReward:  toDec("0"),
				BonusProposerReward: toDec("0"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: false,
		},
		{
			name: "negative community tax",
			fields: fields{
				CommunityTax:        toDec("-0.1"),
				BaseProposerReward:  toDec("0"),
				BonusProposerReward: toDec("0"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: true,
		},
		{
			name: "negative mca tax",
			fields: fields{
				CommunityTax:        toDec("0.1"),
				BaseProposerReward:  toDec("0"),
				BonusProposerReward: toDec("-0.1"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: true,
		},
		{
			name: "negative base proposer reward (must not matter)",
			fields: fields{
				CommunityTax:        toDec("0.1"),
				BaseProposerReward:  toDec("0"),
				BonusProposerReward: toDec("-0.1"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: false,
		},
		{
			name: "negative bonus proposer reward (must not matter)",
			fields: fields{
				CommunityTax:        toDec("0.1"),
				BaseProposerReward:  toDec("0"),
				BonusProposerReward: toDec("-0.1"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: false,
		},
		{
			name: "total sum greater than 1 (must not matter)",
			fields: fields{
				CommunityTax:        toDec("0.2"),
				BaseProposerReward:  toDec("0.5"),
				BonusProposerReward: toDec("0.4"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: false,
		},
		{
			name: "community tax greater than 1",
			fields: fields{
				CommunityTax:        toDec("1.1"),
				BaseProposerReward:  toDec("0"),
				BonusProposerReward: toDec("0"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: true,
		},
		{
			name: "community tax nil",
			fields: fields{
				CommunityTax:        sdkmath.LegacyDec{},
				BaseProposerReward:  toDec("0"),
				BonusProposerReward: toDec("0"),
				WithdrawAddrEnabled: false,
				McaTax:              toDec("0.1"),
				McaAddress:          "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := types.Params{
				CommunityTax:        tt.fields.CommunityTax,
				WithdrawAddrEnabled: tt.fields.WithdrawAddrEnabled,
				McaTax:              tt.fields.McaTax,
				McaAddress:          tt.fields.McaAddress,
			}
			if err := p.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultParams(t *testing.T) {
	require.NoError(t, types.DefaultParams().ValidateBasic())
}
