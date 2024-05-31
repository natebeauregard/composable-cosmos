package types

import "github.com/ethereum/go-ethereum/common"

func (rp ReceiptProof) Has(key []byte) (bool, error) {
	if len(key) != 32 {
		return false, ErrInvalidReceiptKey
	}
	_, ok := rp.Proof[common.Bytes2Hex(key[0:32])]
	return ok, nil
}

func (rp ReceiptProof) Get(key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrInvalidReceiptKey
	}
	value, ok := rp.Proof[common.Bytes2Hex(key[0:32])]
	if !ok {
		return nil, ErrReceiptNotFound
	}
	return value, nil
}

func (rp ReceiptProof) Put(key []byte, value []byte) error {
	if len(key) != 32 {
		return ErrInvalidReceiptKey
	}
	rp.Proof[common.Bytes2Hex(key[0:32])] = value
	return nil
}

func (rp ReceiptProof) Delete(key []byte) error {
	if len(key) != 32 {
		return ErrInvalidReceiptKey
	}
	delete(rp.Proof, common.Bytes2Hex(key[0:32]))
	return nil
}
