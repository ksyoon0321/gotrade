package strategy

import (
	"fmt"

	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/market"
	"github.com/ksyoon0321/gotrade/util"
)

type StrategyCrossHigh struct {
}

func NewStrategyCrossHigh() IStrategy {
	return &StrategyCrossHigh{}
}

func (s *StrategyCrossHigh) GetStrategyName() string {
	return "NewStrategyCrossHigh"
}

func (s *StrategyCrossHigh) Predict(m market.IMarket, c *coin.Coin) coin.BuySellCoin {
	bsCoin := coin.NewBuySellCoin()
	bsCoin.Priority = -1
	bsCoin.StrategyName = s.GetStrategyName()

	if c.GetAvgSet().M1 == 0.0 {
		return bsCoin
	}

	high := coin.LatestHigh(c.GetCurrentHistory(), 10)
	if high <= 0.0 {
		return bsCoin
	}

	//latestHigh 함수는 요청갯수 만큼의 데이터가 없으면 0을 반환하므로 최소 10개 이상 있는 상황으로
	//배열 갯수 검사 필요 없다.
	list := c.GetCurrentHistory().PeekArrayFromTail(3)
	prev := list[2].(map[string]interface{})
	if prev["trade_price"].(float64) > high {
		return bsCoin
	}

	if high < c.GetCurrent() {
		bsCoin.StopLoss = getIfSmallerThanV(c.GetMinSet().M5, util.GetAfterPercent(c.GetCurrent(), -5))
		bsCoin.BoughtPrice = 0
		bsCoin.TakeProfit = getIfBiggerThanV(util.GetAfterPercent(c.GetCurrent(), 3), c.GetMaxSet().M30)
		bsCoin.WillBuyPrice = util.GetAfterPercent(c.GetCurrent(), 2)
		bsCoin.Priority = 1.13

		bsCoin.PrintInfo("ID = " + c.GetId() + ", Current = " + fmt.Sprintf("%.4f", c.GetCurrent()))
	}
	return bsCoin
}
