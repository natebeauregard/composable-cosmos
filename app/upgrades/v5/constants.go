package v5

import (
	"github.com/notional-labs/composable/v6/app/upgrades"
	txboundary "github.com/notional-labs/composable/v6/x/tx-boundary/types"

	store "cosmossdk.io/store"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v5"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{txboundary.ModuleName},
	},
}
