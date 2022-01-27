package strategy

import (
	"fmt"

	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/market"
	"github.com/ksyoon0321/gotrade/util"
)

type StrategyTrends struct {
}

func NewStrategyTrends() IStrategy {
	return &StrategyTrends{}
}

func (s *StrategyTrends) GetStrategyName() string {
	return "StrategyTrends"
}

func (s *StrategyTrends) Predict(m market.IMarket, c *coin.Coin) coin.BuySellCoin {
	bsCoin := coin.NewBuySellCoin()

	bsCoin.Priority = -1
	bsCoin.StrategyName = s.GetStrategyName()

	hist := c.GetCurrentHistory()
	list := hist.PeekArrayFromTail(4)

	//fmt.Println(" list = ", list, c.GetId())
	if len(list) == 4 {
		rasing := coin.IsStillRaising(c.GetCurrentHistory(), 3)

		rasing = true
		if rasing {
			tradeprice := list[0].(map[string]interface{})["trade_price"].(float64)
			bid := util.GetBidPrice(tradeprice)

			bsCoin.Priority = 0.9
			bsCoin.StopLoss = tradeprice - (bid * 50)
			bsCoin.BoughtPrice = 0
			bsCoin.TakeProfit = tradeprice + (bid * 100)
			bsCoin.WillBuyPrice = tradeprice + (bid * 200)

			bsCoin.PrintInfo("ID = " + c.GetId() + ", Current = " + fmt.Sprintf("%.4f", c.GetCurrent()))
		}
	}
	return bsCoin
}
