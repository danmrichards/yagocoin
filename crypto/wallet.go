package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/danmrichards/yagocoin/base58"

	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	walletFile         = "wallet.dat"
	addressChecksumLen = 4
)

// Wallet stores private and public keys.
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// GetAddress returns the wallet address.
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Base58Encode(fullPayload)

	return address
}

// HashPubKey hashes and returns the public key.
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress check if address if valid.
func ValidateAddress(address string) bool {
	// Decode the hash.
	pubKeyHash := base58.Base58Decode([]byte(address))

	// Get the checksum and version.
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]

	// Extract the public key hash from the address.
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]

	// Create new checksum and compare.
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// NewWallet creates and returns a new Wallet.
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// GetPublicKeyHash returns the public key hash from a base58 encoded address.
func GetPublicKeyHash(address []byte) []byte {
	pubKeyHash := base58.Base58Decode(address)
	return pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
}

// Generates and returns a SHA256 checksum for the given payload.
// Hash will be of the length defined by addressChecksumLen.
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// Creates a new ecdsa key pair.
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	// In ecdsa public keys are on a curve hence the public key is a combination
	// of the x and y co-ordinates.
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
