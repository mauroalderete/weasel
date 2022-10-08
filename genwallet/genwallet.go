package genwallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mauroalderete/weasel/wallet"
)

func RandomWallet() (*wallet.Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed get privateKey: %v", err)
	}

	w, err := wallet.NewFromPrivateKey(*privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new wallet: %v", err)
	}

	return w, nil
}
