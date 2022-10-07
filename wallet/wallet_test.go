package wallet

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/mauroalderete/weasel/client"
// )

// const (
// 	pk     = "....."
// 	server = "https://cloudflare-eth.com"
// )

// func TestQueryBalanceNonZero(t *testing.T) {

// 	privateKey, err := crypto.HexToECDSA(pk)
// 	if err != nil {
// 		t.Errorf("expect error nil, got: %v", err)
// 		return
// 	}

// 	w, err := NewFromPrivateKey(*privateKey)
// 	if err != nil {
// 		t.Errorf("expect error nil, got: %v", err)
// 		return
// 	}

// 	c := &client.Client{}
// 	err = c.Connect(server)
// 	if err != nil {
// 		t.Errorf("failed to connect to ethclient: %v", err)
// 		return
// 	}

// 	err = w.Bind(c)
// 	if err != nil {
// 		t.Errorf("failed to bind client to wallet: %v", err)
// 		return
// 	}

// 	err = w.Update()
// 	if err != nil {
// 		t.Errorf("expect error nil, got: %v", err)
// 		return
// 	}

// 	const tolerance float64 = 0.0000000000000000001
// 	var balance float64 = 0.007276439452
// 	var pendingbalance float64 = 0

// 	fmt.Printf("public: %s\n", w.Address().Hex())
// 	fmt.Printf("be: %s\n", w.balance.Eth())
// 	fmt.Printf("bw: %s\n", w.balance.Wei())
// 	fmt.Printf("be: %s\n", w.pendingBalance.Eth())
// 	fmt.Printf("bw: %s\n", w.pendingBalance.Wei())

// 	if val, _ := w.Balance().Eth().Float64(); val-balance <= tolerance {
// 		t.Errorf("expect balance %f, got %f", balance, val)
// 	}

// 	if val, _ := w.PendingBalance().Eth().Float64(); val-pendingbalance <= tolerance {
// 		t.Errorf("expect pending balance %f, got %f", balance, val)
// 	}
// }
