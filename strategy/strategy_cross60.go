package strategy

import (
	"fmt"

	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/market"
	"github.com/ksyoon0321/gotrade/util"
)

type StrategyCross60 struct {
}

func NewStrategyCross60() IStrategy {
	return &StrategyCross60{}
}

func (s *StrategyCross60) GetStrategyName() string {
	return "StrategyCross60"
}

func (s *StrategyCross60) Predict(m market.IMarket, c *coin.Coin) coin.BuySellCoin {
	bsCoin := coin.NewBuySellCoin()
	bsCoin.Priority = -1
	bsCoin.StrategyName = s.GetStrategyName()

	if c.GetAvgSet().M1 == 0.0 {
		return bsCoin
	}

	if c.GetMaxSet().M60 < c.GetCurrent() && c.GetAvgSet().M60 > c.GetCurrent() {
		bsCoin.StopLoss = getIfSmallerThanV(c.GetMinSet().M5, util.GetAfterPercent(c.GetCurrent(), -3))
		bsCoin.BoughtPrice = 0
		bsCoin.TakeProfit = getIfBiggerThanV(util.GetAfterPercent(c.GetCurrent(), 2), c.GetMaxSet().M60)
		bsCoin.WillBuyPrice = util.GetAfterPercent(c.GetCurrent(), 1)
		bsCoin.Priority = 1.12

		bsCoin.PrintInfo("ID = " + c.GetId() + ", Current = " + fmt.Sprintf("%.4f", c.GetCurrent()))
	}
	return bsCoin
}
