package cmd

import (
	"fmt"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var (
	from, to string
	amount   int

	sendCmd = &cobra.Command{
		Use:     "send",
		Short:   "Send an amount of coins from one address to another",
		Run:     send,
		Args:    cobra.ExactArgs(0),
		PreRun:  cmdPreRun,
		PostRun: cmdPostRun,
	}
)

func init() {
	sendCmd.Flags().StringVarP(&from, "from", "f", "", "Address to send the coins from")
	sendCmd.Flags().StringVarP(&to, "to", "t", "", "Address to send the coins to")
	sendCmd.Flags().IntVarP(&amount, "amount", "a", 0, "Amount of coins to send")
	rootCmd.AddCommand(sendCmd)
}

// Send an amount of coins from one address to another.
func send(cmd *cobra.Command, _ []string) {
	// Validate the from adress.
	if from == "" {
		fmt.Printf("Invalid or missing from address\n")
		fmt.Println()

		cmd.Usage()
		return
	}

	// Validate the to adress.
	if to == "" {
		fmt.Printf("Invalid or missing to address\n")
		fmt.Println()

		cmd.Usage()
		return
	}

	// Validate the to amount.
	if amount == 0 {
		fmt.Printf("Invalid or missing amount\n")
		fmt.Println()

		cmd.Usage()
		return
	}

	tx := crypto.NewUTxOTransaction(from, to, amount, bc)
	bc.MineBlock([]*crypto.Transaction{tx})

	fmt.Println("Success!")
}
