package cmd

import (
	"fmt"
	"log"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/danmrichards/yagocoin/server"
	"github.com/spf13/cobra"
)

var (
	from, to string
	amount   int
	mineNow  bool

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
	sendCmd.Flags().BoolVarP(&mineNow, "mine", "m", false, "Mine immediately on the same node")
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

	uTxOSet := crypto.UTxOSet{bc}

	wallets, err := crypto.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := crypto.NewUTxOTransaction(&wallet, to, amount, &uTxOSet)

	if mineNow {
		cbTx := crypto.NewCoinbaseTx(from, "")
		txs := []*crypto.Transaction{cbTx, tx}

		newBlock := bc.MineBlock(txs)
		uTxOSet.Update(newBlock)
	} else {
		server.SendTx(server.KnownNodes[0], tx)
	}

	fmt.Println("Success!")
}
