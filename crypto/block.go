package crypto

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block represents the core structure of a block.
type Block struct {
	Timestamp     time.Time
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Serialize serializes a block as gob for storage.
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// NewBlock creates a new block.
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now(), []byte(data), prevBlockHash, []byte{}, 0}

	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates a new "genesis" block to start a chain.
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// DeserializeBlock deserializes a block from gob.
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
