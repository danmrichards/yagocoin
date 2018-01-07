package cmd

import (
	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var (
	bc *crypto.Blockchain

	address string

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
	// Open the connection to the blockchain db.
	bc = crypto.NewBlockchain()
}

func cmdPostRun(_ *cobra.Command, _ []string) {
	// Close the connection to the blockchain db.
	bc.Close()
}
