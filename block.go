package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Represents the core structure of a block.
type Block struct {
	Timestamp     time.Time
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// Generate and set a new SHA256 hash for the block.
func (b *Block) SetHash() {
	headers := bytes.Join([][]byte{
		b.PrevBlockHash,
		b.Data,
		[]byte(strconv.FormatInt(b.Timestamp.Unix(), 10)),
	}, []byte{})

	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// Create a new block.
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

// Creates a new "genesis" block to start a chain.
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
