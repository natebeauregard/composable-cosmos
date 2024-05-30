package keeper_test

import (
	"encoding/binary"
	"github.com/CosmWasm/wasmd/x/wasm/keeper/wasmtesting"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/notional-labs/composable/v6/x/xcvm/keeper"
	"strconv"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	wasmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/notional-labs/composable/v6/app"
	"github.com/notional-labs/composable/v6/app/helpers"
	testhelper "github.com/notional-labs/composable/v6/x/xcvm/testing"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
	prysmtypes "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	"github.com/stretchr/testify/suite"
)

type TransferIntentTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.ComposableApp

	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	mockVM      *wasmtesting.MockWasmEngine
}

func (suite *TransferIntentTestSuite) SetupTest() {
	suite.app = helpers.SetupComposableAppWithValSet(suite.T())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "centauri-1", Time: time.Now().UTC()})

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 1)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	//suite.mockVM = wasmtesting.NewIBCContractMockWasmEngine()
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TransferIntentTestSuite))
}

// E2E tests with mocked solver behavior for sending and verifying transfer intents
func (suite *TransferIntentTestSuite) TestTransferIntent() {
	suite.SetupTest()

	// generate user and solver accounts
	accounts := app.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(10000000))
	userAddress := accounts[0]
	solverAddress := accounts[1]

	const destinationEthAddress string = "0xe6D38aEa101B30C7c26e533A7F7Dd22b82D1467d"
	const solverEthAddress string = "0x02DB85F48Ffcf5F5Ea1fCF078eb5ABf468e53fAb"
	const blockHash string = "0x3f07a9c83155594c000642e7d60e8a8a00038d03e9849171a05ed0e2d47acbb3"
	const receiptHash string = "0x633a90413361fe1889d1e5ab4cb222608d224c458b30289b8390496a3fab29d8"

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
	ethClientStateBz, err := ethClientState.Marshal()
	suite.Require().NoError(err)
	clientState := wasmtypes.NewClientState(ethClientStateBz, testhelper.Code, clienttypes.NewHeight(0, 0))
	consensusState := &wasmtypes.ConsensusState{}
	ethClientId, err := suite.app.IBCKeeper.ClientKeeper.CreateClient(suite.ctx, clientState, consensusState)
	suite.Require().NoError(err)

	// Send transfer intent message from user
	const tokenAmount uint64 = 10000
	msgSendTransferIntent := types.MsgSendTransferIntent{
		FromAddress:        userAddress.String(),
		DestinationAddress: destinationEthAddress,
		ClientId:           ethClientId,
		Amount:             sdk.NewUint(tokenAmount),
	}
	_, err = suite.app.XCvmKeeper.SendTransferIntent(suite.ctx, &msgSendTransferIntent)
	suite.Require().NoError(err)

	// Verify the transfer intent is stored properly in the store
	intentId := uint64(0)
	transferIntent, err := suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, intentId)
	expectedTransferIntent := types.TransferIntent{
		SourceAddress:      msgSendTransferIntent.FromAddress,
		DestinationAddress: msgSendTransferIntent.DestinationAddress,
		Amount:             msgSendTransferIntent.Amount,
		ClientId:           msgSendTransferIntent.ClientId,
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expectedTransferIntent, transferIntent)

	// Verify the correct event is emitted
	expectedTransferIntentEvent := sdk.NewEvent(
		types.EventAddTransferIntent,
		sdk.NewAttribute(types.AttributeKeyIntentId, strconv.FormatUint(intentId, 10)),
		sdk.NewAttribute(types.AttributeKeyClientId, transferIntent.ClientId),
		sdk.NewAttribute(types.AttributeKeySourceAddress, transferIntent.SourceAddress),
		sdk.NewAttribute(types.AttributeKeyDestinationAddress, transferIntent.DestinationAddress),
		sdk.NewAttribute(types.AttributeKeyAmount, transferIntent.Amount.String()),
	)
	events := suite.ctx.EventManager().Events()
	suite.Require().Equal(expectedTransferIntentEvent, events[len(events)-1])

	// Solver listens to the event, posts collateral and executes the intent
	// TODO: add collateral logic once implemented

	// TODO: update client/consensus state with solver behavior

	// Solver sends a message to verify the intent was properly executed
	tokenAmountBz := make([]byte, 8)
	binary.BigEndian.PutUint64(tokenAmountBz, tokenAmount)
	logs := []*gethtypes.Log{
		{
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
	txReceiptBz, err := txReceipt.MarshalBinary()
	suite.Require().NoError(err)

	receiptProof := make(keeper.ReceiptProof)
	receiptProof[common.HexToHash(receiptHash)] = txReceiptBz
	receiptProofBz, err := rlp.EncodeToBytes(receiptProof)
	suite.Require().NoError(err)

	blockHeader := &gethtypes.Header{
		ReceiptHash: common.HexToHash(receiptHash),
	}
	blockHeaderBz, err := rlp.EncodeToBytes(blockHeader)
	suite.Require().NoError(err)

	msgVerifyTransferIntentProof := types.MsgVerifyTransferIntentProof{
		Signer:           solverAddress.String(),
		IntentId:         intentId,
		TxReceipt:        txReceiptBz,
		ReceiptSignature: []byte(solverEthAddress),
		BlockHeader:      blockHeaderBz,
		ReceiptProof:     receiptProofBz,
		BeaconBlockBody:  beaconBlockBodyBz,
	}

	// Assert proof was verified correctly
	_, err = suite.app.XCvmKeeper.VerifyTransferIntentProof(suite.ctx, &msgVerifyTransferIntentProof)
	suite.Require().NoError(err)

	// Assert that transfer intent is purged from the store after being executed
	_, err = suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, intentId)
	suite.Require().Error(err)
}
