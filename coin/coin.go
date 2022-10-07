package coin

import (
	"math/big"

	"github.com/mauroalderete/weasel/convert/weieth"
)

type CoinMarshable struct {
	Wei int64   `json:"wei"`
	Eth float64 `json:"eth"`
}

func (cm *CoinMarshable) Coin() Coin {

	c := Coin{}

	wei := big.NewInt(cm.Wei)
	c.SetWei(*wei)

	return c
}

type Coin struct {
	wei big.Int
	eth big.Float
}

func (c *Coin) Wei() *big.Int {
	wei := c.wei
	return &wei
}

func (c *Coin) Eth() *big.Float {
	eth := c.eth
	return &eth
}

func (c *Coin) SetWei(wei big.Int) error {
	eth, err := weieth.WeiToEth(wei)
	if err != nil {
		return err
	}

	c.wei = wei
	c.eth = *eth

	return nil
}

func (c *Coin) SetEth(eth big.Float) error {
	wei, err := weieth.EthToWei(eth)
	if err != nil {
		return err
	}

	c.wei = *wei
	c.eth = eth

	return nil
}

func (c *Coin) Marshable() CoinMarshable {
	wei := c.wei.Int64()
	eth, _ := c.eth.Float64()

	cm := CoinMarshable{
		Wei: wei,
		Eth: eth,
	}

	return cm
}
