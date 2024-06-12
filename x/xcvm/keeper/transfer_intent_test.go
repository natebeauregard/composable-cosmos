package keeper_test

import (
	"crypto/ecdsa"
	"encoding/binary"
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
}

func (suite *TransferIntentTestSuite) SetupTest() {
	suite.app = helpers.SetupComposableAppWithValSet(suite.T())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "centauri-1", Time: time.Now().UTC()})

	// get wasm ethereum light client
	wasmBytes, err := os.ReadFile("testdata/icsxx_ethereum_cw.wasm")
	suite.Require().NoError(err)
	govAddress := sdk.AccAddress(address.Module("gov" /* github.com/cosmos/cosmos-sdk/x/gov.ModuleName */)).String()
	resp, err := suite.app.Wasm08Keeper.PushNewWasmCode(suite.ctx, &wasmtypes.MsgPushNewWasmCode{
		Signer: govAddress,
		Code:   wasmBytes,
	})
	suite.Require().NoError(err)
	suite.wasmCodeId = resp.CodeId
	suite.moduleAddr = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName).GetAddress()
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

	ethClientState := &types.ClientState{}
	ethClientId := createTestClient(suite, ethClientState)

	const tokenAmount uint64 = 10000
	tokens := types.TransferTokens{
		Erc20Address: tokenAddress,
		Amount:       sdk.NewUint(tokenAmount),
	}

	const bountyDenom = "stake"
	bountyAmount := sdk.NewInt(1000)

	// Send transfer intent message from user
	msgSendTransferIntent := types.MsgSendTransferIntent{
		FromAddress:        userAddress.String(),
		DestinationAddress: destinationEthAddress,
		ClientId:           ethClientId,
		TimeoutHeight:      suite.ctx.BlockHeight() + 100,
		TransferTokens:     &tokens,
		Bounty:             sdk.NewCoin(bountyDenom, bountyAmount),
	}
	_, err := suite.app.XCvmKeeper.SendTransferIntent(suite.ctx, &msgSendTransferIntent)
	suite.Require().NoError(err)

	// Verify the transfer intent is stored properly in the store
	intentId := uint64(0)
	transferIntent, err := suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, intentId)
	expectedTransferIntent := types.TransferIntent{
		SourceAddress:      msgSendTransferIntent.FromAddress,
		DestinationAddress: msgSendTransferIntent.DestinationAddress,
		TransferTokens:     msgSendTransferIntent.TransferTokens,
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

	solverPrivateKey, err := crypto.GenerateKey()
	suite.Require().NoError(err)
	publicKey, ok := solverPrivateKey.Public().(*ecdsa.PublicKey)
	suite.Require().True(ok)
	solverEthAddress := crypto.PubkeyToAddress(*publicKey).Hex()
	publicKeyCompressed := crypto.CompressPubkey(publicKey)

	// construct beacon block body and header to use for the light client state
	beaconBlockBody := &prysmtypes.BeaconBlockBody{
		RandaoReveal: make([]byte, 96),
		Eth1Data: &prysmtypes.Eth1Data{
			DepositRoot:  make([]byte, 32),
			DepositCount: 0,
			BlockHash:    common.FromHex(blockHash),
		},
		Graffiti: make([]byte, 32),
	}
	beaconBlockBodyBz, err := beaconBlockBody.MarshalSSZ()
	suite.Require().NoError(err)
	bodyRoot, err := beaconBlockBody.HashTreeRoot()
	suite.Require().NoError(err)

	ethClientState := &types.ClientState{
		Inner: &types.LightClientState{
			FinalizedHeader: &types.BeaconBlockHeader{
				BodyRoot: bodyRoot[:],
			},
		},
	}
	ethClientId := createTestClient(suite, ethClientState)

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
	}
	suite.app.XCvmKeeper.AddTransferIntent(suite.ctx, transferIntent, intentId)

	// create ERC20 transfer event log to use for intent proof
	tokenAmountBz := make([]byte, 8)
	binary.BigEndian.PutUint64(tokenAmountBz, tokenAmount)
	logs := []*gethtypes.Log{
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
	txReceipt := gethtypes.Receipt{
		Logs:      logs,
		BlockHash: common.HexToHash(blockHash),
	}
	txReceiptBz, err := txReceipt.MarshalJSON() // TODO: investigate why MarshalJSON needs to be used for passing into the Msg and why MarshalBinary loses some receipt information when unmarshalling
	suite.Require().NoError(err)
	txReceiptBinary, err := txReceipt.MarshalBinary()
	suite.Require().NoError(err)

	// create receipt proof to verify provided receipt was included in the block
	txReceiptHash := crypto.Keccak256(txReceiptBinary)
	receiptProof := types.ReceiptProof{
		Proof: make(map[string][]byte),
	}
	receiptTrie := trie.NewEmpty(trie.NewDatabase(rawdb.NewMemoryDatabase()))
	receiptTrie.Update(txReceiptHash[:], txReceiptBinary)
	err = receiptTrie.Prove(txReceiptHash[:], 0, receiptProof)
	suite.Require().NoError(err)
	receiptProofBz, err := receiptProof.Marshal()
	suite.Require().NoError(err)

	receiptHash := receiptTrie.Hash()
	blockHeader := &gethtypes.Header{
		ReceiptHash: receiptHash,
	}
	blockHeaderBz, err := rlp.EncodeToBytes(blockHeader)
	suite.Require().NoError(err)

	receiptData := append(txReceiptHash, common.FromHex(blockHash)...)
	receiptDataHash := crypto.Keccak256(receiptData)
	receiptSignature, err := crypto.Sign(receiptDataHash, solverPrivateKey)
	suite.Require().NoError(err)

	// create Msg to verify intent execution proof
	msgVerifyTransferIntentProof := types.MsgVerifyTransferIntentProof{
		Signer:           solverAddress.String(),
		IntentId:         intentId,
		TxReceipt:        txReceiptBz,
		ReceiptSignature: receiptSignature,
		PublicKey:        publicKeyCompressed,
		BlockHeader:      blockHeaderBz,
		ReceiptProof:     receiptProofBz,
		BeaconBlockBody:  beaconBlockBodyBz,
	}

	// Assert intent was verified correctly
	_, err = suite.app.XCvmKeeper.VerifyTransferIntentProof(suite.ctx, &msgVerifyTransferIntentProof)
	suite.Require().NoError(err)

	// Verify that the bounty was transferred to the solver's account and deducted from the module's account
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.moduleAddr, bountyDenom)
	suite.Require().Equal(sdk.NewInt(0), moduleBalance.Amount)
	solverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, solverAddress, bountyDenom)
	suite.Require().Equal(startingBaseAccountBalance.Add(bountyAmount), solverBalance.Amount)

	// Assert that transfer intent is purged from the store after being executed
	_, err = suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, intentId)
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

	ethClientState := &types.ClientState{}
	ethClientId := createTestClient(suite, ethClientState)

	const bountyDenom = "stake"
	bountyAmount := sdk.NewInt(1000)

	// Send transfer intent message from user
	const intentBlockDuration int64 = 100
	msgSendTransferIntent := types.MsgSendTransferIntent{
		FromAddress:        userAddress.String(),
		DestinationAddress: destinationEthAddress,
		ClientId:           ethClientId,
		TimeoutHeight:      suite.ctx.BlockHeight() + intentBlockDuration,
		TransferTokens:     &types.TransferTokens{},
		Bounty:             sdk.NewCoin(bountyDenom, bountyAmount),
	}
	_, err := suite.app.XCvmKeeper.SendTransferIntent(suite.ctx, &msgSendTransferIntent)
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
	_, err = suite.app.XCvmKeeper.TriggerTransferIntentTimeout(suite.ctx, &msgTriggerTransferIntentTimeout)
	suite.Require().ErrorContains(err, types.ErrPrematureTimeoutTrigger.Error())

	// Let the desired amount of time pass
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + intentBlockDuration)

	// Trigger the transfer intent timeout with a different user
	msgTriggerTransferIntentTimeout.Sender = otherAddress.String()
	_, err = suite.app.XCvmKeeper.TriggerTransferIntentTimeout(suite.ctx, &msgTriggerTransferIntentTimeout)
	suite.Require().ErrorContains(err, types.ErrInvalidSenderAddress.Error())

	// Trigger the transfer intent timeout after the desired time has passed
	msgTriggerTransferIntentTimeout.Sender = userAddress.String()
	_, err = suite.app.XCvmKeeper.TriggerTransferIntentTimeout(suite.ctx, &msgTriggerTransferIntentTimeout)
	suite.Require().NoError(err)

	// Verify that the bounty was transferred back to the user's account and deducted from the module's account
	moduleBalance = suite.app.BankKeeper.GetBalance(suite.ctx, suite.moduleAddr, bountyDenom)
	suite.Require().Equal(sdk.NewInt(0), moduleBalance.Amount)
	userBalance = suite.app.BankKeeper.GetBalance(suite.ctx, userAddress, bountyDenom)
	suite.Require().Equal(startingUserBalance, userBalance.Amount)

	// Assert that transfer intent is purged from the store after being executed
	_, err = suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, intentId)
	suite.Require().Error(err)
}

// Creates a test client with the provided Ethereum client state and returns the created client's ID
func createTestClient(suite *TransferIntentTestSuite, ethClientState *types.ClientState) string {
	ethClientStateBz, err := ethClientState.Marshal()
	suite.Require().NoError(err)

	clientState := wasmtypes.NewClientState(ethClientStateBz, suite.wasmCodeId, clienttypes.NewHeight(0, 0))
	consensusState := &wasmtypes.ConsensusState{}

	ethClientId, err := suite.app.IBCKeeper.ClientKeeper.CreateClient(suite.ctx, clientState, consensusState)
	suite.Require().NoError(err)

	return ethClientId
}
