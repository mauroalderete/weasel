//go:build !test

package main

import (
	"fmt"
	"log"

	_ "embed"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mauroalderete/weasel/wallet"
)

//go:generate build/get_version.sh
//go:embed version.txt
var version string

func main() {

	const server = "https://cloudflare-eth.com"

	fmt.Printf("Connecting... ")
	client, err := ethclient.Dial(server)
	if err != nil {
		fmt.Printf("[FAIL]\n")
		log.Fatalf("failed to connect to ethclient: %v", err)
	}

	_ = client
	fmt.Printf("[OK]\n")

	fmt.Printf("Opening wallet... ")
	w, err := wallet.New("0x71c7656ec7ab88b098defb751b7401b5f6d8976f", client)
	if err != nil {
		fmt.Printf("[FAIL]\n")
		log.Fatalf("failed to instance a new wallet: %v", err)
	}

	err = w.Update()
	if err != nil {
		fmt.Printf("[FAIL]\n")
		log.Fatalf("failed to open the wallet with address %s: %v", w.Address().Hex(), err)
	}

	fmt.Printf("[OK]\n\n")

	fmt.Printf("Address: %s\n", w.Address().Hex())
	fmt.Printf("Balance: %sWEI\n", w.Balance().Wei().String())
	fmt.Printf("Balance: %sETH\n", w.Balance().Eth().String())
	fmt.Printf("Pending balance: %sWEI\n", w.PendingBalance().Wei().String())
	fmt.Printf("Pending balance: %sETH\n", w.PendingBalance().Eth().String())

	fmt.Printf("Bye\n")
}
