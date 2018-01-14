package cmd

import (
	"fmt"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/spf13/cobra"
)

var reindexUTxOCmd = &cobra.Command{
	Use:     "reindexutxo",
	Short:   "Rebuilds the UTXO set",
	Run:     reindexUTxO,
	Args:    cobra.ExactArgs(0),
	PreRun:  cmdPreRun,
	PostRun: cmdPostRun,
}

func init() {
	rootCmd.AddCommand(reindexUTxOCmd)
}

// Rebuilds the UTXO set.
func reindexUTxO(_ *cobra.Command, _ []string) {
	uTxOSet := crypto.UTxOSet{bc}
	uTxOSet.Reindex()

	count := uTxOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
