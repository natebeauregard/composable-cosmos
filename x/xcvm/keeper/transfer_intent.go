package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibccore "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/golang/protobuf/proto"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
	// prysmtypes "github.com/prysmaticlabs/prysm/proto/eth/v1"
)

type receiptProof struct {
	Receipts map[[32]byte][]byte
}

func (rp receiptProof) Has(key []byte) (bool, error) {
	if len(key) != 32 {
		return false, types.ErrInvalidReceiptKey
	}
	var keyArr [32]byte
	copy(keyArr[:], key[:32])
	_, ok := rp.Receipts[keyArr]
	return ok, nil
}

func (rp receiptProof) Get(key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, types.ErrInvalidReceiptKey
	}
	var keyArr [32]byte
	copy(keyArr[:], key[:32])
	value, ok := rp.Receipts[keyArr]
	if !ok {
		return nil, types.ErrReceiptNotFound
	}
	return value, nil
}

func (k Keeper) SendEthTransferIntent(ctx sdk.Context, msg *types.MsgSendTransferIntent) error {
	clientId := msg.ClientId

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

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
	if err := txReceipt.UnmarshalBinary(msg.TxReceipt); err != nil {
		return types.ErrInvalidTxReceipt
	}
	var blockHeader gethtypes.Header
	if err := rlp.DecodeBytes(msg.BlockHeader, &blockHeader); err != nil {
		return err
	}

	var receiptProof receiptProof
	if err := json.Unmarshal(msg.ReceiptProof, &receiptProof); err != nil {
		return err
	}

	// Get binary representation of txReceipt rlp encoding
	txReceiptBz, err := txReceipt.MarshalBinary()
	if err != nil {
		return err
	}
	txReceiptHash := crypto.Keccak256(txReceiptBz)

	receiptsRoot := blockHeader.ReceiptHash
	if _, err = trie.VerifyProof(receiptsRoot, txReceiptHash, receiptProof); err != nil {
		return err
	}

	clientId := transferIntent.ClientId
	clientState, err := k.GetClientState(ctx, clientId)
	if err != nil {
		return err
	}
	clientStateBz, err := proto.Marshal(clientState)
	if err != nil {
		return fmt.Errorf("marshal client state: %v", err)
	}
	ethClientState := new(types.ClientState)
	if err := proto.Unmarshal(clientStateBz, ethClientState); err != nil {
		return fmt.Errorf("unmarshal client state bytes: %v", err)
	}

	var beaconBlockBodyRoot [32]byte
	beaconBlockBodyRootSlice := ethClientState.GetInner().GetFinalizedHeader().GetBodyRoot()
	copy(beaconBlockBodyRoot[:], beaconBlockBodyRootSlice)

	// TODO: investigate prysm dependency error
	// var beaconBlockBody prysmtypes.BeaconBlockBody
	// if err := beaconBlockBody.UnmarshalSSZ(msg.BeaconBlockBody); err != nil {
	// 	return fmt.Errorf("unmarshal beacon block body: %v", err)
	// }

	// beaconBlockBodyHash, err := beaconBlockBody.HashTreeRoot()
	// if beaconBlockBodyHash != beaconBlockBodyRoot {
	// 	return types.ErrBlockBodyMismatch
	// }

	// blockHash := common.BytesToHash(beaconBlockBody.GetEth1Data().GetBlockHash())
	// if blockHash != txReceipt.BlockHash {
	// 	return types.ErrBlockHashMismatch
	// }

	if err := VerifyTransferEvent(txReceipt, *transferIntent, string(msg.ReceiptSignature)); err != nil {
		return err
	}

	// Purge resolved transfer intent after proof verification?
	store.Delete(types.GetPendingTransferIntentKeyById(msg.IntentId))

	// TODO: unlock bounty for solver?

	return nil
}

func VerifyTransferEvent(txReceipt gethtypes.Receipt, intent types.TransferIntent, solverAddress string) error {
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
	if transferEvent.From != common.HexToAddress(solverAddress) {
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
