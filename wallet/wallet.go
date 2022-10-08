package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mauroalderete/weasel/client"
	"github.com/mauroalderete/weasel/coin"
	"github.com/mauroalderete/weasel/wallet/walletkey"
)

type Wallet struct {
	gateway        *client.Client
	address        common.Address
	privateKey     walletkey.PrivateKey
	publicKey      walletkey.PublicKey
	balance        coin.Coin
	pendingBalance coin.Coin
}

func (w *Wallet) String() string {
	return fmt.Sprintf("{%s}[%s](%f;%f)", w.privateKey.Hex(), w.publicKey.Hex(), w.balance.Eth(), w.pendingBalance.Eth())
}

func (w *Wallet) Bind(client *client.Client) error {
	if client == nil {
		return fmt.Errorf("client is required")
	}

	w.gateway = client

	return nil
}

func (w *Wallet) PrivateKey() *walletkey.PrivateKey {
	pk := w.privateKey
	return &pk
}

func (w *Wallet) PublicKey() *walletkey.PublicKey {
	pk := w.publicKey
	return &pk
}

func (w *Wallet) Address() common.Address {
	return w.address
}

func (w *Wallet) Balance() *coin.Coin {
	c := w.balance
	return &c
}

func (w *Wallet) PendingBalance() *coin.Coin {
	c := w.pendingBalance
	return &c
}

func (w *Wallet) Update() error {

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

func (w *Wallet) QueryBalance(blockNumber *big.Int) error {

	if w.gateway == nil {
		return fmt.Errorf("wallet is not connect to a client")
	}

	balance, err := w.gateway.Client().BalanceAt(context.Background(), w.address, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to query balance from account %s: %v", w.address.Hex(), err)
	}

	err = w.balance.SetWei(*balance)
	if err != nil {
		return fmt.Errorf("failed to store %sWEI as balance requested from account %s: %v", balance.String(), w.address.Hex(), err)
	}

	return nil
}

func (w *Wallet) QueryPendingBalance() error {

	if w.gateway == nil {
		return fmt.Errorf("wallet is not connect to a client")
	}

	balance, err := w.gateway.Client().PendingBalanceAt(context.Background(), w.address)
	if err != nil {
		return fmt.Errorf("failed to query pending balance from account %s: %v", w.address.Hex(), err)
	}

	err = w.pendingBalance.SetWei(*balance)
	if err != nil {
		return fmt.Errorf("failed to store %sWEI as pending balance requested from account %s: %v", balance.String(), w.address.Hex(), err)
	}

	return nil
}

func (w *Wallet) QueryHeader(number *big.Int) (*types.Header, error) {
	header, err := w.gateway.Client().HeaderByNumber(context.Background(), number)
	if err == ethereum.NotFound {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query header by block number %d: %v", number, err)
	}

	return header, nil
}

func (w *Wallet) Block(number *big.Int) (*types.Block, error) {

	block, err := w.gateway.Client().BlockByNumber(context.Background(), number)
	if err != nil {
		return nil, fmt.Errorf("failed to query a block by number %d: %v", number, err)
	}

	return block, nil
}

func NewFromPrivateKey(privateKey ecdsa.PrivateKey) (*Wallet, error) {

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	address := gethcrypto.PubkeyToAddress(*publicKeyECDSA)

	w := &Wallet{}
	w.address = address

	a := walletkey.PrivateKey{}
	a.SetKey(privateKey)
	w.privateKey = a

	b := walletkey.PublicKey{}
	err := b.SetKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error store public key: %v", err)
	}
	w.publicKey = b

	return w, nil
}
