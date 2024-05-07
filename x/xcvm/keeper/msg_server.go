package keeper

import (
	"context"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

// ExecuteCvmProgram implements types.MsgServer.
func (k Keeper) ExecuteCvmProgram(ctx context.Context, msg *types.MsgExecuteCvmProgram) (*types.MsgExecuteCvmProgramResponse, error) {
	program := msg.Program
	// TODO: store composable source port/source channel for Ethereum IBC connection in module?
	// TODO: figure out proper sdk context
	return nil, k.executeCvmProgram(program, sdk.NewContext(nil, tmproto.Header{}, true, nil), msg.FromAddress, "placeholderIBCPort", "placeholderIBCChannel")
}
