package walletkey

import (
	"crypto"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type PrivateKey struct {
	key      ecdsa.PrivateKey
	keyBytes []byte
	keyHex   string
}

func (wpk *PrivateKey) Key() ecdsa.PrivateKey {
	return wpk.key
}

func (wpk *PrivateKey) Bytes() []byte {
	return wpk.keyBytes
}

func (wpk *PrivateKey) Hex() string {
	return wpk.keyHex
}

func (wpk *PrivateKey) SetKey(key ecdsa.PrivateKey) {
	wpk.key = key
	wpk.keyBytes = gethcrypto.FromECDSA(&wpk.key)
	wpk.keyHex = hexutil.Encode(wpk.keyBytes[2:])
}

type PublicKey struct {
	key      crypto.PublicKey
	keyBytes []byte
	keyHex   string
}

func (wpk *PublicKey) Key() crypto.PublicKey {
	return wpk.key
}

func (wpk *PublicKey) Bytes() []byte {
	return wpk.keyBytes
}

func (wpk *PublicKey) Hex() string {
	return wpk.keyHex
}

func (wpk *PublicKey) SetKey(key crypto.PublicKey) error {
	publicKeyECDSA, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	wpk.key = key
	wpk.keyBytes = gethcrypto.FromECDSAPub(publicKeyECDSA)
	wpk.keyHex = hexutil.Encode(wpk.keyBytes[2:])

	return nil
}
