package v4

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	"github.com/notional-labs/composable/v6/app/upgrades"
	ratelimitmoduletypes "github.com/notional-labs/composable/v6/x/ratelimit/types"

	store "cosmossdk.io/store"

	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v4"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{wasmtypes.StoreKey, ibchookstypes.StoreKey, ratelimitmoduletypes.StoreKey, icahosttypes.StoreKey},
		Deleted: []string{},
	},
}
