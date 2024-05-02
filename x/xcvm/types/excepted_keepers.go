package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

type TransferKeeper interface {
	GetReceiveEnabled(ctx sdk.Context) bool
	Transfer(goCtx context.Context, msg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error)
	HasDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) bool
	SetDenomTrace(ctx sdk.Context, denomTrace transfertypes.DenomTrace)
}
