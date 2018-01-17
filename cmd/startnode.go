package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/danmrichards/yagocoin/crypto"
	"github.com/danmrichards/yagocoin/server"
	"github.com/spf13/cobra"
)

var (
	minerAddress string

	startNodeCmd = &cobra.Command{
		Use:   "startnode",
		Short: "Start a node with ID specified in NODE_ID env. var.",
		Run:   startNode,
		Args:  cobra.ExactArgs(0),
	}
)

func init() {
	startNodeCmd.Flags().StringVarP(&minerAddress, "miner", "m", "", "Enable mining mode and send reward to address")
	rootCmd.AddCommand(startNodeCmd)
}

// Start a node with ID specified in NODE_ID env. var.
func startNode(_ *cobra.Command, _ []string) {
	nodeID = os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if crypto.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}

	server.StartServer(nodeID, minerAddress)
}
