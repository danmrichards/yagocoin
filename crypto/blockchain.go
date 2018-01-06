package crypto

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	blocksBucket = "blocks"
	dbFile       = "blockchain.db"
	fileMode     = 0600
	hashKey      = "l"
)

// Blockchain represents the chain of blocks.
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator is used to iterate over the crypto.
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// AddBlock creates a new block and adds it to the crypto.
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// Get the hash of the last block in the DB.
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(hashKey))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// Mine a new block and add to the DB.
	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte(hashKey), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// Iterator returns a new iterator for the current crypto.
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// Close closes the blockchain database connection.
func (bc *Blockchain) Close() error {
	return bc.db.Close()
}

// Next returns next block starting from the tip.
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	// Get the current block.
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

// NewBlockchain creates a new crypto with a genesis block.
func NewBlockchain() *Blockchain {
	var tip []byte

	// Open or create crypto db.
	db, err := bolt.Open(dbFile, fileMode, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// Get the block bucket.
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			fmt.Println("no existing blockchain found, creating a new one...")
			fmt.Println()

			// Create new genesis block and bucket if the bucket did not exist.
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte(hashKey), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = genesis.Hash
		} else {
			// Bucket exists, get the tip.
			tip = b.Get([]byte(hashKey))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}
