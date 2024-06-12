package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

type msgServer struct {
	Keeper
}

func (k Keeper) SendTransferIntent(goCtx context.Context, msg *types.MsgSendTransferIntent) (*types.MsgSendTransferIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err := k.SendEthTransferIntent(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendTransferIntentResponse{}, nil
}

func (k Keeper) VerifyTransferIntentProof(goCtx context.Context, msg *types.MsgVerifyTransferIntentProof) (*types.MsgVerifyTransferIntentProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err := k.VerifyEthTransferIntentProof(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgVerifyTransferIntentProofResponse{}, nil
}

func (k Keeper) TriggerTransferIntentTimeout(goCtx context.Context, msg *types.MsgTriggerTransferIntentTimeout) (*types.MsgTriggerTransferIntentTimeoutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.TriggerEthTransferIntentTimeout(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgTriggerTransferIntentTimeoutResponse{}, nil
}
