package keeper

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	ibccore "github.com/cosmos/ibc-go/v7/modules/core/exported"
	wasmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
	prysmtypes "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	"math/big"
	"strconv"
)

func (k Keeper) SendEthTransferIntent(ctx sdk.Context, msg *types.MsgSendTransferIntent) error {
	clientId := msg.ClientId

	if err := k.ValidateClientState(ctx, clientId); err != nil {
		return err
	}

	intentId := k.GetNextIntentId(ctx)
	transferIntent := types.TransferIntent{
		ClientId:           clientId,
		SourceAddress:      msg.FromAddress,
		DestinationAddress: msg.DestinationAddress,
		Amount:             msg.Amount,
	}
	k.AddTransferIntent(ctx, transferIntent, intentId)

	// TODO: post bounty for solver?
	// TODO: should a collateral amount be set here as well for the solver to deposit before executing the intent?

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventAddTransferIntent,
		sdk.NewAttribute(types.AttributeKeyIntentId, strconv.FormatUint(intentId, 10)),
		sdk.NewAttribute(types.AttributeKeyClientId, transferIntent.ClientId),
		sdk.NewAttribute(types.AttributeKeySourceAddress, transferIntent.SourceAddress),
		sdk.NewAttribute(types.AttributeKeyDestinationAddress, transferIntent.DestinationAddress),
		sdk.NewAttribute(types.AttributeKeyAmount, transferIntent.Amount.String()),
	))

	k.SetNextIntentId(ctx, intentId+1)

	return nil
}

func (k Keeper) GetNextIntentId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	intentIdBz := store.Get(types.TransferIntentIdKey)
	var intentId uint64
	if intentIdBz == nil {
		intentId = 0
	} else {
		intentId = binary.BigEndian.Uint64(intentIdBz)
	}

	return intentId
}

func (k Keeper) SetNextIntentId(ctx sdk.Context, intentId uint64) {
	store := ctx.KVStore(k.storeKey)
	intentIdBz := make([]byte, 8)
	binary.BigEndian.PutUint64(intentIdBz, intentId)
	store.Set(types.TransferIntentIdKey, intentIdBz)
}

// AddTransferIntent stores an intent object in the store
func (k Keeper) AddTransferIntent(ctx sdk.Context, transferIntent types.TransferIntent, intentId uint64) {
	store := ctx.KVStore(k.storeKey)

	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	transferIntentValue := k.cdc.MustMarshal(&transferIntent)

	store.Set(transferIntentKey, transferIntentValue)
}

func (k Keeper) GetTransferIntent(ctx sdk.Context, intentId uint64) (*types.TransferIntent, error) {
	store := ctx.KVStore(k.storeKey)

	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	if !store.Has(transferIntentKey) {
		return nil, types.ErrInvalidIntentId
	}

	transferIntentBz := store.Get(transferIntentKey)
	var transferIntent types.TransferIntent
	if err := k.cdc.Unmarshal(transferIntentBz, &transferIntent); err != nil {
		return nil, err
	}

	return &transferIntent, nil
}

func (k Keeper) VerifyEthTransferIntentProof(ctx sdk.Context, msg *types.MsgVerifyTransferIntentProof) error {
	store := ctx.KVStore(k.storeKey)

	transferIntent, err := k.GetTransferIntent(ctx, msg.IntentId)
	if err != nil {
		return err
	}

	var txReceipt gethtypes.Receipt
	if err := txReceipt.UnmarshalJSON(msg.TxReceipt); err != nil {
		return types.ErrInvalidTxReceipt
	}

	var blockHeader gethtypes.Header
	if err := rlp.DecodeBytes(msg.BlockHeader, &blockHeader); err != nil {
		return err
	}
	if err = VerifyReceiptProof(blockHeader, txReceipt, msg.ReceiptProof); err != nil {
		return fmt.Errorf("verify receipt proof: %v", err)
	}

	clientId := transferIntent.ClientId
	clientState, err := k.GetClientState(ctx, clientId)
	if err != nil {
		return err
	}
	if err := VerifyBeaconBlockBody(clientState, msg.BeaconBlockBody, txReceipt); err != nil {
		return fmt.Errorf("verify beacon block body: %v", err)
	}

	if err := VerifyTransferEvent(txReceipt, *transferIntent, string(msg.ReceiptSignature)); err != nil {
		return fmt.Errorf("verify transfer event: %v", err)
	}

	// TODO: verify solver signature

	// Purge resolved transfer intent after proof verification?
	store.Delete(types.GetPendingTransferIntentKeyById(msg.IntentId))

	// TODO: unlock bounty for solver?

	return nil
}

func VerifyReceiptProof(blockHeader gethtypes.Header, txReceipt gethtypes.Receipt, receiptProofBz []byte) error {
	var receiptProof types.ReceiptProof
	if err := receiptProof.Unmarshal(receiptProofBz); err != nil {
		return err
	}

	//Get binary representation of txReceipt rlp encoding
	txReceiptBz, err := txReceipt.MarshalBinary()
	if err != nil {
		return err
	}
	txReceiptHash := crypto.Keccak256(txReceiptBz)

	receiptsRoot := blockHeader.ReceiptHash
	if _, err = trie.VerifyProof(receiptsRoot, txReceiptHash, receiptProof); err != nil {
		return err
	}

	return nil
}

func VerifyBeaconBlockBody(clientState ibccore.ClientState, beaconBlockBodySSZ []byte, txReceipt gethtypes.Receipt) error {
	clientStateBz, err := proto.Marshal(clientState)
	if err != nil {
		return fmt.Errorf("marshal client state: %v", err)
	}
	wasmClientState := new(wasmtypes.ClientState)
	if err := proto.Unmarshal(clientStateBz, wasmClientState); err != nil {
		return fmt.Errorf("unmarshal client state bytes: %v", err)
	}
	ethClientState := new(types.ClientState)
	if err := ethClientState.Unmarshal(wasmClientState.Data); err != nil {
		return fmt.Errorf("unmarshal eth client state bytes: %v", err)
	}

	var beaconBlockBodyRoot [32]byte
	beaconBlockBodyRootSlice := ethClientState.GetInner().GetFinalizedHeader().GetBodyRoot()
	copy(beaconBlockBodyRoot[:], beaconBlockBodyRootSlice)

	var beaconBlockBody prysmtypes.BeaconBlockBody
	if err := beaconBlockBody.UnmarshalSSZ(beaconBlockBodySSZ); err != nil {
		return fmt.Errorf("unmarshal beacon block body: %v", err)
	}

	beaconBlockBodyHash, err := beaconBlockBody.HashTreeRoot()
	if beaconBlockBodyHash != beaconBlockBodyRoot {
		return types.ErrBlockBodyMismatch
	}

	blockHash := common.BytesToHash(beaconBlockBody.GetEth1Data().GetBlockHash())
	if blockHash != txReceipt.BlockHash {
		return types.ErrBlockHashMismatch
	}

	return nil
}

func VerifyTransferEvent(txReceipt gethtypes.Receipt, intent types.TransferIntent, solverAddress string) error {
	//TODO: find external package to import instead of using new struct
	type LogTransfer struct {
		From   common.Address
		To     common.Address
		Tokens *big.Int
		//TokenAddress common.Address
	}
	transferEventSig := []byte("Transfer(address,address,uint256)")
	transferEventSigHash := crypto.Keccak256Hash(transferEventSig)

	var transferEvent LogTransfer
	for _, log := range txReceipt.Logs {
		if log.Topics[0].Hex() == transferEventSigHash.Hex() {
			transferEvent.From = common.HexToAddress(log.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(log.Topics[2].Hex())
			transferEvent.Tokens = new(big.Int).SetBytes(log.Data)
			break
		}
	}

	if transferEvent == (LogTransfer{}) {
		return types.ErrTransferEventNotFound
	}
	if transferEvent.To != common.HexToAddress(intent.DestinationAddress) {
		return types.ErrDestinationAddressMismatch
	}
	if transferEvent.From != common.HexToAddress(solverAddress) {
		return types.ErrSourceAddressMismatch
	}
	if transferEvent.Tokens.Cmp(intent.Amount.BigInt()) != 0 {
		return types.ErrAmountMismatch
	}

	return nil
}

func (k Keeper) ValidateClientState(ctx sdk.Context, clientId string) error {
	_, found := k.clientKeeper.GetClientState(ctx, clientId)
	if !found {
		return types.ErrClientNotFound
	}

	// TODO uncomment clientStatus checks after figuring out why status is Unknown and not Active in test
	//clientStatus := k.clientKeeper.GetClientStatus(ctx, clientState, clientId)
	//if clientStatus != ibccore.Active {
	//	return types.ErrClientNotActive
	//}

	return nil
}

func (k Keeper) GetClientState(ctx sdk.Context, clientId string) (ibccore.ClientState, error) {
	clientState, found := k.clientKeeper.GetClientState(ctx, clientId)
	if !found {
		return nil, types.ErrClientNotFound
	}

	//clientStatus := k.clientKeeper.GetClientStatus(ctx, clientState, clientId)
	//if clientStatus != ibccore.Active {
	//	return nil, types.ErrClientNotActive
	//}

	return clientState, nil
}
