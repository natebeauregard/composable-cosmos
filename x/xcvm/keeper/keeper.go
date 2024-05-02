package keeper

import (
	"fmt"

	"github.com/notional-labs/composable/v6/x/xcvm/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"

	"github.com/cometbft/cometbft/libs/log"
)

// Keeper struct
type Keeper struct {
	cdc         codec.Codec
	storeKey    storetypes.StoreKey
	ICS4Wrapper porttypes.ICS4Wrapper
	// transferKeeper types.TransferKeeper

	// the address capable of executing a privileged message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper returns keeper
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, authority string) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		authority: authority,
	}
}

// Logger returns logger
func (Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
