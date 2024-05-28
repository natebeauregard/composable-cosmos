package keeper_test

import (
	"fmt"
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

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.ComposableApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = helpers.SetupComposableAppWithValSet(suite.T())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "centauri-1", Time: time.Now().UTC()})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// E2E tests with mocked solver behavior for sending and verifying transfer intents
func (suite *KeeperTestSuite) TestTransferIntent() {
	suite.SetupTest()

	// generate user and solver accounts
	accounts := app.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(10000000))
	userAddress := accounts[0]

	const ethDesinationAddress = "0xe6D38aEa101B30C7c26e533A7F7Dd22b82D1467d"

	// add eth-client as valid light client
	// TODO: create set client/consensus state?
	ethClientId, err := suite.app.IBCKeeper.ClientKeeper.CreateClient(suite.ctx, nil, nil)
	if err != nil {
		suite.Fail(fmt.Sprintf("failed to create eth client: %s", err))
	}
	intentId := uint64(0)

	// Send transfer intent message from user
	msgSendTransferIntent := types.MsgSendTransferIntent{
		FromAddress:        userAddress.String(),
		DestinationAddress: ethDesinationAddress,
		ClientId:           ethClientId,
		Amount:             sdk.NewUint(10000),
	}
	if err := suite.app.XCvmKeeper.SendEthTransferIntent(suite.ctx, &msgSendTransferIntent); err != nil {
		suite.Fail(fmt.Sprintf("failed to send transfer intent: %s", err))
	}

	// Verify the transfer intent is stored properly in the store
	transferIntent, err := suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, intentId)
	expectedTransferIntent := types.TransferIntent{
		SourceAddress:      msgSendTransferIntent.FromAddress,
		DestinationAddress: msgSendTransferIntent.DestinationAddress,
		Amount:             msgSendTransferIntent.Amount,
		ClientId:           msgSendTransferIntent.ClientId,
	}
	if err != nil {
		suite.Fail("transfer intent does not match expected")
	}
	suite.Equal(expectedTransferIntent, transferIntent)

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
	suite.Equal(expectedTransferIntentEvent, events[len(events)-1])

	// Solver listens to the event, posts collateral and executes the intent

	// Solver sends a message to verify the intent was properly executed

}

// var (
// 	successfulMsgSendTransferIntent = types.MsgSendTransferIntent{
// 		FromAddress:        "source",
// 		DestinationAddress: "destination",
// 		ClientId:           "client",
// 		Amount:             sdk.NewUint(10),
// 	}
// 	successfulIntent = types.TransferIntent{
// 		SourceAddress:      "source",
// 		DestinationAddress: "destination",
// 		Amount:             sdk.NewUint(10),
// 		ClientId:           "client",
// 	}
// 	successfulIntentEvent = sdk.Event{
// 		Type: "transfer_intent",
// 		// TODO: attributes
// 	}
// )

// SendValidatedTransferIntent tests
// func (suite *KeeperTestSuite) TestSendValidatedTransferIntent() {
// 	for _, tc := range []struct {
// 		desc                string
// 		expectedIntent      types.TransferIntent
// 		expectedIntentEvent sdk.Event
// 		malleate            func() error
// 		shouldErr           bool
// 		expectedErr         string
// 	}{
// 		{
// 			desc:                "Case success",
// 			expectedIntent:      successfulIntent,
// 			expectedIntentEvent: successfulIntentEvent,
// 			malleate: func() error {
// 				return suite.app.XCvmKeeper.SendValidatedTransferIntent(suite.ctx, &successfulMsgSendTransferIntent)
// 			},
// 			shouldErr: false,
// 		},
// 	} {
// 		suite.Run(tc.desc, func() {
// 			suite.SetupTest()
// 			err := tc.malleate()
// 			if !tc.shouldErr {
// 				res, _ := suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, 0)
// 				suite.Equal(res, tc.expectedIntent)
// 				// suite.ctx.EventManager().ABCIEvents()
// 			} else {
// 				suite.Equal(err.Error(), tc.expectedErr)
// 			}
// 		})
// 	}
// }
