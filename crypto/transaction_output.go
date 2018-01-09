package crypto

import "bytes"

// TxOutput represents a transaction output.
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock signs the output.
func (out *TxOutput) Lock(address []byte) {
	out.PubKeyHash = GetPublicKeyHash(address)
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey.
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTxOutput create a new TXOutput.
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}
