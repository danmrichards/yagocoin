package cmd

import (
	"fmt"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var (
	createBlockchainCmd = &cobra.Command{
		Use:   "createblockchain",
		Short: "Create a new blockchain",
		Run:   createBlockchain,
		Args:  cobra.ExactArgs(0),
	}
)

func init() {
	createBlockchainCmd.Flags().StringVarP(&address, "address", "a", "", "Address to send the genesis block reward to")
	rootCmd.AddCommand(createBlockchainCmd)
}

// Create a new blockchain.
func createBlockchain(cmd *cobra.Command, _ []string) {
	// Validate the address.
	if address == "" {
		fmt.Printf("Invalid or missing address\n")
		fmt.Println()

		cmd.Usage()
		return
	}

	bc := crypto.CreateBlockchain(address)
	bc.Close()

	fmt.Println("Done!")
}
