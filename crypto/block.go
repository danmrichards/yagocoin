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
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
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

// HashTransactions returns a hash of the transactions in the block.
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// NewBlock creates a new block.
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{time.Now(), transactions, prevBlockHash, []byte{}, 0, height}

	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates a new "genesis" block to start a chain.
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
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
