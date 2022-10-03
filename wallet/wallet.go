package wallet

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mauroalderete/weasel/coin"
)

type WalletPrivateKey struct {
	key      ecdsa.PrivateKey
	keyBytes []byte
	keyHex   string
}

func (wpk *WalletPrivateKey) Key() ecdsa.PrivateKey {
	return wpk.key
}

func (wpk *WalletPrivateKey) Bytes() []byte {
	return wpk.keyBytes
}

func (wpk *WalletPrivateKey) Hex() string {
	return wpk.keyHex
}

func (wpk *WalletPrivateKey) SetKey(key ecdsa.PrivateKey) {
	wpk.key = key
	wpk.keyBytes = gethcrypto.FromECDSA(&wpk.key)
	wpk.keyHex = hexutil.Encode(wpk.keyBytes[2:])
}

type WalletPublicKey struct {
	key      crypto.PublicKey
	keyBytes []byte
	keyHex   string
}

func (wpk *WalletPublicKey) Key() crypto.PublicKey {
	return wpk.key
}

func (wpk *WalletPublicKey) Bytes() []byte {
	return wpk.keyBytes
}

func (wpk *WalletPublicKey) Hex() string {
	return wpk.keyHex
}

func (wpk *WalletPublicKey) SetKey(key crypto.PublicKey) error {
	publicKeyECDSA, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	wpk.key = key
	wpk.keyBytes = gethcrypto.FromECDSAPub(publicKeyECDSA)
	wpk.keyHex = hexutil.Encode(wpk.keyBytes[2:])

	return nil
}

type Wallet struct {
	client         *ethclient.Client
	address        common.Address
	privateKey     WalletPrivateKey
	publicKey      WalletPublicKey
	balance        coin.Coin
	pendingBalance coin.Coin
}

func (w *Wallet) Bind(client *ethclient.Client) error {
	if client == nil {
		return fmt.Errorf("client is required")
	}

	w.client = client

	return nil
}

func (w *Wallet) PrivateKey() WalletPrivateKey {
	return w.privateKey
}

func (w *Wallet) PublicKey() WalletPublicKey {
	return w.publicKey
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

	if w.client == nil {
		return fmt.Errorf("wallet is not connect to a client")
	}

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

func (w *Wallet) QueryPendingBalance() error {

	if w.client == nil {
		return fmt.Errorf("wallet is not connect to a client")
	}

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

func NewFromPrivateKey(privateKey ecdsa.PrivateKey) (*Wallet, error) {

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	address := gethcrypto.PubkeyToAddress(*publicKeyECDSA)

	w := &Wallet{}
	w.address = address

	a := WalletPrivateKey{}
	a.SetKey(privateKey)
	w.privateKey = a

	b := WalletPublicKey{}
	err := b.SetKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error store public key: %v", err)
	}
	w.publicKey = b

	return w, nil
}
