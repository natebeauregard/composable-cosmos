package keeper

import (
	"time"

	"github.com/notional-labs/composable/v6/x/transfermiddleware/types"

	errorsmod "cosmossdk.io/errors"

	storetypes "cosmossdk.io/store"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/cometbft/cometbft/libs/log"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type Keeper struct {
	cdc            codec.BinaryCodec
	storeKey       storetypes.StoreKey
	paramSpace     paramtypes.Subspace
	ICS4Wrapper    porttypes.ICS4Wrapper
	bankKeeper     types.BankKeeper
	transferKeeper types.TransferKeeper

	// the address capable of executing a AddParachainIBCTokenInfo and RemoveParachainIBCTokenInfo message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper returns a new instance of the x/ibchooks keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	codec codec.BinaryCodec,
	ics4Wrapper porttypes.ICS4Wrapper,
	transferKeeper types.TransferKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:       storeKey,
		paramSpace:     paramSpace,
		transferKeeper: transferKeeper,
		bankKeeper:     bankKeeper,
		cdc:            codec,
		ICS4Wrapper:    ics4Wrapper,
		authority:      authority,
	}
}

// TODO: testing
// AddParachainIBCTokenInfo add new parachain token information token to chain state.
func (k Keeper) AddParachainIBCInfo(ctx sdk.Context, ibcDenom, channelID, nativeDenom, assetID string) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.GetKeyParachainIBCTokenInfoByAssetID(assetID)) {
		return errorsmod.Wrapf(types.ErrMultipleMapping, "duplicate assetID")
	}
	if store.Has(types.GetKeyNativeDenomAndIbcSecondaryIndex(ibcDenom)) {
		return errorsmod.Wrapf(types.ErrMultipleMapping, "duplicate IBC denom")
	}
	if store.Has(types.GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom)) {
		return errorsmod.Wrapf(types.ErrMultipleMapping, "duplicate native denom")
	}

	info := types.ParachainIBCTokenInfo{
		IbcDenom:    ibcDenom,
		ChannelID:   channelID,
		NativeDenom: nativeDenom,
		AssetId:     assetID,
	}

	bz, err := k.cdc.Marshal(&info)
	if err != nil {
		return err
	}

	store.Set(types.GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom), bz)
	store.Set(types.GetKeyParachainIBCTokenInfoByAssetID(assetID), bz)
	store.Set(types.GetKeyNativeDenomAndIbcSecondaryIndex(ibcDenom), []byte(nativeDenom))
	return nil
}

// TODO: testing
// AddParachainIBCInfoToRemoveList add parachain token information token to remove list.
func (k Keeper) AddParachainIBCInfoToRemoveList(ctx sdk.Context, nativeDenom string) (time.Time, error) {
	params := k.GetParams(ctx)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom)) {
		return time.Time{}, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "Token %v info not found", nativeDenom)
	}

	// Add to remove list
	removeTime := ctx.BlockTime().Add(params.Duration)
	removeToken := types.RemoveParachainIBCTokenInfo{
		NativeDenom: nativeDenom,
		RemoveTime:  removeTime,
	}

	bz, err := k.cdc.Marshal(&removeToken)
	if err != nil {
		return time.Time{}, err
	}

	store.Set(types.GetKeyParachainIBCTokenRemoveListByNativeDenom(nativeDenom), bz)
	return removeTime, nil
}

// TODO: testing
// IterateRemoveListInfo iterate all parachain token in remove list.
func (k Keeper) IterateRemoveListInfo(ctx sdk.Context, cb func(removeInfo types.RemoveParachainIBCTokenInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyParachainIBCTokenRemoveListByNativeDenom)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var removeInfo types.RemoveParachainIBCTokenInfo
		k.cdc.MustUnmarshal(iterator.Value(), &removeInfo)
		if cb(removeInfo) {
			break
		}
	}
}

// TODO: testing
// RemoveParachainIBCTokenInfo remove parachain token information from chain state.
func (k Keeper) RemoveParachainIBCInfo(ctx sdk.Context, nativeDenom string) error {
	if !k.hasParachainIBCTokenInfo(ctx, nativeDenom) {
		return types.NotRegisteredNativeDenom
	}

	// get the IBCdenom
	tokenInfo := k.GetParachainIBCTokenInfoByNativeDenom(ctx, nativeDenom)
	ibcDenom := tokenInfo.IbcDenom
	assetID := tokenInfo.AssetId

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom))
	store.Delete(types.GetKeyParachainIBCTokenInfoByAssetID(assetID))
	store.Delete(types.GetKeyNativeDenomAndIbcSecondaryIndex(ibcDenom))

	return nil
}

func (k Keeper) SetAllowRlyAddress(ctx sdk.Context, rlyAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetKeyByRlyAddress(rlyAddress), []byte{1})
}

func (k Keeper) DeleteAllowRlyAddress(ctx sdk.Context, rlyAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetKeyByRlyAddress(rlyAddress))
}

func (k Keeper) HasAllowRlyAddress(ctx sdk.Context, rlyAddress string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKeyByRlyAddress(rlyAddress)

	if store.Has(key) {
		return true
	}

	prefixStore := prefix.NewStore(store, types.KeyRlyAddress)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	// there are not records => so it is permissionless
	return !iter.Valid()
}

func (k Keeper) IterateAllowRlyAddress(ctx sdk.Context, cb func(rlyAddress string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.KeyRlyAddress)
	iterator := sdk.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		rlyAddress := string(iterator.Key())
		if cb(rlyAddress) {
			break
		}
	}
}

func (k Keeper) HasParachainIBCTokenInfoByNativeDenom(ctx sdk.Context, nativeDenom string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom)

	return store.Has(key)
}

func (k Keeper) HasParachainIBCTokenInfoByAssetID(ctx sdk.Context, assetID string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKeyParachainIBCTokenInfoByAssetID(assetID)

	return store.Has(key)
}

// TODO: testing
// GetParachainIBCTokenInfo add new information about parachain token to chain state.
func (k Keeper) GetParachainIBCTokenInfoByNativeDenom(ctx sdk.Context, nativeDenom string) (info types.ParachainIBCTokenInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetKeyParachainIBCTokenInfoByNativeDenom(nativeDenom))

	k.cdc.Unmarshal(bz, &info) //nolint:errcheck // TODO: handle error

	return info
}

func (k Keeper) GetParachainIBCTokenInfoByAssetID(ctx sdk.Context, assetID string) (info types.ParachainIBCTokenInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetKeyParachainIBCTokenInfoByAssetID(assetID))

	k.cdc.Unmarshal(bz, &info) //nolint:errcheck // TODO: handle error

	return info
}

func (k Keeper) GetNativeDenomByIBCDenomSecondaryIndex(ctx sdk.Context, ibcDenom string) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetKeyNativeDenomAndIbcSecondaryIndex(ibcDenom))

	return string(bz)
}

func (k Keeper) GetTotalEscrowedToken(ctx sdk.Context) (coins sdk.Coins) {
	k.IterateParaTokenInfos(ctx, func(index int64, info types.ParachainIBCTokenInfo) (stop bool) {
		escrowIbcCoin := k.bankKeeper.GetBalance(ctx, transfertypes.GetEscrowAddress(transfertypes.PortID, info.ChannelID), info.IbcDenom)
		escrowNativeCoin := k.bankKeeper.GetBalance(ctx, transfertypes.GetEscrowAddress(transfertypes.PortID, info.ChannelID), info.NativeDenom)
		coins = append(coins, escrowIbcCoin, escrowNativeCoin)
		return false
	})

	return coins
}

func (Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+exported.ModuleName+"-"+types.ModuleName)
}
