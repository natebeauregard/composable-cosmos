package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDestinationAddress = errorsmod.Register(ModuleName, 1, "invalid destination address")
	ErrClientNotFound            = errorsmod.Register(ModuleName, 2, "client not found")
	ErrClientNotActive           = errorsmod.Register(ModuleName, 3, "client not active")
	ErrInvalidIntentId           = errorsmod.Register(ModuleName, 4, "invalid intent id")
	ErrInvalidTxReceipt          = errorsmod.Register(ModuleName, 5, "invalid transaction receipt")
	ErrInvalidBlockHeaders       = errorsmod.Register(ModuleName, 6, "invalid block headers")

	ErrInvalidReceiptKey = errorsmod.Register(ModuleName, 7, "invalid receipt key")
	ErrReceiptNotFound   = errorsmod.Register(ModuleName, 8, "receipt not found")
	ErrBlockHashMismatch = errorsmod.Register(ModuleName, 9, "block hash mismatch")
	ErrBlockBodyMismatch = errorsmod.Register(ModuleName, 10, "block body mismatch")

	ErrTransferEventNotFound      = errorsmod.Register(ModuleName, 11, "transfer event not found")
	ErrTokenAddressMismatch       = errorsmod.Register(ModuleName, 12, "token address mismatch")
	ErrDestinationAddressMismatch = errorsmod.Register(ModuleName, 13, "destination address mismatch")
	ErrSourceAddressMismatch      = errorsmod.Register(ModuleName, 14, "source address mismatch")
	ErrAmountMismatch             = errorsmod.Register(ModuleName, 15, "amount mismatch")

	ErrInvalidReceiptSignature = errorsmod.Register(ModuleName, 16, "invalid receipt signature")
	ErrReceiptAlreadyProcessed = errorsmod.Register(ModuleName, 17, "receipt already processed")

	ErrInvalidSenderAddress       = errorsmod.Register(ModuleName, 18, "invalid sender address")
	ErrPrematureTimeoutTrigger    = errorsmod.Register(ModuleName, 19, "timeout must be triggered after current block time")
	ErrProofSubmittedAfterTimeout = errorsmod.Register(ModuleName, 20, "intent execution proof cannot be submitted after timeout")
)
