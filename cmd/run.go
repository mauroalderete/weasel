/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mauroalderete/weasel/client"
	"github.com/mauroalderete/weasel/genwallet"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executes weasel to explore accounts",
	Long: `Start weasel service to generate random accounts and explore his activity.
Store accounts with activity.
For example:

weasel run --thread 12

Weasel is a tool to search accounts with activity generating a private key randomly.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().Int32P("thread", "t", 1, "Number of threads >0 to execute. Each thread handle his own connection and own search. By default is 1.")
	runCmd.Flags().StringP("gateway", "g", "https://cloudflare-eth.com", "Ethereum gateway to connect. By default is https://cloudflare-eth.com")

	runCmd.RunE = runMain
}

func runMain(cmd *cobra.Command, args []string) error {

	gateway := cmd.Flag("gateway").Value.String()

	err := Search(gateway)
	if err != nil {
		return fmt.Errorf("failed execute search: %v", err)
	}

	return nil
}

func Search(server string) error {
	c := &client.Client{}
	err := c.Connect(server)
	if err != nil {
		return fmt.Errorf("failed to connect to ethclient: %v", err)
	}
	defer c.Close()

	w, err := genwallet.RandomWallet()
	if err != nil {
		return fmt.Errorf("failed instance a new random wallet: %v", err)
	}

	err = w.Bind(c.Client())
	if err != nil {
		return fmt.Errorf("failed to bind client to wallet: %v", err)
	}

	err = w.Update()
	if err != nil {
		return fmt.Errorf("failed to update wallet info: %v", err)
	}

	fmt.Printf("Address: %s\n", w.Address().Hex())
	fmt.Printf("Balance: %sWEI\n", w.Balance().Wei().String())
	fmt.Printf("Balance: %sETH\n", w.Balance().Eth().String())
	fmt.Printf("Pending balance: %sWEI\n", w.PendingBalance().Wei().String())
	fmt.Printf("Pending balance: %sETH\n", w.PendingBalance().Eth().String())

	return nil
}
