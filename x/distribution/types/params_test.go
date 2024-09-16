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
		{"success", fields{toDec("0.1"), toDec("0"), toDec("0"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, false},
		{"negative community tax", fields{toDec("-0.1"), toDec("0"), toDec("0"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, true},
		{"negative mca tax", fields{toDec("0.1"), toDec("0"), toDec("-0.1"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, true},
		{"negative base proposer reward (must not matter)", fields{toDec("0.1"), toDec("0"), toDec("-0.1"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, false},
		{"negative bonus proposer reward (must not matter)", fields{toDec("0.1"), toDec("0"), toDec("-0.1"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, false},
		{"total sum greater than 1 (must not matter)", fields{toDec("0.2"), toDec("0.5"), toDec("0.4"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, false},
		{"community tax greater than 1", fields{toDec("1.1"), toDec("0"), toDec("0"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, true},
		{"community tax nil", fields{sdkmath.LegacyDec{}, toDec("0"), toDec("0"), false, toDec("0.1"), "cosmos15m77x4pe6w9vtpuqm22qxu0ds7vn4ehz9dd9u2"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := types.Params{
				CommunityTax:        tt.fields.CommunityTax,
				WithdrawAddrEnabled: tt.fields.WithdrawAddrEnabled,
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
