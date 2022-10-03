//go:build !test

package main

import (
	"fmt"
	"log"

	_ "embed"

	"github.com/bgadrian/go-mnemonic/bip39"
	"github.com/mauroalderete/weasel/client"
	"github.com/mauroalderete/weasel/genwallet"
)

//go:generate build/get_version.sh
//go:embed version.txt
var version string

func main() {

	const server = "https://cloudflare-eth.com"

	Look(server)

	m, err := bip39.NewMnemonicRandom(128, "")
	if err != nil {
		log.Fatal("Whoops something went wrong to make mnemonic generator", err)
	}

	password, err := m.GetSentence()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(password)

	fmt.Printf("Bye\n")
}

func Look(server string) {
	c := &client.Client{}
	err := c.Connect(server)
	if err != nil {
		log.Fatalf("failed to connect to ethclient: %v", err)
	}
	defer c.Close()

	w, err := genwallet.RandomWallet()
	if err != nil {
		log.Fatal("failed instance a new random wallet", err)
	}

	err = w.Bind(c.Client())
	if err != nil {
		log.Fatalf("failed to bind client to wallet: %v", err)
	}

	err = w.Update()
	if err != nil {
		log.Fatalf("failed to update wallet info: %v", err)
	}

	fmt.Printf("Address: %s\n", w.Address().Hex())
	fmt.Printf("Balance: %sWEI\n", w.Balance().Wei().String())
	fmt.Printf("Balance: %sETH\n", w.Balance().Eth().String())
	fmt.Printf("Pending balance: %sWEI\n", w.PendingBalance().Wei().String())
	fmt.Printf("Pending balance: %sETH\n", w.PendingBalance().Eth().String())
}
