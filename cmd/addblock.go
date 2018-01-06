package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	blockData string

	addBlockCmd = &cobra.Command{
		Use:     "addblock",
		Short:   "Add a block to the blockchain",
		Long:    "Adds a block to the blockchain with the specified data",
		Run:     addBlock,
		PreRun:  cmdPreRun,
		PostRun: cmdPostRun,
	}
)

func init() {
	addBlockCmd.Flags().StringVarP(&blockData, "data", "d", "", "block data")
	rootCmd.AddCommand(addBlockCmd)
}

// Add a block to the crypto.
func addBlock(cmd *cobra.Command, _ []string) {
	// Validate the block data.
	if blockData == "" {
		fmt.Printf("invalid or missing block data\n")
		fmt.Println()

		cmd.Usage()
		return
	}

	bc.AddBlock(blockData)
	fmt.Println("Successfully added block to the blockchain!")
}
