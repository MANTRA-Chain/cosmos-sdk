package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens performs reward and fee distribution to all validators based
// on the F1 fee distribution specification, with MCA tax deducted before community tax.
func (k Keeper) AllocateTokens(ctx context.Context, totalPreviousPower int64, bondedVotes []abci.VoteInfo) error {
	// Fetch and clear the collected fees for distribution.
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// Transfer collected fees to the distribution module account.
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		return err
	}

	// Get the current fee pool.
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return err
	}

	// Deduct MCA tax first.
	mcaTaxRate, err := k.GetMcaTax(ctx)
	if err != nil {
		return err
	}
	mcaTaxAmount := feesCollected.MulDecTruncate(mcaTaxRate)
	remainingAfterMca := feesCollected.Sub(mcaTaxAmount)

	// Get the MCA address from the params.
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	mcaAddressStr := params.McaAddress

	mcaAddress, err := sdk.AccAddressFromBech32(mcaAddressStr)
	if err != nil {
		return err
	}
	mcaTaxCoins, leftoverMcaTax := mcaTaxAmount.TruncateDecimal()
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mcaAddress, mcaTaxCoins)
	if err != nil {
		return err
	}

	// Deduct community tax from the remaining funds after MCA tax.
	communityTaxRate, err := k.GetCommunityTax(ctx)
	if err != nil {
		return err
	}
	communityTaxAmount := remainingAfterMca.MulDecTruncate(communityTaxRate)
	remainingAfterCommunityTax := remainingAfterMca.Sub(communityTaxAmount)

	// Add the community tax to the community pool.
	feePool.CommunityPool = feePool.CommunityPool.Add(communityTaxAmount...)

	// If there is no validator power, add the remaining funds to the community pool and return.
	if totalPreviousPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(remainingAfterCommunityTax...)
		return k.FeePool.Set(ctx, feePool)
	}

	// Distribute the remaining tokens to validators based on voting power.
	totalDistributed := sdk.NewDecCoins()
	for _, vote := range bondedVotes {
		validator, err := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		if err != nil {
			return err
		}

		// Calculate the fraction of power for each validator.
		powerFraction := math.LegacyNewDec(vote.Validator.Power).QuoTruncate(math.LegacyNewDec(totalPreviousPower))
		reward := remainingAfterCommunityTax.MulDecTruncate(powerFraction)

		// Allocate tokens to the validator.
		err = k.AllocateTokensToValidator(ctx, validator, reward)
		if err != nil {
			return err
		}

		totalDistributed = totalDistributed.Add(reward...)
	}

	// Any leftover due to rounding is added to the community pool.
	leftover := remainingAfterCommunityTax.Sub(totalDistributed)
	if !leftover.IsZero() {
		feePool.CommunityPool = feePool.CommunityPool.Add(leftover...)
	}

	// Add leftover MCA tax to the community pool.
	if !leftoverMcaTax.IsZero() {
		feePool.CommunityPool = feePool.CommunityPool.Add(leftoverMcaTax...)
	}

	// Update the fee pool.
	return k.FeePool.Set(ctx, feePool)
}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx context.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) error {
	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	if err != nil {
		return err
	}

	// update current commission
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator()),
		),
	)
	currentCommission, err := k.GetValidatorAccumulatedCommission(ctx, valBz)
	if err != nil {
		return err
	}

	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	err = k.SetValidatorAccumulatedCommission(ctx, valBz, currentCommission)
	if err != nil {
		return err
	}

	// update current rewards
	currentRewards, err := k.GetValidatorCurrentRewards(ctx, valBz)
	if err != nil {
		return err
	}

	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	err = k.SetValidatorCurrentRewards(ctx, valBz, currentRewards)
	if err != nil {
		return err
	}

	// update outstanding rewards
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator()),
		),
	)

	outstanding, err := k.GetValidatorOutstandingRewards(ctx, valBz)
	if err != nil {
		return err
	}

	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	return k.SetValidatorOutstandingRewards(ctx, valBz, outstanding)
}
