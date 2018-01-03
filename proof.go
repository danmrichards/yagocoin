package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// Maximum counter for proof-of-work algorithm.
var maxNonce = math.MaxInt64

// Sets the upper boundary of the hash target.
const targetBits = 24

// Represents a proof-of-work.
type Proof struct {
	block  *Block
	target *big.Int
}

// Prepare the data needed to hash a block.
func (p *Proof) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.Data,
			IntToHex(p.block.Timestamp.Unix()),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Performs a proof-of-work.
func (p *Proof) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("mining the block containing \"%s\"\n", p.block.Data)
	for nonce < maxNonce {
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validates a blocks proof-of-work.
func (p *Proof) Validate() bool {
	var hashInt big.Int

	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(p.target) == -1
}

// Builds a new proof.
func NewProof(b *Block) *Proof {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &Proof{b, target}
}
