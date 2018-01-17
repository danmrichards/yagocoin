package cmd

import (
	"fmt"
	"os"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var (
	bc *crypto.Blockchain

	address string

	nodeID string

	rootCmd = &cobra.Command{
		Use:   "yagocoin",
		Short: "Yet Another Go Coin",
		Long: `A proof-of-concept cryptocurrency written in Go.

Ships with a basic CLI tool for adding to and viewing the block chain.

Based on a simple blockchain as described at https://jeiwan.cc`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func cmdPreRun(_ *cobra.Command, _ []string) {
	nodeID = os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	// Open the connection to the blockchain db.
	bc = crypto.NewBlockchain(nodeID)
}

func cmdPostRun(_ *cobra.Command, _ []string) {
	// Close the connection to the blockchain db.
	bc.Close()
}
