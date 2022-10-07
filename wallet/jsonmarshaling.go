package wallet

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mauroalderete/weasel/coin"
)

type walletJsonMarshable struct {
	Datetime       string             `json:"datetime"`
	Client         string             `json:"client"`
	Address        string             `json:"address"`
	PrivateKey     string             `json:"privateKey"`
	PublicKey      string             `json:"publicKey"`
	Balance        coin.CoinMarshable `json:"balance"`
	PendingBalance coin.CoinMarshable `json:"pendingBalance"`
}

func JsonMarshal(w Wallet) ([]byte, error) {

	wm := &walletJsonMarshable{
		Datetime:       time.Now().String(),
		Client:         w.gateway.Url(),
		Address:        w.Address().Hex(),
		PrivateKey:     w.PrivateKey().Hex(),
		PublicKey:      w.PublicKey().Hex(),
		Balance:        w.Balance().Marshable(),
		PendingBalance: w.PendingBalance().Marshable(),
	}

	payload, err := json.Marshal(wm)
	if err != nil {
		return nil, fmt.Errorf("failed to marshaling json")
	}

	return payload, nil
}
