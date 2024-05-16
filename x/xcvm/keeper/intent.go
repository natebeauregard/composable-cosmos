package keeper

import (
	"encoding/binary"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	ibccore "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

func (k Keeper) SendValidatedTransferIntent(ctx sdk.Context, msg *types.MsgSendTransferIntent) error {
	clientId := msg.ClientId

	_, err := k.ValidateClientId(ctx, clientId)
	if err != nil {
		return err
	}

	intentId := k.GetNextIntentId(ctx)
	transferIntent := types.TransferIntent{
		ClientId:           clientId,
		SourceAddress:      msg.FromAddress,
		DestinationAddress: msg.DestinationAddress,
		Amount:             msg.Amount,
	}
	k.AddTransferIntent(ctx, transferIntent, intentId)

	// TODO: post bounty for solver?
	// TODO: should a collateral amount be set here as well for the solver to deposit before executing the intent?

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

	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	transferIntentValue := k.cdc.MustMarshal(&transferIntent)

	store.Set(transferIntentKey, transferIntentValue)
}

func (k Keeper) VerifyIntentProof(ctx sdk.Context, msg *types.MsgVerifyTransferIntentProof) error {
	store := ctx.KVStore(k.storeKey)

	intentId := msg.IntentId
	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	if !store.Has(transferIntentKey) {
		return types.ErrInvalidIntentId
	}

	transferIntentBz := store.Get(transferIntentKey)
	var transferIntent types.TransferIntent
	err := k.cdc.Unmarshal(transferIntentBz, &transferIntent)
	if err != nil {
		return err
	}

	clientId := transferIntent.ClientId
	clientState, err := k.ValidateClientId(ctx, clientId)
	if err != nil {
		return err
	}

	height := clienttypes.NewHeight(msg.ProofHeight.RevisionNumber, msg.ProofHeight.RevisionHeight)
	merklePath := commitmenttypes.NewMerklePath(msg.MerklePath.KeyPath...)
	clientStore := k.clientKeeper.ClientStore(ctx, clientId)
	err = clientState.VerifyMembership(
		ctx,
		clientStore,
		k.cdc,
		height,
		msg.TimeDelay,
		msg.BlockDelay,
		msg.Proof,
		merklePath,
		msg.Value,
	)
	if err != nil {
		return err
	}

	// Purge resolved transfer intent after proof verification?
	// store.Delete(transferIntentKey)

	// TODO: unlock bounty for solver?

	return nil
}

func (k Keeper) ValidateClientId(ctx sdk.Context, clientId string) (ibccore.ClientState, error) {
	clientState, found := k.clientKeeper.GetClientState(ctx, clientId)
	if !found {
		return nil, types.ErrClientNotFound
	}

	clientStatus := k.clientKeeper.GetClientStatus(ctx, clientState, clientId)
	if clientStatus != ibccore.Active {
		return nil, types.ErrClientNotActive
	}

	return clientState, nil
}
