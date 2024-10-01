package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MultiBankHooks combine multiple bank hooks, all hook functions are run in array sequence
type MultiBankHooks struct {
	bankhooks []BankHooks
}

// NewMultiBankHooks takes a list of BankHooks and returns a MultiBankHooks
func NewMultiBankHooks(hooks ...BankHooks) *MultiBankHooks {
	return &MultiBankHooks{
		bankhooks: hooks,
	}
}

// TrackBeforeSend runs the TrackBeforeSend hooks in order for each BankHook in a MultiBankHooks struct
func (h MultiBankHooks) TrackBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	hooks := h.bankhooks
	for i := range hooks {
		hooks[i].TrackBeforeSend(ctx, from, to, amount)
	}
}

// BlockBeforeSend runs the BlockBeforeSend hooks in order for each BankHook in a MultiBankHooks struct
func (h MultiBankHooks) BlockBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	hooks := h.bankhooks
	for i := range hooks {
		err := hooks[i].BlockBeforeSend(ctx, from, to, amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MultiBankHooks) Append(hook BankHooks) {
	r.bankhooks = append(r.bankhooks, hook)
}
