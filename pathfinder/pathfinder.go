package pathfinder

import (
	"fmt"

	"github.com/mauroalderete/weasel/client"
	"github.com/mauroalderete/weasel/genwallet"
	"github.com/mauroalderete/weasel/wallet"
)

type Pathfinder struct {
	client *client.Client
	wallet *wallet.Wallet
	match  bool
}

func (p *Pathfinder) Connect(gateway string) error {
	p.client = &client.Client{}

	err := p.client.Connect(gateway)
	if err != nil {
		return fmt.Errorf("failed to connect to ethclient: %v", err)
	}

	return nil
}

func (p *Pathfinder) Close() {
	if p.client != nil {
		p.client.Close()
	}
}

func (p *Pathfinder) Search() error {
	w, err := genwallet.RandomWallet()
	if err != nil {
		return fmt.Errorf("failed instance a new random wallet: %v", err)
	}

	err = w.Bind(p.client)
	if err != nil {
		return fmt.Errorf("failed to bind client to wallet: %v", err)
	}

	err = w.Update()
	if err != nil {
		return fmt.Errorf("failed to update wallet info: %v", err)
	}

	p.wallet = w

	p.match, err = p.searchMatch()
	if err != nil {
		return fmt.Errorf("failed to comprobate if the wallet match: %v", err)
	}

	return nil
}

func (p *Pathfinder) Wallet() *wallet.Wallet {
	w := *p.wallet
	return &w
}

func (p *Pathfinder) Match() bool {
	return p.match
}

func (p *Pathfinder) searchMatch() (bool, error) {

	if p.wallet == nil {
		return false, fmt.Errorf("there is not a wallet found to comprobate if match or not")
	}

	b := p.wallet.Balance().Wei()
	pb := p.wallet.PendingBalance().Wei()

	return b.Int64() != 0 || pb.Int64() != 0, nil
}
