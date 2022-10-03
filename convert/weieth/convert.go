package weieth

import (
	"fmt"
	"math/big"

	"github.com/mauroalderete/weasel/convert/weieth/rate"
)

func WeiToEth(wei big.Int) (*big.Float, error) {
	fwei := &big.Float{}
	_, ok := fwei.SetString(wei.String())
	if !ok {
		return nil, fmt.Errorf("failed convert %sWEI to ETH", wei.String())
	}

	eth := &big.Float{}
	eth.Quo(fwei, rate.ToBigFloat())

	return eth, nil
}

func EthToWei(eth big.Float) (*big.Int, error) {

	fwei := big.Float{}
	fwei.Mul(&eth, rate.ToBigFloat())

	iwei, _ := fwei.Int64()
	wei := big.NewInt(iwei)

	return wei, nil
}
