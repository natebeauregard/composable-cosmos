package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/xcvm/types"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

func (k Keeper) executeTransferMsg(ctx sdk.Context, transferMsg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error) {
	if err := transferMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("bad msg %v", err.Error())
	}
	return k.transferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)
}

func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) {
	var fungibleTokenPacketData transfertypes.FungibleTokenPacketData

	err := transfertypes.ModuleCdc.UnmarshalJSON(packet.Data, &fungibleTokenPacketData)
	if err != nil {
		return
	}

	var cvmProgram types.Program
	err = json.Unmarshal([]byte(fungibleTokenPacketData.Memo), &cvmProgram)
	if err != nil {
		return
	}

	k.executeCvmProgram(cvmProgram, ctx, packet.SourceChannel, packet.DestinationPort, packet.DestinationChannel)
}
