package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "xcvm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// Keys to use for the keeper store.
var (
	TransferIntentIdKey            = []byte{0x00}
	PendingTransferIntentKeyPrefix = []byte{0x01}
)

func GetPendingTransferIntentKeyById(intentId uint64) []byte {
	intentIdBz := make([]byte, 8)
	binary.BigEndian.PutUint64(intentIdBz, intentId)
	return append(PendingTransferIntentKeyPrefix, intentIdBz...)
}
