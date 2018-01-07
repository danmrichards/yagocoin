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

// SetID sets the ID of a transaction.
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// TxInput respresents a transaction input.
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// CanUnlockOutputWith checks whether the address initiated the transaction.
func (in *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// TxOutput represents a transaction output.
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// CanBeUnlockedWith checks if the output can be unlocked with the data.
func (out *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTx creates a new 'coinbase' transaction. This is a special type
// of transactions, which doesn’t require previously existing outputs. It
// creates outputs (i.e. coins) out of nowhere becoming the reward miners get
// for mining new blocks.
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{subsidy, to}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetID()

	return &tx
}

// NewUTxOTransaction creates a new transaction.
func NewUTxOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

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
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TxOutput{amount, to})

	// Change.
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
