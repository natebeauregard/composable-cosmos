package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type LogTransfer struct {
	From         common.Address
	To           common.Address
	Tokens       *big.Int
	TokenAddress common.Address
}
