package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

const subsidy = 10 // The amount of reward for mining.

// Transaction represents a yagocoin transaction.
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// IsCoinbase checks whether the transaction is a 'coinbase' transaction.
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// Serialize returns a serialized Transaction.
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// Hash returns the hash of the Transaction.
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// NewCoinbaseTx creates a new 'coinbase' transaction. This is a special type
// of transactions, which doesn’t require previously existing outputs. It
// creates outputs (i.e. coins) out of nowhere becoming the reward miners get
// for mining new blocks.
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TxInput{[]byte{}, -1, nil, []byte(data)}
	txOut := NewTxOutput(subsidy, to)

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}
	tx.ID = tx.Hash()

	return &tx
}

// NewUTxOTransaction creates a new transaction.
func NewUTxOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		fmt.Println("ERROR: Not enough funds")
		os.Exit(1)
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTxOutput(amount, to))

	// Change.
	if acc > amount {
		outputs = append(outputs, *NewTxOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()

	return &tx
}
