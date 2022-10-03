package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mauroalderete/weasel/coin"
)

type wallet struct {
	client         *ethclient.Client
	address        common.Address
	balance        coin.Coin
	pendingBalance coin.Coin
}

func (w *wallet) Address() common.Address {
	return w.address
}

func (w *wallet) Balance() *coin.Coin {
	c := w.balance
	return &c
}

func (w *wallet) PendingBalance() *coin.Coin {
	c := w.pendingBalance
	return &c
}

func (w *wallet) Update() error {

	err := w.QueryBalance(nil)
	if err != nil {
		return fmt.Errorf("failed to update last balance: %v", err)
	}

	err = w.QueryPendingBalance()
	if err != nil {
		return fmt.Errorf("failed to update pending balance: %v", err)
	}

	return nil
}

func (w *wallet) QueryBalance(blockNumber *big.Int) error {
	balance, err := w.client.BalanceAt(context.Background(), w.address, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to query balance from account %s: %v", w.address.Hex(), err)
	}

	err = w.balance.SetWei(*balance)
	if err != nil {
		return fmt.Errorf("failed to store %sWEI as balance requested from account %s: %v", balance.String(), w.address.Hex(), err)
	}

	return nil
}

func (w *wallet) QueryPendingBalance() error {
	balance, err := w.client.PendingBalanceAt(context.Background(), w.address)
	if err != nil {
		return fmt.Errorf("failed to query pending balance from account %s: %v", w.address.Hex(), err)
	}

	err = w.pendingBalance.SetWei(*balance)
	if err != nil {
		return fmt.Errorf("failed to store %sWEI as pending balance requested from account %s: %v", balance.String(), w.address.Hex(), err)
	}

	return nil
}

func New(address string, client *ethclient.Client) (*wallet, error) {

	if client == nil {
		return nil, fmt.Errorf("client is required")
	}

	w := &wallet{}
	w.client = client
	w.address = common.HexToAddress(address)

	return w, nil
}
