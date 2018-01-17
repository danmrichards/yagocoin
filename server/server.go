package server

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/danmrichards/yagocoin/crypto"
)

const (
	protocol      = "tcp"
	nodeVersion   = 1
	commandLength = 12
)

var (
	nodeAddress     string
	miningAddress   string
	KnownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	mempool         = make(map[string]crypto.Transaction)
)

type addr struct {
	AddrList []string
}

// block represents a message to transfer a block.
type block struct {
	AddrFrom string
	Block    []byte
}

// getBlocks represents a request to get a list of all block hashes.
type getBlocks struct {
	AddrFrom string
}

// getData represents a request to get a specific block or transaction.
type getData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

// inv represents an inventory of block hashes.
type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

// tx represents a message to transfer a transaction.
type tx struct {
	AddFrom     string
	Transaction []byte
}

// version represents the current state of a blockchain that a node has.
type version struct {
	Version    int
	BestHeight int    // Length of the nodes blockchain.
	AddrFrom   string // The address of the sender
}

// commandToBytes returns a byte array representing a command.
// Our server messages are byte arrays, the first 12 of which specify the name
// of the command the message represents.
func commandToBytes(command string) []byte {
	var outBytes [commandLength]byte

	for i, c := range command {
		outBytes[i] = byte(c)
	}

	return outBytes[:]
}

// bytesToCommand extracts and returns the command name from a message.
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

// requestBlocks sends the list of block hashes for each node.
func requestBlocks() {
	for _, node := range KnownNodes {
		sendGetBlocks(node)
	}
}

// sendBlock sends a message representing a block.
func sendBlock(addr string, b *crypto.Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

// sendData sends a set of data to the specified address.
func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		KnownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

// sendInv sends an inventory to the given address.
func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

// sendGetBlocks sens a 'get blocks' message to an address.
func sendGetBlocks(address string) {
	payload := gobEncode(getBlocks{nodeAddress})
	request := append(commandToBytes("getBlocks"), payload...)

	sendData(address, request)
}

// sendGetData sends a 'get data' request to an address.
func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getData{nodeAddress, kind, id})
	request := append(commandToBytes("getData"), payload...)

	sendData(address, request)
}

// SendTx sends a message representing a transaction.
func SendTx(addr string, tnx *crypto.Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

// sendVersion responds with the state of this nodes blockchain.
func sendVersion(addr string, bc *crypto.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

// handleAddr handles a request to get an address list.
func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(KnownNodes))
	requestBlocks()
}

// handleBlock handles a request to get a specific block.
func handleBlock(request []byte, bc *crypto.Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := crypto.DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := crypto.UTxOSet{bc}
		UTXOSet.Reindex()
	}
}

// handleInv handles a request informing the node of an inventory.
func handleInv(request []byte) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		var newInTransit [][]byte
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

// handleGetBlocks handles a request to list all block hashes.
func handleGetBlocks(request []byte, bc *crypto.Blockchain) {
	var buff bytes.Buffer
	var payload getBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

// handleGetData handles a request to get a specific block or transaction.
func handleGetData(request []byte, bc *crypto.Blockchain) {
	var buff bytes.Buffer
	var payload getData

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		SendTx(payload.AddrFrom, &tx)
	}
}

// handleTx handles a request to get a specific transaction.
func handleTx(request []byte, bc *crypto.Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := crypto.DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*crypto.Transaction

			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := crypto.NewCoinbaseTx(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := crypto.UTxOSet{bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range KnownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

// handleVersion handles a request for blockchain version information.
func handleVersion(request []byte, bc *crypto.Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

// handleConnection handles a new server connection and dispatches the relevant
// function based on the command received.
func handleConnection(conn net.Conn, bc *crypto.Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request)
	case "getBlocks":
		handleGetBlocks(request, bc)
	case "getData":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

// StartServer starts a node server.
func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := crypto.NewBlockchain(nodeID)

	// If this is not the central node, send a request to it to check if the
	// blockchain is up to date.
	if nodeAddress != KnownNodes[0] {
		sendVersion(KnownNodes[0], bc)
	}

	// Handle commands as they come in.
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

// goEncode encodes data as gob.
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// nodeIsKnown checks if the given address is known to us.
func nodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}

	return false
}
