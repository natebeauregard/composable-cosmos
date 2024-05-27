package slashing

import (
	"github.com/notional-labs/composable/v6/bech32-migration/utils"

	storetypes "cosmossdk.io/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

func MigrateAddressBech32(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	ctx.Logger().Info("Migration of address bech32 for slashing module begin")
	validatorSigningInfoCount := uint64(0)
	utils.IterateStoreByPrefix(ctx, storeKey, types.ValidatorSigningInfoKeyPrefix, func(bz []byte) []byte {
		validatorSigningInfo := types.ValidatorSigningInfo{}
		cdc.MustUnmarshal(bz, &validatorSigningInfo)
		validatorSigningInfo.Address = utils.ConvertConsAddr(validatorSigningInfo.Address)
		validatorSigningInfoCount++
		return cdc.MustMarshal(&validatorSigningInfo)
	})

	ctx.Logger().Info(
		"Migration of address bech32 for slashing module done",
		"validator_signing_info_count", validatorSigningInfoCount,
	)
}
