package strategy

import (
	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/market"
)

type IStrategy interface {
	Predict(m market.IMarket, c *coin.Coin) coin.BuySellCoin
	GetStrategyName() string
}
