package cmd

import (
	"fmt"
	"strconv"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var printChainCmd = &cobra.Command{
	Use:     "printchain",
	Short:   "Print all the blocks of the blockchain",
	Run:     printChain,
	Args:    cobra.ExactArgs(0),
	PreRun:  cmdPreRun,
	PostRun: cmdPostRun,
}

func init() {
	rootCmd.AddCommand(printChainCmd)
}

// Print all the blocks of the blockchain.
func printChain(_ *cobra.Command, _ []string) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)

		pow := crypto.NewProof(block)
		fmt.Printf("Proof: %s\n", strconv.FormatBool(pow.Validate()))

		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
