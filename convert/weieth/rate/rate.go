package rate

import (
	"math"
	"math/big"
)

var rateFloat64Value float64

func ToFloat() float64 {
	return rateFloat64Value
}

func ToBigFloat() *big.Float {
	return big.NewFloat(rateFloat64Value)
}

func init() {
	rateFloat64Value = math.Pow10(18)
}
