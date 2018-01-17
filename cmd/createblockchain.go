package cmd

import (
	"fmt"
	"log"
	"os"

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
	nodeID = os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	// Validate the address.
	if address == "" {
		fmt.Printf("Invalid or missing address\n")
		fmt.Println()

		cmd.Usage()
		return
	}

	if !crypto.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := crypto.CreateBlockchain(address, nodeID)
	defer bc.Close()

	UTXOSet := crypto.UTxOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}
