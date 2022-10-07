package repository

import (
	"fmt"
	"io"

	"github.com/mauroalderete/weasel/wallet"
)

type Repository interface {
	Save(wallet.Wallet) error
	Close()
}

type Base struct {
	write io.Writer
}

func (r *Base) Save(w wallet.Wallet) error {
	payload, err := wallet.JsonMarshal(w)
	if err != nil {
		return fmt.Errorf("failed to prepare the payload to save the wallet: %v", err)
	}

	total := len(payload)
	n := 0

	for n < total {
		n, err = r.write.Write(payload[n:])
		if err != nil {
			return fmt.Errorf("failed to save the wallet from %d bytes to forward: %v", n, err)
		}
	}

	return nil
}

func (r *Base) Close() {
}

func New(write io.Writer) (*Base, error) {

	if write == nil {
		return nil, fmt.Errorf("a writer instanced is required")
	}

	r := Base{}

	r.write = write

	return &r, nil
}
