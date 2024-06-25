package types

import (
	"fmt"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	prysmtypes "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
)

func (header *BeaconBlockHeader) Hash() ([32]byte, error) {
	prysmHeader := prysmtypes.BeaconBlockHeader{
		Slot:          primitives.Slot(header.Slot),
		ProposerIndex: primitives.ValidatorIndex(header.ProposerIndex),
		ParentRoot:    header.ParentRoot,
		StateRoot:     header.StateRoot,
		BodyRoot:      header.BodyRoot,
	}
	hash, err := prysmHeader.HashTreeRoot()
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to hash beacon block header: %v", err)
	}
	return hash, nil
}
