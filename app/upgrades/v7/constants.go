package v6

import (
	"github.com/notional-labs/composable/v6/app/upgrades"

	store "cosmossdk.io/store"
)

const (
	UpgradeName = "v7"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
