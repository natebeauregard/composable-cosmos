package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

const (
	TypeMsgSendTransferIntent        = "send_transfer_intent"
	TypeMsgVerifyTransferIntentProof = "verify_transfer_intent_proof"
)

var _ sdk.Msg = &MsgSendTransferIntent{}

func NewMsgSendTransferIntent(
	fromAddress string,
	destinationAddress string,
	clientId string,
	amount math.Uint,
) *MsgSendTransferIntent {
	return &MsgSendTransferIntent{
		FromAddress:        fromAddress,
		DestinationAddress: destinationAddress,
		ClientId:           clientId,
		Amount:             amount,
	}
}

// Type Implements Msg.
func (MsgSendTransferIntent) Type() string { return TypeMsgSendTransferIntent }

// GetSigners returns the expected signers for a MsgSendTransferIntent message.
func (msg *MsgSendTransferIntent) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgSendTransferIntent) ValidateBasic() error {
	// validate from address
	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}

	// validate destination address
	if !common.IsHexAddress(msg.DestinationAddress) {
		return ErrInvalidDestinationAddress
	}

	// validate clientId
	if err := host.ClientIdentifierValidator(msg.ClientId); err != nil {
		return err
	}

	return nil
}

var _ sdk.Msg = &MsgVerifyTransferIntentProof{}

func NewMsgVerifyTransferIntentProof(
	signer string,
	proof []byte,
	intentId uint64,
) *MsgVerifyTransferIntentProof {
	return &MsgVerifyTransferIntentProof{
		Signer:   signer,
		Proof:    proof,
		IntentId: intentId,
	}
}

// Type Implements Msg.
func (MsgVerifyTransferIntentProof) Type() string { return TypeMsgVerifyTransferIntentProof }

// GetSigners returns the expected signers for a MsgVerifyTransferIntentProof message.
func (msg *MsgVerifyTransferIntentProof) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgVerifyTransferIntentProof) ValidateBasic() error {
	// validate signer
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return errorsmod.Wrap(err, "invalid signer address")
	}

	return nil
}
