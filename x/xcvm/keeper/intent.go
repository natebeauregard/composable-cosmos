package keeper

import (
	"encoding/binary"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibccore "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

func (k Keeper) SendValidatedTransferIntent(ctx sdk.Context, msg *types.MsgSendTransferIntent) error {
	clientId := msg.ClientId
	clientState, found := k.clientKeeper.GetClientState(ctx, clientId)
	if !found {
		return types.ErrClientNotFound
	}

	clientStatus := k.clientKeeper.GetClientStatus(ctx, clientState, clientId)
	if clientStatus != ibccore.Active {
		return types.ErrClientNotActive
	}

	intentId := k.GetNextIntentId(ctx)
	transferIntent := types.TransferIntent{
		ClientId:           clientId,
		SourceAddress:      msg.FromAddress,
		DestinationAddress: msg.DestinationAddress,
		Amount:             msg.Amount,
	}
	k.AddTransferIntent(ctx, transferIntent, intentId)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventAddTransferIntent,
		sdk.NewAttribute(types.AttributeKeyIntentId, strconv.FormatUint(intentId, 10)),
		sdk.NewAttribute(types.AttributeKeyClientId, transferIntent.ClientId),
		sdk.NewAttribute(types.AttributeKeySourceAddress, transferIntent.SourceAddress),
		sdk.NewAttribute(types.AttributeKeyDestinationAddress, transferIntent.DestinationAddress),
		sdk.NewAttribute(types.AttributeKeyAmount, transferIntent.Amount.String()),
	))

	k.SetNextIntentId(ctx, intentId+1)

	return nil
}

func (k Keeper) GetNextIntentId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	intentIdBz := store.Get(types.TransferIntentIdKey)
	intentId := binary.BigEndian.Uint64(intentIdBz)

	return intentId
}

func (k Keeper) SetNextIntentId(ctx sdk.Context, intentId uint64) {
	store := ctx.KVStore(k.storeKey)

	intentIdBz := make([]byte, 8)
	binary.BigEndian.PutUint64(intentIdBz, intentId)

	store.Set(types.TransferIntentIdKey, intentIdBz)
}

// Stores an intent object in the store
func (k Keeper) AddTransferIntent(ctx sdk.Context, transferIntent types.TransferIntent, intentId uint64) {
	store := ctx.KVStore(k.storeKey)

	intentIdBz := make([]byte, 8)
	binary.BigEndian.PutUint64(intentIdBz, intentId)

	transferIntentKey := types.GetPendingTransferIntentKeyById(intentIdBz)
	transferIntentValue := k.cdc.MustMarshal(&transferIntent)

	store.Set(transferIntentKey, transferIntentValue)
}

func (k Keeper) VerifyIntentProof(ctx sdk.Context, msg *types.MsgVerifyTransferIntentProof) error {
	// TODO implement verifying intent proofs
	return nil
}
