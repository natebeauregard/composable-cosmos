package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDestinationAddress = errorsmod.Register(ModuleName, 1, "invalid destination address")
	ErrClientNotFound            = errorsmod.Register(ModuleName, 2, "client not found")
	ErrClientNotActive           = errorsmod.Register(ModuleName, 3, "client not active")
	ErrInvalidIntentId           = errorsmod.Register(ModuleName, 4, "invalid intent id")
	ErrInvalidTransferIntent     = errorsmod.Register(ModuleName, 5, "invalid transfer intent")
	ErrInvalidTxReceipt          = errorsmod.Register(ModuleName, 6, "invalid transaction receipt")
	ErrInvalidBlockHeaders       = errorsmod.Register(ModuleName, 7, "invalid block headers")

	ErrInvalidReceiptKey = errorsmod.Register(ModuleName, 8, "invalid receipt key")
	ErrReceiptNotFound   = errorsmod.Register(ModuleName, 9, "receipt not found")
	ErrBlockHashMismatch = errorsmod.Register(ModuleName, 10, "block hash mismatch")
	ErrBlockBodyMismatch = errorsmod.Register(ModuleName, 11, "block body mismatch")

	ErrTransferEventNotFound      = errorsmod.Register(ModuleName, 12, "transfer event not found")
	ErrTokenAddressMismatch       = errorsmod.Register(ModuleName, 13, "token address mismatch")
	ErrDestinationAddressMismatch = errorsmod.Register(ModuleName, 14, "destination address mismatch")
	ErrSourceAddressMismatch      = errorsmod.Register(ModuleName, 15, "source address mismatch")
	ErrAmountMismatch             = errorsmod.Register(ModuleName, 16, "amount mismatch")

	ErrInvalidReceiptSignature = errorsmod.Register(ModuleName, 17, "invalid receipt signature")
	ErrReceiptAlreadyProcessed = errorsmod.Register(ModuleName, 18, "receipt already processed")

	ErrInvalidSenderAddress          = errorsmod.Register(ModuleName, 19, "invalid sender address")
	ErrPrematureTimeoutTrigger       = errorsmod.Register(ModuleName, 20, "timeout must be triggered after current block time")
	ErrProofSubmittedAfterTimeout    = errorsmod.Register(ModuleName, 21, "intent execution proof cannot be submitted after timeout")
	ErrProofSubmittedBeforeExecution = errorsmod.Register(ModuleName, 22, "intent execution occurred before intent submission")
)
