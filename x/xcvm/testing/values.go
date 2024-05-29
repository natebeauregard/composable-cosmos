package testing

import (
	"errors"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	wasmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

var (
	// Represents the code of the wasm contract used in the tests with a mock vm.
	WasmMagicNumber                    = []byte("\x00\x61\x73\x6D")
	Code                               = append(WasmMagicNumber, []byte("0123456780123456780123456780")...)
	MockClientStateBz                  = []byte("client-state-data")
	MockConsensusStateBz               = []byte("consensus-state-data")
	MockTendermitClientState           = CreateMockTendermintClientState(clienttypes.NewHeight(1, 10))
	MockTendermintClientHeader         = &ibctm.Header{}
	MockTendermintClientMisbehaviour   = ibctm.NewMisbehaviour("client-id", MockTendermintClientHeader, MockTendermintClientHeader)
	MockTendermintClientConsensusState = ibctm.NewConsensusState(time.Now(), commitmenttypes.NewMerkleRoot([]byte("hash")), []byte("nextValsHash"))
	MockValidProofBz                   = []byte("valid proof")
	MockInvalidProofBz                 = []byte("invalid proof")
	MockUpgradedClientStateProofBz     = []byte("upgraded client state proof")
	MockUpgradedConsensusStateProofBz  = []byte("upgraded consensus state proof")

	ErrMockContract = errors.New("mock contract error")
	ErrMockVM       = errors.New("mock vm error")
)

// CreateMockTendermintClientState returns a valid Tendermint client state for use in tests.
func CreateMockTendermintClientState(height clienttypes.Height) *ibctm.ClientState {
	return ibctm.NewClientState(
		"chain-id",
		ibctm.DefaultTrustLevel,
		ibctesting.TrustingPeriod,
		ibctesting.UnbondingPeriod,
		ibctesting.MaxClockDrift,
		height,
		commitmenttypes.GetSDKSpecs(),
		ibctesting.UpgradePath,
	)
}

// CreateMockEthereumClientState returns a valid Ethereum client state for use in tests.
func CreateMockEthereumClientState(bodyRoot []byte) *types.ClientState {
	return &types.ClientState{
		Inner: &types.LightClientState{
			FinalizedHeader: &types.BeaconBlockHeader{
				Slot:          0,
				ProposerIndex: 0,
				ParentRoot:    nil,
				StateRoot:     nil,
				BodyRoot:      bodyRoot,
			},
			LatestFinalizedEpoch: 0,
			CurrentSyncCommittee: nil,
			NextSyncCommittee:    nil,
			StatePeriod:          0,
		},
		FrozenHeight:   nil,
		LatestHeight:   0,
		IbcCoreAddress: "",
		NextUpgradeId:  0,
		XPhantom:       nil,
	}
}

// CreateMockClientStateBz returns valid client state bytes for use in tests.
func CreateMockClientStateBz(cdc codec.BinaryCodec, checksum wasmvmtypes.Checksum) []byte {
	wrappedClientStateBz := clienttypes.MustMarshalClientState(cdc, MockTendermitClientState)
	mockClientSate := wasmtypes.NewClientState(wrappedClientStateBz, checksum, MockTendermitClientState.LatestHeight)
	return clienttypes.MustMarshalClientState(cdc, mockClientSate)
}

// CreateMockContract returns a well formed (magic number prefixed) wasm contract the given code.
func CreateMockContract(code []byte) []byte {
	return append(WasmMagicNumber, code...)
}
