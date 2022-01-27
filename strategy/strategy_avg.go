package strategy

import (
	"fmt"
	"strconv"

	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/market"
	"github.com/ksyoon0321/gotrade/util"
)

type StrategyAvg struct {
	mintype int
}

func NewStrategyAvg(min int) IStrategy {
	return &StrategyAvg{
		mintype: min,
	}
}

func (s *StrategyAvg) GetStrategyName() string {
	return "StrategyAvg" + strconv.Itoa(s.mintype)
}
func (s *StrategyAvg) Predict(m market.IMarket, c *coin.Coin) coin.BuySellCoin {
	bsCoin := coin.NewBuySellCoin()
	bsCoin.Priority = -1
	bsCoin.StrategyName = s.GetStrategyName()

	if c.GetAvgSet().M1 == 0.0 {
		return bsCoin
	}

	//fmt.Println("======== AVG : ", c.GetAvgSet().M60, c.GetAvgSet().M5)
	//going up
	q := c.GetCandleQueue(s.mintype)
	var avgvl float64
	var prior float64
	switch s.mintype {
	case coin.MIN1:
		avgvl = c.GetAvgSet().M1
		prior = 1.01
	case coin.MIN5:
		avgvl = c.GetAvgSet().M5
		prior = 1.02
	case coin.MIN30:
		avgvl = c.GetAvgSet().M30
		prior = 1.1
	case coin.MIN60:
		avgvl = c.GetAvgSet().M60
		prior = 1.11
	case coin.MIN240:
		avgvl = c.GetAvgSet().M240
		prior = 1.2
	}

	if avgvl == 0 {
		return bsCoin
	}

	list := q.PeekArrayFromTail(2)

	if len(list) < 2 {
		return bsCoin
	}

	item := list[1].(coin.Candle)
	avgPrev := item.GetTradePrice()

	if c.GetCurrent() > avgvl && avgvl > avgPrev {
		//fmt.Println("++++++++++", c.GetId(), "AVG ", s.mintype, ": ", util.GetPrintFloat64(avgvl), ", CURRENT : ", util.GetPrintFloat64(c.GetCurrent()), " PREV : ", util.GetPrintFloat64(avgPrev))

		bid := util.GetBidPrice(c.GetCurrent())
		bsCoin.StopLoss = getIfSmallerThanV(c.GetCurrent()-(bid*10), util.GetAfterPercent(c.GetCurrent(), -3))
		bsCoin.BoughtPrice = 0
		bsCoin.TakeProfit = getIfBiggerThanV(util.GetAfterPercent(c.GetCurrent(), 2), c.GetCurrent()+(bid*10))
		bsCoin.WillBuyPrice = util.GetAfterPercent(c.GetCurrent(), 1)
		bsCoin.Priority = prior

		bsCoin.PrintInfo("ID = " + c.GetId() + ", Current = " + fmt.Sprintf("%.4f", c.GetCurrent()))
	}
	return bsCoin
}

func getIfBiggerThanV(v, compareV float64) float64 {
	if v > compareV {
		return v
	}
	return compareV
}

func getIfSmallerThanV(v, compareV float64) float64 {
	if v > compareV {
		return compareV
	}
	return v
}
