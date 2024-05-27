package keeper

import (
	"fmt"

	"github.com/notional-labs/composable/v6/x/xcvm/types"

	storetypes "cosmossdk.io/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/libs/log"
)

// Keeper struct
type Keeper struct {
	cdc          codec.Codec
	storeKey     storetypes.StoreKey
	clientKeeper types.ClientKeeper

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
