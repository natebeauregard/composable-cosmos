package v6

import (
	"github.com/notional-labs/composable/v6/app/upgrades"

	store "cosmossdk.io/store"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

const (
	// UpgradeName defines the on-chain upgrade name for the composable upgrade.
	UpgradeName = "v6"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{authz.ModuleName},
	},
}
