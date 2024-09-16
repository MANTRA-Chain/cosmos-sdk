package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultParams returns default distribution parameters
func DefaultParams() Params {
	return Params{
		CommunityTax:        math.LegacyNewDecWithPrec(2, 2), // 2%
		BaseProposerReward:  math.LegacyZeroDec(),            // deprecated
		BonusProposerReward: math.LegacyZeroDec(),            // deprecated
		WithdrawAddrEnabled: true,
		McaTax:              math.LegacyNewDecWithPrec(40, 2),                // 40%
		McaAddress:          "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls", // TODO: change this to the actual mca address
	}
}

// ValidateBasic performs basic validation on distribution parameters.
func (p Params) ValidateBasic() error {
	err := validateCommunityTax(p.CommunityTax)
	if err != nil {
		return err
	}
	err = validateMcaTax(p.McaTax)
	if err != nil {
		return err
	}
	err = validateMcaAddress(p.McaAddress)
	if err != nil {
		return err
	}
	return nil
}

func validateCommunityTax(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("community tax must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("community tax must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("community tax too large: %s", v)
	}

	return nil
}

func validateMcaTax(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("mca tax must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("mca tax must be positive: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("mca tax too large: %s", v)
	}

	return nil
}

func validateMcaAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("invalid MCA address: %s", err)
	}

	return nil
}

func validateWithdrawAddrEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
