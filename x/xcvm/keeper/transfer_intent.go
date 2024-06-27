package keeper

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
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

// SendEthTransferIntent creates and stores a new transfer intent and emits an event to alert potential solvers
func (k Keeper) SendEthTransferIntent(ctx sdk.Context, msg *types.MsgSendTransferIntent) error {
	clientId := msg.ClientId

	if err := k.validateClientState(ctx, clientId); err != nil {
		return fmt.Errorf("validate client state: %v", err)
	}

	intentId := k.getNextIntentId(ctx)
	transferIntent := types.TransferIntent{
		ClientId:           clientId,
		SourceAddress:      msg.FromAddress,
		DestinationAddress: msg.DestinationAddress,
		TimeoutHeight:      msg.TimeoutHeight,
		TransferTokens:     msg.TransferTokens,
		Bounty:             msg.Bounty,
	}
	k.AddTransferIntent(ctx, transferIntent, intentId)

	userAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return fmt.Errorf("acc address conversion: %v", err)
	}
	coins := sdk.NewCoins(msg.Bounty)

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, userAddress, types.ModuleName, coins); err != nil {
		return fmt.Errorf("send coins from account to module: %v", err)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventAddTransferIntent,
		sdk.NewAttribute(types.AttributeKeyIntentId, strconv.FormatUint(intentId, 10)),
		sdk.NewAttribute(types.AttributeKeyClientId, transferIntent.ClientId),
		sdk.NewAttribute(types.AttributeKeySourceAddress, transferIntent.SourceAddress),
		sdk.NewAttribute(types.AttributeKeyDestinationAddress, transferIntent.DestinationAddress),
		sdk.NewAttribute(types.AttributeKeyTimeout, strconv.FormatInt(transferIntent.TimeoutHeight, 10)),
		sdk.NewAttribute(types.AttributeKeyAmount, transferIntent.TransferTokens.String()),
		sdk.NewAttribute(types.AttributeKeyBounty, transferIntent.Bounty.String()),
	))

	k.SetNextIntentId(ctx, intentId+1)

	return nil
}

// getNextIntentId retrieves the next intent id from the store
func (k Keeper) getNextIntentId(ctx sdk.Context) uint64 {
	kvStore := ctx.KVStore(k.storeKey)

	intentIdBz := kvStore.Get(types.TransferIntentIdKey)
	var intentId uint64
	if intentIdBz == nil {
		intentId = 0
	} else {
		intentId = binary.BigEndian.Uint64(intentIdBz)
	}

	return intentId
}

// SetNextIntentId sets the next intent id in the store
func (k Keeper) SetNextIntentId(ctx sdk.Context, intentId uint64) {
	kvStore := ctx.KVStore(k.storeKey)
	intentIdBz := make([]byte, 8)
	binary.BigEndian.PutUint64(intentIdBz, intentId)
	kvStore.Set(types.TransferIntentIdKey, intentIdBz)
}

// AddTransferIntent stores an intent object in the store
func (k Keeper) AddTransferIntent(ctx sdk.Context, transferIntent types.TransferIntent, intentId uint64) {
	kvStore := ctx.KVStore(k.storeKey)

	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	transferIntentValue := k.cdc.MustMarshal(&transferIntent)

	kvStore.Set(transferIntentKey, transferIntentValue)
}

// GetTransferIntent retrieves a transfer intent from the store
func (k Keeper) GetTransferIntent(ctx sdk.Context, intentId uint64) (*types.TransferIntent, error) {
	kvStore := ctx.KVStore(k.storeKey)

	transferIntentKey := types.GetPendingTransferIntentKeyById(intentId)
	if !kvStore.Has(transferIntentKey) {
		return nil, types.ErrInvalidIntentId
	}

	transferIntentBz := kvStore.Get(transferIntentKey)
	var transferIntent types.TransferIntent
	if err := k.cdc.Unmarshal(transferIntentBz, &transferIntent); err != nil {
		return nil, err
	}

	return &transferIntent, nil
}

// VerifyEthTransferIntentProof verifies the proof submitted by the solver to claim the bounty for executing the transfer intent
func (k Keeper) VerifyEthTransferIntentProof(ctx sdk.Context, msg *types.MsgVerifyTransferIntentProof) error {
	kvStore := ctx.KVStore(k.storeKey)

	// verify that the transfer intent exists
	transferIntent, err := k.GetTransferIntent(ctx, msg.IntentId)
	if err != nil {
		return fmt.Errorf("get transfer intent: %v", err)
	}

	// verify that the transfer intent request is not timed out
	currentBlockHeight := ctx.BlockHeight()
	if currentBlockHeight >= transferIntent.TimeoutHeight {
		return types.ErrProofSubmittedAfterTimeout
	}

	// verify that the transaction receipt is valid and included in the execution block header
	var txReceipt gethtypes.Receipt
	if err := txReceipt.UnmarshalJSON(msg.TxReceipt); err != nil {
		return types.ErrInvalidTxReceipt
	}
	var blockHeader gethtypes.Header
	if err := rlp.DecodeBytes(msg.BlockHeader, &blockHeader); err != nil {
		return fmt.Errorf("decode block header: %v", err)
	}
	txReceiptHash, err := getTxReceiptHash(txReceipt)
	if err != nil {
		return fmt.Errorf("get tx receipt hash: %v", err)
	}
	if err = verifyReceiptProof(blockHeader, txReceiptHash, msg.ReceiptProof); err != nil {
		return fmt.Errorf("verify receipt proof: %v", err)
	}

	clientId := transferIntent.ClientId
	clientState, err := k.getClientState(ctx, clientId)
	if err != nil {
		return fmt.Errorf("get client state: %v", err)
	}

	clientStateBeaconBlockHeader, err := getClientStateBeaconBlockHeader(clientState)
	if err != nil {
		return fmt.Errorf("get client state beacon block header: %v", err)
	}

	// If the light client is ahead of the block where the intent was executed, verify that the provided beacon block headers are valid
	var intentBeaconBlockHeader *types.BeaconBlockHeader
	if len(msg.BeaconBlockHeaders) > 0 {
		// Include the light client state as the last block header in the chain to verify
		beaconBlockHeaders := append(msg.BeaconBlockHeaders, clientStateBeaconBlockHeader)
		if err := verifyPreviousBeaconBlockHeaders(beaconBlockHeaders); err != nil {
			return fmt.Errorf("verify beacon block headers: %v", err)
		}
		intentBeaconBlockHeader = msg.BeaconBlockHeaders[0]
	} else {
		intentBeaconBlockHeader = clientStateBeaconBlockHeader
	}

	// Verify that the provided beacon block body matches the beacon block header's body root and the transaction receipt's block hash
	if err := verifyBeaconBlockBody(*intentBeaconBlockHeader, msg.BeaconBlockBody, txReceipt); err != nil {
		return fmt.Errorf("verify beacon block body: %v", err)
	}

	// Verify that the solver signed over the receipt
	solverPublicKey, err := crypto.DecompressPubkey(msg.PublicKey)
	if err != nil {
		return fmt.Errorf("decompress public key: %v", err)
	}
	if err := verifyReceiptSignature(solverPublicKey, txReceiptHash, txReceipt.BlockHash, msg.ReceiptSignature); err != nil {
		return fmt.Errorf("verify receipt signature: %v", err)
	}

	// Verify that the ERC-20 transfer log emitted in the transaction receipt includes the correct destination address and amount specified in the transfer intent
	if err := verifyTransferEvent(txReceipt, *transferIntent, solverPublicKey); err != nil {
		return fmt.Errorf("verify transfer event: %v", err)
	}

	// Verify that the receipt is unique
	if err := verifyReceiptUniqueness(kvStore, txReceiptHash, txReceipt.BlockHash); err != nil {
		return fmt.Errorf("verify receipt uniqueness: %v", err)
	}

	// Get solver's account address
	accAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return fmt.Errorf("acc address conversion: %v", err)
	}
	coins := sdk.NewCoins(transferIntent.Bounty)

	// Transfer bounty to solver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddress, coins); err != nil {
		return fmt.Errorf("unlock bounty for solver: %v", err)
	}

	// Purge resolved transfer intent after proof verification
	kvStore.Delete(types.GetPendingTransferIntentKeyById(msg.IntentId))

	return nil
}

// TriggerEthTransferIntentTimeout releases the bounty to the original sender if the intent has timed out
func (k Keeper) TriggerEthTransferIntentTimeout(ctx sdk.Context, msg *types.MsgTriggerTransferIntentTimeout) error {
	kvStore := ctx.KVStore(k.storeKey)

	// Check that the transfer intent exists and is valid
	transferIntent, err := k.GetTransferIntent(ctx, msg.IntentId)
	if err != nil {
		return fmt.Errorf("get transfer intent: %v", err)
	}

	// Check if the sender is the original sender of the intent
	if msg.Sender != transferIntent.SourceAddress {
		return types.ErrInvalidSenderAddress
	}

	// Check if the intent is able to be timed out
	if ctx.BlockHeight() < transferIntent.TimeoutHeight {
		return types.ErrPrematureTimeoutTrigger
	}

	// Get user's account address
	accAddress, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return fmt.Errorf("acc address conversion: %v", err)
	}
	coins := sdk.NewCoins(transferIntent.Bounty)

	// Release bounty to original sender
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddress, coins); err != nil {
		return fmt.Errorf("release unclaimed bounty to original sender: %v", err)
	}

	// Remove transfer intent from store
	kvStore.Delete(types.GetPendingTransferIntentKeyById(msg.IntentId))

	return nil
}

// verifyReceiptSignature verifies that the solver signed over its proof using the same Ethereum account with which it fulfills the intent
func verifyReceiptSignature(solverPublicKey *ecdsa.PublicKey, txReceiptHash []byte, blockHash common.Hash, receiptSig []byte) error {
	encPublicKey := crypto.FromECDSAPub(solverPublicKey)
	receiptDataHash := crypto.Keccak256(append(txReceiptHash, blockHash.Bytes()...))

	// Remove the recovery id from the receipt signature before verifying
	if !crypto.VerifySignature(encPublicKey, receiptDataHash, receiptSig[:64]) {
		return types.ErrInvalidReceiptSignature
	}
	return nil
}

// verifyReceiptProof verifies that the receipt is included in the block header's receipt trie
func verifyReceiptProof(blockHeader gethtypes.Header, txReceiptHash []byte, receiptProofBz []byte) error {
	var receiptProof types.ReceiptProof
	if err := receiptProof.Unmarshal(receiptProofBz); err != nil {
		return fmt.Errorf("unmarshal receipt proof: %v", err)
	}

	receiptsRoot := blockHeader.ReceiptHash
	if _, err := trie.VerifyProof(receiptsRoot, txReceiptHash, receiptProof); err != nil {
		return fmt.Errorf("verify proof: %v", err)
	}

	return nil
}

// getTxReceiptHash returns the hash of the transaction receipt
func getTxReceiptHash(txReceipt gethtypes.Receipt) ([]byte, error) {
	//Get binary representation of txReceipt rlp encoding
	txReceiptBz, err := txReceipt.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("get tx receipt hash: %v", err)
	}
	txReceiptHash := crypto.Keccak256(txReceiptBz)
	return txReceiptHash, nil
}

// verifyReceiptUniqueness verifies that the receipt has not been processed before
func verifyReceiptUniqueness(store store.KVStore, txReceiptHash []byte, blockHash common.Hash) error {
	receiptKey := types.GetUsedReceiptKey(txReceiptHash, blockHash)
	if store.Has(receiptKey) {
		return types.ErrReceiptAlreadyProcessed
	}

	store.Set(receiptKey, []byte{0x01})
	return nil
}

// verifyPreviousBeaconBlockHeaders verifies that the provided previous beacon block headers are valid and in the correct order
func verifyPreviousBeaconBlockHeaders(beaconBlockHeaders []*types.BeaconBlockHeader) error {
	intentBeaconBlockHeader := *beaconBlockHeaders[0]
	headerHash, err := intentBeaconBlockHeader.Hash()
	if err != nil {
		return fmt.Errorf("hash intent beacon block header: %v", err)
	}

	var currentParentRoot [32]byte
	for _, header := range beaconBlockHeaders[1:] {
		copy(currentParentRoot[:], header.GetParentRoot())
		if currentParentRoot != headerHash {
			return types.ErrInvalidBlockHeaders
		}
		headerHash, err = header.Hash()
		if err != nil {
			return fmt.Errorf("hash beacon block header: %v", err)
		}
	}

	return nil
}

// verifyBeaconBlockBody verifies that the provided beacon block body matches the beacon block header's body root and the transaction receipt's block hash
func verifyBeaconBlockBody(beaconBlockHeader types.BeaconBlockHeader, beaconBlockBodySSZ []byte, txReceipt gethtypes.Receipt) error {
	var beaconBlockBodyRoot [32]byte
	beaconBlockBodyRootSlice := beaconBlockHeader.GetBodyRoot()
	copy(beaconBlockBodyRoot[:], beaconBlockBodyRootSlice)

	var beaconBlockBody prysmtypes.BeaconBlockBody
	if err := beaconBlockBody.UnmarshalSSZ(beaconBlockBodySSZ); err != nil {
		return fmt.Errorf("unmarshal beacon block body: %v", err)
	}

	beaconBlockBodyHash, err := beaconBlockBody.HashTreeRoot()
	if err != nil {
		return fmt.Errorf("hash beacon block body: %v", err)
	}
	if beaconBlockBodyHash != beaconBlockBodyRoot {
		return types.ErrBlockBodyMismatch
	}

	blockHash := common.BytesToHash(beaconBlockBody.GetEth1Data().GetBlockHash())
	if blockHash != txReceipt.BlockHash {
		return types.ErrBlockHashMismatch
	}

	return nil
}

// getClientStateBeaconBlockHeader extracts the beacon block header from the client state
func getClientStateBeaconBlockHeader(clientState ibccore.ClientState) (*types.BeaconBlockHeader, error) {
	clientStateBz, err := proto.Marshal(clientState)
	if err != nil {
		return nil, fmt.Errorf("marshal client state: %v", err)
	}
	wasmClientState := new(wasmtypes.ClientState)
	if err := proto.Unmarshal(clientStateBz, wasmClientState); err != nil {
		return nil, fmt.Errorf("unmarshal client state bytes: %v", err)
	}
	typedClientState := new(codectypes.Any)
	if err := typedClientState.Unmarshal(wasmClientState.Data); err != nil {
		return nil, fmt.Errorf("unmarshal typed eth client state bytes: %v", err)
	}
	ethClientState := new(types.ClientState)
	if err := ethClientState.Unmarshal(typedClientState.Value); err != nil {
		return nil, fmt.Errorf("unmarshal eth client state bytes: %v", err)
	}

	return ethClientState.GetInner().GetFinalizedHeader(), nil
}

// verifyTransferEvent verifies that the ERC-20 transfer log emitted in the transaction receipt includes the correct destination address and amount specified in the transfer intent
func verifyTransferEvent(txReceipt gethtypes.Receipt, intent types.TransferIntent, solverPublicKey *ecdsa.PublicKey) error {
	type LogTransfer struct {
		From         common.Address
		To           common.Address
		Tokens       *big.Int
		TokenAddress common.Address
	}
	transferEventSig := []byte("Transfer(address,address,uint256)")
	transferEventSigHash := crypto.Keccak256Hash(transferEventSig)

	var transferEvent LogTransfer
	for _, log := range txReceipt.Logs {
		if log.Topics[0].Hex() == transferEventSigHash.Hex() {
			transferEvent.From = common.HexToAddress(log.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(log.Topics[2].Hex())
			transferEvent.Tokens = new(big.Int).SetBytes(log.Data)
			transferEvent.TokenAddress = log.Address
			break
		}
	}

	if transferEvent == (LogTransfer{}) {
		return types.ErrTransferEventNotFound
	}
	if transferEvent.TokenAddress != common.HexToAddress(intent.TransferTokens.Erc20Address) {
		return types.ErrTokenAddressMismatch
	}
	if transferEvent.To != common.HexToAddress(intent.DestinationAddress) {
		return types.ErrDestinationAddressMismatch
	}
	solverAddress := crypto.PubkeyToAddress(*solverPublicKey).Hex()
	if transferEvent.From != common.HexToAddress(solverAddress) {
		return types.ErrSourceAddressMismatch
	}
	if transferEvent.Tokens.Cmp(intent.TransferTokens.Amount.BigInt()) != 0 {
		return types.ErrAmountMismatch
	}

	return nil
}

// validateClientState checks if the client state exists and is active
func (k Keeper) validateClientState(ctx sdk.Context, clientId string) error {
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

// getClientState retrieves the client state if it exists and validates that it is active
func (k Keeper) getClientState(ctx sdk.Context, clientId string) (ibccore.ClientState, error) {
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
