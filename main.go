package main

import (
	"fmt"
	"strconv"
)

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 YGC to Dan")
	bc.AddBlock("Send 2 more YCG to Dan")

	for _, block := range bc.blocks {
		fmt.Printf("prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)

		pow := NewProof(block)
		fmt.Printf("proof: %s\n", strconv.FormatBool(pow.Validate()))

		fmt.Println()
	}
}
