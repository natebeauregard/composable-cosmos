package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker of xcvm module.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// TODO: investigate event emission at the beginning of each block
}
