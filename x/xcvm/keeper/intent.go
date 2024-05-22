package keeper

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
	ibccore "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	ssz "github.com/ferranbt/fastssz"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

type receiptProof struct {
	receipts map[[32]byte][]byte
}

func (rp receiptProof) Has(key []byte) (bool, error) {
	if len(key) != 32 {
		return false, types.ErrInvalidReceiptKey
	}
	var keyArr [32]byte
	copy(keyArr[:], key[:32])
	_, ok := rp.receipts[keyArr]
	return ok, nil
}

func (rp receiptProof) Get(key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, types.ErrInvalidReceiptKey
	}
	var keyArr [32]byte
	copy(keyArr[:], key[:32])
	value, ok := rp.receipts[keyArr]
	if !ok {
		return nil, types.ErrReceiptNotFound
	}
	return value, nil
}

func (k Keeper) SendValidatedTransferIntent(ctx sdk.Context, msg *types.MsgSendTransferIntent) error {
	clientId := msg.ClientId

	err := k.ValidateClientState(ctx, clientId)
	if err != nil {
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
	intentId := binary.BigEndian.Uint64(intentIdBz)

	return intentId
}

func (k Keeper) SetNextIntentId(ctx sdk.Context, intentId uint64) {
	store := ctx.KVStore(k.storeKey)
	intentIdBz := make([]byte, 8)
	binary.BigEndian.PutUint64(intentIdBz, intentId)
	store.Set(types.TransferIntentIdKey, intentIdBz)
}

// Stores an intent object in the store
func (k Keeper) AddTransferIntent(ctx sdk.Context, transferIntent types.TransferIntent, intentId uint64) {
	store := ctx.KVStore(k.storeKey)

	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	transferIntentValue := k.cdc.MustMarshal(&transferIntent)

	store.Set(transferIntentKey, transferIntentValue)
}

func (k Keeper) VerifyIntentProof(ctx sdk.Context, msg *types.MsgVerifyTransferIntentProof) error {
	store := ctx.KVStore(k.storeKey)

	intentId := msg.IntentId
	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	if !store.Has(transferIntentKey) {
		return types.ErrInvalidIntentId
	}

	transferIntentBz := store.Get(transferIntentKey)
	var transferIntent types.TransferIntent
	err := k.cdc.Unmarshal(transferIntentBz, &transferIntent)
	if err != nil {
		return err
	}

	var txReceipt gethtypes.Receipt
	err = txReceipt.UnmarshalBinary(msg.TxReceipt)
	if err != nil {
		return types.ErrInvalidTxReceipt
	}
	signature := msg.ReceiptSignature
	var blockHeader gethtypes.Header
	err = rlp.DecodeBytes(msg.BlockHeader, &blockHeader)
	if err != nil {
		return err
	}

	var receiptProof receiptProof
	err = json.Unmarshal(msg.ReceiptProof, &receiptProof)
	if err != nil {
		return err
	}

	// Get binary representation of txReceipt rlp encoding
	txReceiptBz, err := txReceipt.MarshalBinary()
	if err != nil {
		return err
	}
	txReceiptHash := crypto.Keccak256(txReceiptBz)

	receiptsRoot := blockHeader.ReceiptHash
	_, err = trie.VerifyProof(receiptsRoot, txReceiptHash, receiptProof)
	if err != nil {
		return err
	}

	clientId := transferIntent.ClientId
	clientState, err := k.GetClientState(ctx, clientId)
	if err != nil {
		return err
	}

	// TODO: how to go from clientState to ethClientState?
	ethClientState := proto.Unmarshal()
	var beaconBlockBodyRoot [32]byte
	beaconBlockBodyRootSlice := ethClientState.GetInner().GetFinalizedHeader().GetBodyRoot()
	copy(beaconBlockBodyRoot[:], beaconBlockBodyRootSlice)

	var beaconBlockBody types.BeaconBlockBody
	err = beaconBlockBody.Unmarshal(msg.BeaconBlockBody)
	if err != nil {
		return err
	}

	beaconBlockBodyHash, err := ssz.HashWithDefaultHasher(beaconBlockBody)
	if beaconBlockBodyHash != beaconBlockBodyRoot {
		return types.ErrBlockBodyMismatch
	}

	blockHash := common.BytesToHash(beaconBlockBody.GetEth1Data().GetBlockHash())
	// TODO: receipt root corresponds to block hash and make sure that receipt is included in the block

	// Purge resolved transfer intent after proof verification
	store.Delete(transferIntentKey)

	// TODO: unlock bounty for solver?

	// addd tests for proof verification

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventVerifyTransferIntentProof,
		sdk.NewAttribute(types.AttributeKeyIntentId, strconv.FormatUint(intentId, 10)),
	))

	return nil
}

func VerifyTransferEvent(txReceipt gethtypes.Receipt, intent types.TransferIntent) error {
	//TODO: find external package to import instead of using new struct
	type LogTransfer struct {
		From         common.Address
		To           common.Address
		Tokens       *big.Int
		TokenAddress common.Address
	}
	transferEventSig := []byte("Transfer(address,address,uint256)")
	transferEventSigHash := crypto.Keccak256Hash(transferEventSig)

	erc20Abi, err := os.ReadFile("erc20.abi.json")
	if err != nil {
		return err
	}
	// TODO: store contract abi as constant instead of needing to ready from JSON each call?
	contractAbi, err := abi.JSON(strings.NewReader(string(erc20Abi)))
	if err != nil {
		return err
	}

	var transferEvent LogTransfer
	for _, log := range txReceipt.Logs {
		if log.Topics[0].Hex() == transferEventSigHash.Hex() {
			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", log.Data)
			if err != nil {
				return err
			}
			transferEvent.TokenAddress = log.Address
			transferEvent.From = common.HexToAddress(log.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(log.Topics[2].Hex())
			break
		}
	}

	if transferEvent == (LogTransfer{}) {
		return types.ErrTransferEventNotFound
	}
	if transferEvent.To != common.HexToAddress(intent.DestinationAddress) {
		return types.ErrDestinationAddressMismatch
	}
	if transferEvent.From != common.HexToAddress(intent.SourceAddress) {
		return types.ErrSourceAddressMismatch
	}
	if transferEvent.Tokens.Cmp(intent.Amount.BigInt()) != 0 {
		return types.ErrAmountMismatch
	}

	return nil
}

func (k Keeper) ValidateClientState(ctx sdk.Context, clientId string) error {
	clientState, found := k.clientKeeper.GetClientState(ctx, clientId)
	if !found {
		return types.ErrClientNotFound
	}

	clientStatus := k.clientKeeper.GetClientStatus(ctx, clientState, clientId)
	if clientStatus != ibccore.Active {
		return types.ErrClientNotActive
	}

	return nil
}

func (k Keeper) GetClientState(ctx sdk.Context, clientId string) (ibccore.ClientState, error) {
	clientState, found := k.clientKeeper.GetClientState(ctx, clientId)
	if !found {
		return nil, types.ErrClientNotFound
	}

	clientStatus := k.clientKeeper.GetClientStatus(ctx, clientState, clientId)
	if clientStatus != ibccore.Active {
		return nil, types.ErrClientNotActive
	}

	return clientState, nil
}
