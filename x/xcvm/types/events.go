package types

// xcvm module event types
const (
	EventAddTransferIntent         = "add-transfer-intent"
	EventVerifyTransferIntentProof = "verify-transfer-intent-proof"

	AttributeKeyIntentId           = "intent-id"
	AttributeKeyClientId           = "client-id"
	AttributeKeySourceAddress      = "source-address"
	AttributeKeyDestinationAddress = "destination-address"
	AttributeKeyAmount             = "amount"
	AttributeKeyBounty             = "bounty"
)
