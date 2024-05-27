package centauri

import (
	bech32authmigration "github.com/notional-labs/composable/v6/bech32-migration/auth"
	bech32govmigration "github.com/notional-labs/composable/v6/bech32-migration/gov"
	bech32slashingmigration "github.com/notional-labs/composable/v6/bech32-migration/slashing"
	bech32stakingmigration "github.com/notional-labs/composable/v6/bech32-migration/staking"

	"cosmossdk.io/math"

	storetypes "cosmossdk.io/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keys map[string]*storetypes.KVStoreKey, codec codec.Codec,
	slashingKeeper *slashingkeeper.Keeper,
	govkeeper *govkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Migration prefix
		ctx.Logger().Info("First step: Migrate addresses stored in bech32 form to use new prefix")
		bech32stakingmigration.MigrateAddressBech32(ctx, keys[stakingtypes.StoreKey], codec)
		bech32slashingmigration.MigrateAddressBech32(ctx, keys[slashingtypes.StoreKey], codec)
		bech32govmigration.MigrateAddressBech32(ctx, keys[govtypes.StoreKey], codec)
		bech32authmigration.MigrateAddressBech32(ctx, keys[authtypes.StoreKey], codec)

		// Slashing params change
		newParamsSet := slashingKeeper.GetParams(ctx)
		newParamsSet.SlashFractionDowntime = math.LegacyNewDecWithPrec(1, 5)
		slashingKeeper.SetParams(ctx, newParamsSet) //nolint:errcheck // TODO: handle error

		// Gov params change: minium deposit
		newGovParamsSet := govkeeper.GetParams(ctx)
		newGovParamsSet.MinInitialDepositRatio = sdk.NewDecWithPrec(1, 2).String() // this is 1%
		govkeeper.SetParams(ctx, newGovParamsSet)                                  //nolint:errcheck // TODO: handle error

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
