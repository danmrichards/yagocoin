package crypto

import (
	"bytes"
	"encoding/gob"
	"log"
)

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

// TXOutputs collects TXOutput
type TxOutputs struct {
	Outputs []TxOutput
}

// Serialize serializes TXOutputs
func (outs TxOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
