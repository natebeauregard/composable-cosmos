package types

var (
	// DefaultIntentId is the byte representation of uint64(0)
	DefaultIntentId = uint64(0)
)

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		IntentId: DefaultIntentId,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func ValidateGenesis(data GenesisState) error {
	return nil
}
