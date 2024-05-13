package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

func (k Keeper) SendIntent(ctx sdk.Context, msg *types.MsgSendTransferIntent) error {
	// TODO implement sending intents
	return nil
}

func (k Keeper) VerifyIntentProof(ctx sdk.Context, msg *types.MsgVerifyTransferIntentProof) error {
	// TODO implement verifying intent proofs
	return nil
}
