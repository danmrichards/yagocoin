package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	getBalanceCmd = &cobra.Command{
		Use:     "getbalance",
		Short:   "Get balance of adress",
		Run:     getBalance,
		Args:    cobra.ExactArgs(0),
		PreRun:  cmdPreRun,
		PostRun: cmdPostRun,
	}
)

func init() {
	getBalanceCmd.Flags().StringVarP(&address, "address", "a", "", "Address to send the balance of")
	rootCmd.AddCommand(getBalanceCmd)
}

// Get balance of an address.
func getBalance(cmd *cobra.Command, _ []string) {
	// Validate the address.
	if address == "" {
		fmt.Printf("Invalid or missing address\n")
		fmt.Println()

		cmd.Usage()
		return
	}

	balance := 0
	UTXOs := bc.FindUTxO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
