package keeper_test

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	wasmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	prysmtypes "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	"os"
	"strconv"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/app"
	"github.com/notional-labs/composable/v6/app/helpers"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
	"github.com/stretchr/testify/suite"
)

type TransferIntentTestSuite struct {
	suite.Suite

	ctx        sdk.Context
	app        *app.ComposableApp
	wasmCodeId []byte
	moduleAddr sdk.AccAddress
	govAddr    sdk.AccAddress
}

func (suite *TransferIntentTestSuite) SetupTest() {
	suite.app = helpers.SetupComposableAppWithValSet(suite.T())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "centauri-1", Time: time.Now().UTC()})

	// get wasm ethereum light client
	wasmBytes, err := os.ReadFile("testdata/icsxx_ethereum_cw.wasm")
	suite.Require().NoError(err)
	govAddress := sdk.AccAddress(address.Module("gov" /* github.com/cosmos/cosmos-sdk/x/gov.ModuleName */))
	resp, err := suite.app.Wasm08Keeper.PushNewWasmCode(suite.ctx, &wasmtypes.MsgPushNewWasmCode{
		Signer: govAddress.String(),
		Code:   wasmBytes,
	})
	suite.Require().NoError(err)
	suite.wasmCodeId = resp.CodeId
	suite.moduleAddr = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName).GetAddress()
	suite.govAddr = govAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TransferIntentTestSuite))
}

// Unit test for sending transfer intents
func (suite *TransferIntentTestSuite) TestSendTransferIntent() {
	suite.SetupTest()

	// generate user account
	startingUserBalance := sdk.NewInt(10000000)
	userAddress := app.AddTestAddrs(suite.app, suite.ctx, 1, startingUserBalance)[0]

	const destinationEthAddress string = "0xe6D38aEa101B30C7c26e533A7F7Dd22b82D1467d"
	const tokenAddress string = "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326"

	_, ethClientId, err := createEthLightClient(suite, make([]byte, 32))

	const tokenAmount uint64 = 10000
	tokens := types.TransferTokens{
		Erc20Address: tokenAddress,
		Amount:       sdk.NewUint(tokenAmount),
	}

	const bountyDenom = "stake"
	bountyAmount := sdk.NewInt(1000)

	// Send transfer intent message from user
	msgSendTransferIntent := types.MsgSendTransferIntent{
		Sender:             userAddress.String(),
		DestinationAddress: destinationEthAddress,
		ClientId:           ethClientId,
		TimeoutHeight:      suite.ctx.BlockHeight() + 100,
		TransferTokens:     &tokens,
		Bounty:             sdk.NewCoin(bountyDenom, bountyAmount),
	}
	_, err = suite.app.XCVMKeeper.SendTransferIntent(suite.ctx, &msgSendTransferIntent)
	suite.Require().NoError(err)

	// Verify the transfer intent is stored properly in the store
	intentId := uint64(0)
	transferIntent, err := suite.app.XCVMKeeper.GetTransferIntent(suite.ctx, intentId)
	expectedTransferIntent := types.TransferIntent{
		SourceAddress:      msgSendTransferIntent.Sender,
		DestinationAddress: msgSendTransferIntent.DestinationAddress,
		TransferTokens:     msgSendTransferIntent.TransferTokens,
		StartingHeight:     suite.ctx.BlockHeight(),
		TimeoutHeight:      msgSendTransferIntent.TimeoutHeight,
		ClientId:           msgSendTransferIntent.ClientId,
		Bounty:             msgSendTransferIntent.Bounty,
	}
	suite.Require().NoError(err)
	suite.Require().Equal(&expectedTransferIntent, transferIntent)

	// Verify the correct event is emitted
	expectedTransferIntentEvent := sdk.NewEvent(
		types.EventAddTransferIntent,
		sdk.NewAttribute(types.AttributeKeyIntentId, strconv.FormatUint(intentId, 10)),
		sdk.NewAttribute(types.AttributeKeyClientId, transferIntent.ClientId),
		sdk.NewAttribute(types.AttributeKeySourceAddress, transferIntent.SourceAddress),
		sdk.NewAttribute(types.AttributeKeyDestinationAddress, transferIntent.DestinationAddress),
		sdk.NewAttribute(types.AttributeKeyTimeout, strconv.FormatInt(transferIntent.TimeoutHeight, 10)),
		sdk.NewAttribute(types.AttributeKeyAmount, transferIntent.TransferTokens.String()),
		sdk.NewAttribute(types.AttributeKeyBounty, transferIntent.Bounty.String()),
	)
	events := suite.ctx.EventManager().Events()
	suite.Require().Equal(expectedTransferIntentEvent, events[len(events)-1])

	// Verify that the bounty was deducted from the user's account and is stored in the XCVM module account
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.moduleAddr, bountyDenom)
	suite.Require().Equal(bountyAmount, moduleBalance.Amount)
	userBalance := suite.app.BankKeeper.GetBalance(suite.ctx, userAddress, bountyDenom)
	suite.Require().Equal(startingUserBalance.Sub(bountyAmount), userBalance.Amount)
}

// Unit test for verifying transfer intent proofs
func (suite *TransferIntentTestSuite) TestVerifyTransferIntentProof() {
	suite.SetupTest()

	// generate user and solver accounts
	startingBaseAccountBalance := sdk.NewInt(10000000)
	accounts := app.AddTestAddrs(suite.app, suite.ctx, 2, startingBaseAccountBalance)
	userAddress := accounts[0]
	solverAddress := accounts[1]

	// setup transfer intent bounty
	const bountyDenom = "stake"
	bountyAmount := sdk.NewInt(1000)
	bounty := sdk.NewCoin(bountyDenom, bountyAmount)
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, userAddress, types.ModuleName, sdk.NewCoins(bounty))
	suite.Require().NoError(err)

	const destinationEthAddress string = "0xe6D38aEa101B30C7c26e533A7F7Dd22b82D1467d"
	const tokenAddress string = "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326"
	const blockHash string = "0x3f07a9c83155594c000642e7d60e8a8a00038d03e9849171a05ed0e2d47acbb3"

	solverPrivateKey, publicKeyCompressed, solverEthAddress := generateKeys()

	// construct beacon block body and header to use for the light client state
	beaconBlockBody := createBeaconBlockBody(blockHash)
	beaconBlockBodyBz, err := beaconBlockBody.MarshalSSZ()
	suite.Require().NoError(err)
	bodyRoot, err := beaconBlockBody.HashTreeRoot()
	suite.Require().NoError(err)

	_, ethClientId, err := createEthLightClient(suite, bodyRoot[:])
	suite.Require().NoError(err)

	const tokenAmount uint64 = 10000
	tokens := types.TransferTokens{
		Erc20Address: tokenAddress,
		Amount:       sdk.NewUint(tokenAmount),
	}

	// Add transfer intent to the XCVM module store
	const intentId uint64 = 0
	transferIntent := types.TransferIntent{
		SourceAddress:      userAddress.String(),
		DestinationAddress: destinationEthAddress,
		ClientId:           ethClientId,
		TransferTokens:     &tokens,
		Bounty:             bounty,
		TimeoutHeight:      suite.ctx.BlockHeight() + 100,
	}
	err = suite.app.XCVMKeeper.AddTransferIntent(suite.ctx, transferIntent, intentId)
	suite.Require().NoError(err)

	// create ERC20 transfer event log to use for intent proof
	tokenAmountBz := make([]byte, 8)
	binary.BigEndian.PutUint64(tokenAmountBz, tokenAmount)
	logs := createERC20Logs(tokenAddress, solverEthAddress, destinationEthAddress, tokenAmountBz)
	txReceipt := gethtypes.Receipt{
		Logs:      logs,
		BlockHash: common.HexToHash(blockHash),
	}
	txReceiptBz, err := txReceipt.MarshalJSON()
	suite.Require().NoError(err)
	txReceiptBinary, err := txReceipt.MarshalBinary()
	suite.Require().NoError(err)

	// create receipt proof to verify provided receipt was included in the block
	txReceiptHash := crypto.Keccak256(txReceiptBinary)
	receiptProofBz, receiptHash, err := createReceiptProof(txReceiptHash, txReceiptBinary)

	blockHeader := &gethtypes.Header{
		ReceiptHash: *receiptHash,
	}
	blockHeaderBz, err := rlp.EncodeToBytes(blockHeader)
	suite.Require().NoError(err)

	receiptData := append(txReceiptHash, common.FromHex(blockHash)...)
	receiptDataHash := crypto.Keccak256(receiptData)
	receiptSignature, err := crypto.Sign(receiptDataHash, solverPrivateKey)
	suite.Require().NoError(err)

	// create Msg to verify intent execution proof
	msgVerifyTransferIntentProof := types.MsgVerifyTransferIntentProof{
		Sender:           solverAddress.String(),
		IntentId:         intentId,
		TxReceipt:        txReceiptBz,
		ReceiptSignature: receiptSignature,
		PublicKey:        publicKeyCompressed,
		BlockHeader:      blockHeaderBz,
		ReceiptProof:     receiptProofBz,
		BeaconBlockBody:  beaconBlockBodyBz,
	}

	// Assert intent was verified correctly
	_, err = suite.app.XCVMKeeper.VerifyTransferIntentProof(suite.ctx, &msgVerifyTransferIntentProof)
	suite.Require().NoError(err)

	// Verify that the bounty was transferred to the solver's account and deducted from the module's account
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.moduleAddr, bountyDenom)
	suite.Require().Equal(sdk.NewInt(0), moduleBalance.Amount)
	solverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, solverAddress, bountyDenom)
	suite.Require().Equal(startingBaseAccountBalance.Add(bountyAmount), solverBalance.Amount)

	// Assert that transfer intent is purged from the store after being executed
	_, err = suite.app.XCVMKeeper.GetTransferIntent(suite.ctx, intentId)
	suite.Require().Error(err)
}

// Unit test for sending transfer intents
func (suite *TransferIntentTestSuite) TestTriggerTransferIntentTimeout() {
	suite.SetupTest()

	// generate user account
	startingUserBalance := sdk.NewInt(10000000)
	accounts := app.AddTestAddrs(suite.app, suite.ctx, 2, startingUserBalance)
	userAddress := accounts[0]
	otherAddress := accounts[1]

	const destinationEthAddress string = "0xe6D38aEa101B30C7c26e533A7F7Dd22b82D1467d"

	_, ethClientId, err := createEthLightClient(suite, make([]byte, 32))
	suite.Require().NoError(err)

	const bountyDenom = "stake"
	bountyAmount := sdk.NewInt(1000)

	// Send transfer intent message from user
	const intentBlockDuration int64 = 100
	msgSendTransferIntent := types.MsgSendTransferIntent{
		Sender:             userAddress.String(),
		DestinationAddress: destinationEthAddress,
		ClientId:           ethClientId,
		TimeoutHeight:      suite.ctx.BlockHeight() + intentBlockDuration,
		TransferTokens:     &types.TransferTokens{},
		Bounty:             sdk.NewCoin(bountyDenom, bountyAmount),
	}
	_, err = suite.app.XCVMKeeper.SendTransferIntent(suite.ctx, &msgSendTransferIntent)
	suite.Require().NoError(err)

	// Verify that the bounty was deducted from the user's account and is stored in the XCVM module account
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.moduleAddr, bountyDenom)
	suite.Require().Equal(bountyAmount, moduleBalance.Amount)
	userBalance := suite.app.BankKeeper.GetBalance(suite.ctx, userAddress, bountyDenom)
	suite.Require().Equal(startingUserBalance.Sub(bountyAmount), userBalance.Amount)

	const intentId uint64 = 0
	msgTriggerTransferIntentTimeout := types.MsgTriggerTransferIntentTimeout{
		Sender:   userAddress.String(),
		IntentId: intentId,
	}

	// Trigger the transfer intent timeout prematurely
	_, err = suite.app.XCVMKeeper.TriggerTransferIntentTimeout(suite.ctx, &msgTriggerTransferIntentTimeout)
	suite.Require().ErrorContains(err, types.ErrPrematureTimeoutTrigger.Error())

	// Let the desired amount of time pass
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + intentBlockDuration)

	// Trigger the transfer intent timeout with a different user
	msgTriggerTransferIntentTimeout.Sender = otherAddress.String()
	_, err = suite.app.XCVMKeeper.TriggerTransferIntentTimeout(suite.ctx, &msgTriggerTransferIntentTimeout)
	suite.Require().ErrorContains(err, types.ErrInvalidSenderAddress.Error())

	// Trigger the transfer intent timeout after the desired time has passed
	msgTriggerTransferIntentTimeout.Sender = userAddress.String()
	_, err = suite.app.XCVMKeeper.TriggerTransferIntentTimeout(suite.ctx, &msgTriggerTransferIntentTimeout)
	suite.Require().NoError(err)

	// Verify that the bounty was transferred back to the user's account and deducted from the module's account
	moduleBalance = suite.app.BankKeeper.GetBalance(suite.ctx, suite.moduleAddr, bountyDenom)
	suite.Require().Equal(sdk.NewInt(0), moduleBalance.Amount)
	userBalance = suite.app.BankKeeper.GetBalance(suite.ctx, userAddress, bountyDenom)
	suite.Require().Equal(startingUserBalance, userBalance.Amount)

	// Assert that transfer intent is purged from the store after being executed
	_, err = suite.app.XCVMKeeper.GetTransferIntent(suite.ctx, intentId)
	suite.Require().Error(err)
}

func (suite *TransferIntentTestSuite) TestE2EIntent_IntentIncludedInCurrentClientState() {
	const executionBlockHash string = "0x3f07a9c83155594c000642e7d60e8a8a00038d03e9849171a05ed0e2d47acbb3"

	// construct beacon block body and header to use for the light client state
	intentBeaconBlockBody := createBeaconBlockBody(executionBlockHash)
	intentBodyRoot, err := intentBeaconBlockBody.HashTreeRoot()
	suite.Require().NoError(err)

	lightClientBeaconBlockHeader := &types.BeaconBlockHeader{
		Slot:          0,
		ProposerIndex: 0,
		ParentRoot:    make([]byte, 32),
		StateRoot:     make([]byte, 32),
		BodyRoot:      intentBodyRoot[:],
	}

	suite.runE2ETest(lightClientBeaconBlockHeader, nil, intentBeaconBlockBody)
}

func (suite *TransferIntentTestSuite) TestE2EIntent_IntentIncludedInPreviousClientState() {
	const executionBlockHash string = "0x3f07a9c83155594c000642e7d60e8a8a00038d03e9849171a05ed0e2d47acbb3"
	const intentBeaconBlockHash string = "0x89c3e4585e9904876cdcfbd4d6bfbd474b8e42c2eb754e19feb7e55e7d6c88eb"
	const intermediateBeaconBlockHash string = "0xb3a4ab7196f46f24c076f3575131c8abae9aeebf2615c14e73aab7962336fba8"

	// construct beacon block body and header to use for the light client state
	intentBeaconBlockBody := createBeaconBlockBody(executionBlockHash)
	intentBodyRoot, err := intentBeaconBlockBody.HashTreeRoot()
	suite.Require().NoError(err)

	lightClientBeaconBlockHeader := &types.BeaconBlockHeader{
		Slot:          0,
		ProposerIndex: 0,
		ParentRoot:    common.FromHex(intermediateBeaconBlockHash),
		StateRoot:     make([]byte, 32),
		BodyRoot:      make([]byte, 32),
	}

	// construct previous beacon block headers to verify an intent execution from a prior block
	intentBeaconBlockHeader := &types.BeaconBlockHeader{
		Slot:          0,
		ProposerIndex: 0,
		ParentRoot:    make([]byte, 32),
		StateRoot:     make([]byte, 32),
		BodyRoot:      intentBodyRoot[:],
	}
	intermediateBeaconBlockHeader := &types.BeaconBlockHeader{
		Slot:          0,
		ProposerIndex: 0,
		ParentRoot:    common.FromHex(intentBeaconBlockHash),
		StateRoot:     make([]byte, 32),
		BodyRoot:      make([]byte, 32),
	}
	previousBeaconBlockHeaders := []*types.BeaconBlockHeader{
		intentBeaconBlockHeader,
		intermediateBeaconBlockHeader,
	}

	suite.runE2ETest(lightClientBeaconBlockHeader, previousBeaconBlockHeaders, intentBeaconBlockBody)
}

func (suite *TransferIntentTestSuite) runE2ETest(lightClientBeaconBlockHeader *types.BeaconBlockHeader, previousBeaconBlockHeaders []*types.BeaconBlockHeader, intentBeaconBlockBody *prysmtypes.BeaconBlockBody) {
	suite.SetupTest()

	// generate user account
	startingBaseAccountBalance := sdk.NewInt(10000000)
	accounts := app.AddTestAddrs(suite.app, suite.ctx, 2, startingBaseAccountBalance)
	userAddress := accounts[0]
	solverAddress := accounts[1]

	const destinationEthAddress string = "0xe6D38aEa101B30C7c26e533A7F7Dd22b82D1467d"
	const tokenAddress string = "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326"
	const blockHash string = "0x3f07a9c83155594c000642e7d60e8a8a00038d03e9849171a05ed0e2d47acbb3"
	const startingBodyRoot string = "6f2d49f08a43eacd98c4f0ae09b7e6a98e784f1b7c5e4c6d1f0f9a6c4d5b2e3a"

	solverPrivateKey, publicKeyCompressed, solverEthAddress := generateKeys()

	originalBodyRoot, err := hex.DecodeString(startingBodyRoot)
	_, ethClientId, err := createEthLightClient(suite, originalBodyRoot)

	const tokenAmount uint64 = 10000
	tokens := types.TransferTokens{
		Erc20Address: tokenAddress,
		Amount:       sdk.NewUint(tokenAmount),
	}

	const bountyDenom = "stake"
	bountyAmount := sdk.NewInt(1000)

	// Send transfer intent message from user
	msgSendTransferIntent := types.MsgSendTransferIntent{
		Sender:             userAddress.String(),
		DestinationAddress: destinationEthAddress,
		ClientId:           ethClientId,
		TimeoutHeight:      suite.ctx.BlockHeight() + 100,
		TransferTokens:     &tokens,
		Bounty:             sdk.NewCoin(bountyDenom, bountyAmount),
	}
	fmt.Println("Transfer intent sent by user")

	_, err = suite.app.XCVMKeeper.SendTransferIntent(suite.ctx, &msgSendTransferIntent)
	suite.Require().NoError(err)

	lightClientUpdate := types.LightClientUpdate{
		AttestedHeader:       lightClientBeaconBlockHeader,
		XSyncCommitteeUpdate: nil,
		FinalizedHeader:      lightClientBeaconBlockHeader,
		ExecutionPayload: &types.ExecutionPayloadProof{
			StateRoot:              make([]byte, 32),
			BlockNumber:            1,
			MultiProof:             nil,
			ExecutionPayloadBranch: nil,
			Timestamp:              0,
		},
		FinalityProof: &types.FinalityProof{
			Epoch:          0,
			FinalityBranch: nil,
		},
		SyncAggregate: &types.SyncAggregate{
			SyncCommitteeBits:      make([]byte, 64),
			SyncCommitteeSignature: make([]byte, 96),
		},
		SignatureSlot: 0,
	}
	lightClientUpdateBz, err := lightClientUpdate.Marshal()
	suite.Require().NoError(err)
	wrappedLightClientUpdate := codectypes.Any{
		TypeUrl: "/ibc.lightclients.ethereum.v1.LightClientUpdate",
		Value:   lightClientUpdateBz,
	}
	wrappedLightClientUpdateBz, err := wrappedLightClientUpdate.Marshal()
	suite.Require().NoError(err)
	updateClientMsg, err := clienttypes.NewMsgUpdateClient(
		ethClientId,
		&wasmtypes.Header{
			Data:   wrappedLightClientUpdateBz,
			Height: clienttypes.NewHeight(0, 1),
		},
		suite.govAddr.String(),
	)
	fmt.Println("Updating light client with new solver execution")
	_, err = suite.app.IBCKeeper.UpdateClient(suite.ctx, updateClientMsg)
	suite.Require().NoError(err)

	// create ERC20 transfer event log to use for intent proof
	tokenAmountBz := make([]byte, 8)
	binary.BigEndian.PutUint64(tokenAmountBz, tokenAmount)
	logs := createERC20Logs(tokenAddress, solverEthAddress, destinationEthAddress, tokenAmountBz)
	txReceipt := gethtypes.Receipt{
		Logs:      logs,
		BlockHash: common.HexToHash(blockHash),
	}
	txReceiptBz, err := txReceipt.MarshalJSON()
	suite.Require().NoError(err)
	txReceiptBinary, err := txReceipt.MarshalBinary()
	suite.Require().NoError(err)

	// create receipt proof to verify provided receipt was included in the block
	txReceiptHash := crypto.Keccak256(txReceiptBinary)
	receiptProofBz, receiptHash, err := createReceiptProof(txReceiptHash, txReceiptBinary)

	blockHeader := &gethtypes.Header{
		ReceiptHash: *receiptHash,
	}
	blockHeaderBz, err := rlp.EncodeToBytes(blockHeader)
	suite.Require().NoError(err)

	receiptData := append(txReceiptHash, common.FromHex(blockHash)...)
	receiptDataHash := crypto.Keccak256(receiptData)
	receiptSignature, err := crypto.Sign(receiptDataHash, solverPrivateKey)
	suite.Require().NoError(err)

	intentBeaconBlockBodyBz, err := intentBeaconBlockBody.MarshalSSZ()
	suite.Require().NoError(err)

	// create Msg to verify intent execution proof
	msgVerifyTransferIntentProof := types.MsgVerifyTransferIntentProof{
		Sender:             solverAddress.String(),
		IntentId:           0,
		TxReceipt:          txReceiptBz,
		ReceiptSignature:   receiptSignature,
		PublicKey:          publicKeyCompressed,
		BlockHeader:        blockHeaderBz,
		ReceiptProof:       receiptProofBz,
		BeaconBlockBody:    intentBeaconBlockBodyBz,
		BeaconBlockHeaders: previousBeaconBlockHeaders,
	}

	// Assert intent was verified correctly
	fmt.Println("Verifying transfer intent proof submitted by solver")
	_, err = suite.app.XCVMKeeper.VerifyTransferIntentProof(suite.ctx, &msgVerifyTransferIntentProof)
	suite.Require().NoError(err)

	// Verify that the bounty was transferred to the solver's account and deducted from the user's account
	solverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, solverAddress, bountyDenom)
	suite.Require().Equal(startingBaseAccountBalance.Add(bountyAmount), solverBalance.Amount)
	userBalance := suite.app.BankKeeper.GetBalance(suite.ctx, userAddress, bountyDenom)
	suite.Require().Equal(startingBaseAccountBalance.Sub(bountyAmount), userBalance.Amount)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.moduleAddr, bountyDenom)
	suite.Require().Equal(sdk.NewInt(0), moduleBalance.Amount)

	// Assert that transfer intent is purged from the store after being executed
	_, err = suite.app.XCVMKeeper.GetTransferIntent(suite.ctx, 0)
	suite.Require().Error(err)
}

func createMockSyncCommittee(BlsPublicKeyLen uint32, SyncCommitteeMembers uint32) *types.SyncCommittee {
	var publicKeys [][]byte
	for i := uint32(0); i < SyncCommitteeMembers; i++ {
		publicKeys = append(publicKeys, make([]byte, BlsPublicKeyLen))
	}
	return &types.SyncCommittee{
		PublicKeys:         publicKeys,
		AggregatePublicKey: make([]byte, BlsPublicKeyLen),
	}
}

func generateKeys() (*ecdsa.PrivateKey, []byte, string) {
	privateKey, _ := crypto.GenerateKey()
	publicKey, _ := privateKey.Public().(*ecdsa.PublicKey)
	publicKeyEthAddress := crypto.PubkeyToAddress(*publicKey).Hex()
	publicKeyCompressed := crypto.CompressPubkey(publicKey)
	return privateKey, publicKeyCompressed, publicKeyEthAddress
}

func createBeaconBlockBody(blockHashHex string) *prysmtypes.BeaconBlockBody {
	return &prysmtypes.BeaconBlockBody{
		RandaoReveal: make([]byte, 96),
		Eth1Data: &prysmtypes.Eth1Data{
			DepositRoot:  make([]byte, 32),
			DepositCount: 0,
			BlockHash:    common.FromHex(blockHashHex),
		},
		Graffiti: make([]byte, 32),
	}
}

func createReceiptProof(txReceiptHash []byte, txReceiptBz []byte) ([]byte, *common.Hash, error) {
	otherTxReceipt := gethtypes.Receipt{}
	otherTxReceiptBz, err := otherTxReceipt.MarshalBinary()
	otherTxReceiptHash := crypto.Keccak256(txReceiptBz)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal placeholder tx receipt: %w", err)
	}

	receiptTrie := trie.NewEmpty(trie.NewDatabase(rawdb.NewMemoryDatabase()))
	receiptTrie.Update(txReceiptHash[:], txReceiptBz)
	receiptTrie.Update(otherTxReceiptHash[:], otherTxReceiptBz)

	receiptProof := types.ReceiptProof{
		Proof: make(map[string][]byte),
	}
	if err := receiptTrie.Prove(txReceiptHash[:], 0, receiptProof); err != nil {
		return nil, nil, err
	}
	receiptProofBz, err := receiptProof.Marshal()
	if err != nil {
		return nil, nil, err
	}

	receiptHash := receiptTrie.Hash()

	return receiptProofBz, &receiptHash, nil
}

// Creates a test client with the provided Ethereum client state and returns the created client's ID
func createEthLightClient(suite *TransferIntentTestSuite, bodyRoot []byte) (*types.ClientState, string, error) {
	syncCommittee := createMockSyncCommittee(48, 512)
	ethClientState := types.ClientState{
		Inner: &types.LightClientState{
			FinalizedHeader: &types.BeaconBlockHeader{
				Slot:          0,
				ProposerIndex: 0,
				ParentRoot:    make([]byte, 32),
				StateRoot:     make([]byte, 32),
				BodyRoot:      bodyRoot[:],
			},
			LatestFinalizedEpoch: 0,
			CurrentSyncCommittee: syncCommittee,
			NextSyncCommittee:    syncCommittee,
		},
	}
	ethClientStateBz, err := ethClientState.Marshal()
	if err != nil {
		return nil, "", err
	}
	const ethClientStateTypeUrl = "/ibc.lightclients.ethereum.v1.ClientState"
	wasmClientState := &codectypes.Any{
		TypeUrl: ethClientStateTypeUrl,
		Value:   ethClientStateBz,
	}
	wasmClientStateBz, err := wasmClientState.Marshal()
	if err != nil {
		return nil, "", err
	}

	clientState := wasmtypes.NewClientState(wasmClientStateBz, suite.wasmCodeId, clienttypes.NewHeight(0, 0))
	consensusState := &wasmtypes.ConsensusState{}

	ethClientId, err := suite.app.IBCKeeper.ClientKeeper.CreateClient(suite.ctx, clientState, consensusState)
	if err != nil {
		return nil, "", err
	}

	return &ethClientState, ethClientId, nil
}

func createERC20Logs(tokenAddress string, solverEthAddress string, destinationEthAddress string, tokenAmountBz []byte) []*gethtypes.Log {
	return []*gethtypes.Log{
		{
			Address: common.HexToAddress(tokenAddress),
			Topics: []common.Hash{
				crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")),
				common.HexToHash(solverEthAddress),
				common.HexToHash(destinationEthAddress),
			},
			Data: tokenAmountBz,
		},
	}
}
