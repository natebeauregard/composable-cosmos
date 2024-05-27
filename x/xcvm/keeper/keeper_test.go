package keeper_test

import (
	"testing"
	"time"

	"github.com/notional-labs/composable/v6/app"
	"github.com/notional-labs/composable/v6/app/helpers"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
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

var (
	successfulMsgSendTransferIntent = types.MsgSendTransferIntent{
		FromAddress:        "source",
		DestinationAddress: "destination",
		ClientId:           "client",
		Amount:             sdk.NewUint(10),
	}
	successfulIntent = types.TransferIntent{
		SourceAddress:      "source",
		DestinationAddress: "destination",
		Amount:             sdk.NewUint(10),
		ClientId:           "client",
	}
	successfulIntentEvent = sdk.Event{
		Type: "transfer_intent",
		// TODO: attributes
	}
)

func (suite *KeeperTestSuite) TestSendValidatedTransferIntent() {
	for _, tc := range []struct {
		desc                string
		expectedIntent      types.TransferIntent
		expectedIntentEvent sdk.Event
		malleate            func() error
		shouldErr           bool
		expectedErr         string
	}{
		{
			desc:                "Case success",
			expectedIntent:      successfulIntent,
			expectedIntentEvent: successfulIntentEvent,
			malleate: func() error {
				return suite.app.XCvmKeeper.SendValidatedTransferIntent(suite.ctx, &successfulMsgSendTransferIntent)
			},
			shouldErr: false,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			err := tc.malleate()
			if !tc.shouldErr {
				res, _ := suite.app.XCvmKeeper.GetTransferIntent(suite.ctx, 0)
				suite.Equal(res, tc.expectedIntent)
				// suite.ctx.EventManager().ABCIEvents()
			} else {
				suite.Equal(err.Error(), tc.expectedErr)
			}
		})
	}
}
