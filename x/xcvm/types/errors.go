package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDestinationAddress = errorsmod.Register(ModuleName, 1, "invalid destination address")
	ErrClientNotFound            = errorsmod.Register(ModuleName, 2, "client not found")
	ErrClientNotActive           = errorsmod.Register(ModuleName, 3, "client not active")
)
