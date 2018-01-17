package cmd

import (
	"fmt"
	"log"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var (
	listAddressesCmd = &cobra.Command{
		Use:     "listaddresses",
		Short:   "Lists all addresses from the wallet file",
		Run:     listAddresses,
		Args:    cobra.ExactArgs(0),
		PreRun:  cmdPreRun,
		PostRun: cmdPostRun,
	}
)

func init() {
	rootCmd.AddCommand(listAddressesCmd)
}

// Lists all addresses from the wallet file.
func listAddresses(_ *cobra.Command, _ []string) {
	wallets, err := crypto.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
