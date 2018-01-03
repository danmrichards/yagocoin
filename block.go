package main

import "time"

// Represents the core structure of a block.
type Block struct {
	Timestamp     time.Time
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Create a new block.
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now(), []byte(data), prevBlockHash, []byte{}, 0}

	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Creates a new "genesis" block to start a chain.
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
