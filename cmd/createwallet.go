package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var (
	createWalletCmd = &cobra.Command{
		Use:   "createwallet",
		Short: "Generates a new key-pair and saves it into the wallet file",
		Run:   createWallet,
		Args:  cobra.ExactArgs(0),
	}
)

func init() {
	rootCmd.AddCommand(createWalletCmd)
}

// Create a new key-pair and save it into the wallet file
func createWallet(cmd *cobra.Command, _ []string) {
	wallets, err := crypto.NewWallets()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("could not create wallet: %s", err)

		fmt.Println("Could not create wallet at this time!")
		return
	}

	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}
