package strategy

import (
	"fmt"

	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/market"
	"github.com/ksyoon0321/gotrade/util"
)

type StrategyJump struct {
}

func NewStrategyJump() IStrategy {
	return &StrategyJump{}
}

func (s *StrategyJump) GetStrategyName() string {
	return "StrategyJump"
}

func (s *StrategyJump) Predict(m market.IMarket, c *coin.Coin) coin.BuySellCoin {
	bsCoin := coin.NewBuySellCoin()

	bsCoin.Priority = -1
	bsCoin.StrategyName = s.GetStrategyName()

	hist := c.GetCurrentHistory()
	list := hist.PeekArrayFromTail(2)

	//fmt.Println(" list = ", list, c.GetId())
	if len(list) == 2 {
		prev := list[1].(map[string]interface{})["acc_trade_price_24h"].(float64)
		last := list[0].(map[string]interface{})["acc_trade_price_24h"].(float64)

		if prev <= 0.0 || last <= 0.0 {
			return bsCoin
		}

		//fmt.Println("per = ", fmt.Sprintf("%.4f", (last-prev)/last), c.GetId())
		incPer := (last - prev) / last
		if incPer > 0.005 {
			tradeprice := list[0].(map[string]interface{})["trade_price"].(float64)
			prevtradeprice := list[1].(map[string]interface{})["trade_price"].(float64)

			if prevtradeprice > tradeprice {
				return bsCoin
			}

			bid := util.GetBidPrice(tradeprice)

			//fmt.Println("======= JUMP OVER 0.1% : ", prev, last, c.GetId())
			bsCoin.Priority = 1.4
			bsCoin.StopLoss = getIfSmallerThanV(tradeprice-(bid*5), util.GetAfterPercent(c.GetCurrent(), -3))
			bsCoin.BoughtPrice = 0
			bsCoin.TakeProfit = getIfBiggerThanV(util.GetAfterPercent(c.GetCurrent(), 1), tradeprice+(bid*5))
			bsCoin.WillBuyPrice = util.GetAfterPercent(c.GetCurrent(), 2)

			bsCoin.PrintInfo("ID = " + c.GetId() + ", Current = " + fmt.Sprintf("%.4f", c.GetCurrent()))
		}
	}
	return bsCoin
}
