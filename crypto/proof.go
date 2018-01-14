package crypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/danmrichards/yagocoin/domain"
)

// Maximum counter for proof-of-work algorithm.
var maxNonce = math.MaxInt64

// Sets the upper boundary of the hash target.
const targetBits = 16

// Proof represents a proof-of-work.
type Proof struct {
	block  *Block
	target *big.Int
}

// Prepare the data needed to hash a block.
func (p *Proof) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.HashTransactions(),
			domain.IntToHex(p.block.Timestamp.Unix()),
			domain.IntToHex(int64(targetBits)),
			domain.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof-of-work run for a block. We iteratively hash the block
// and compare the hash to the proof target. The comparison is done by writing
// the block hash bytes to a big int. If the block hash is less than the target
// then the proof is valid and vice versa.
func (p *Proof) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Println("Mining a new block")
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

// Validate validates a blocks proof-of-work.
func (p *Proof) Validate() bool {
	var hashInt big.Int

	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(p.target) == -1
}

// NewProof builds a new proof with a target.  The target for the proof is 256
// minus our targetBits value, because 256 is the length of the hash we get
// from a block.
func NewProof(b *Block) *Proof {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &Proof{b, target}
}
