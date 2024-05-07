package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	types1 "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/notional-labs/composable/v6/x/xcvm/types"
)

// TODO: call Composable smart contract to execute CVM program instead?
func (k Keeper) executeCvmProgram(cvmProgram types.Program, ctx sdk.Context, from string, sourcePort string, sourceChannel string) error {

	// TODO: add multiple instruction executions
	//
	// instructionPointer := 0
	// for _, instruction := range cvmProgram.Instructions {
	// 	instructionPointer++

	err := k.executeInstruction(*cvmProgram.Instructions[0], ctx, from, sourcePort, sourceChannel)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) executeInstruction(instruction types.Instruction, ctx sdk.Context, from string, sourcePort string, sourceChannel string) error {
	// TODO: implement switch or more elegant way to handle instruction execution
	if spawn := instruction.GetSpawn(); spawn != nil {
		err := k.executeSpawn(*spawn, ctx, from, sourcePort, sourceChannel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) executeSpawn(spawn types.SpawnInstruction, ctx sdk.Context, from string, sourcePort string, sourceChannel string) error {
	// TODO: add logic for escrow and transfer tokens

	// nativeTransferToken := sdk.NewCoin(fungibleTokenPacketData.Denom, transferAmount)
	// ibcTransferToken := sdk.NewCoin(parachainInfo.IbcDenom, transferAmount)

	// escrowAddress := transfertypes.GetEscrowAddress(sourcePort, sourceChannel)
	// err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, escrowAddress, transfertypes.ModuleName, sdk.NewCoins(nativeTransferToken))
	// if err != nil {
	// 	return err
	// }

	// burn native token
	// Get Coin from escrow address
	// err = k.bankKeeper.BurnCoins(ctx, transfertypes.ModuleName, sdk.NewCoins(nativeTransferToken))
	// if err != nil {
	// 	return err
	// }

	// err = k.bankKeeper.SendCoins(ctx, escrowAddress, sender, sdk.NewCoins(ibcTransferToken))
	// if err != nil {
	// 	return err
	// }

	memo, err := json.Marshal(spawn.Program)
	if err != nil {
		return err
	}

	// TODO: add default timeout Height/Timestamp?
	transferMsg := transfertypes.MsgTransfer{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		// Token:         ibcTransferToken,
		Sender: from,
		// Receiver: fungibleTokenPacketData.Receiver,
		TimeoutHeight:    types1.Height{},
		TimeoutTimestamp: 0,
		Memo:             string(memo),
	}
	// TODO: handle result?
	_, err = k.executeTransferMsg(ctx, &transferMsg)
	if err != nil {
		return err
	}

	return nil
}
