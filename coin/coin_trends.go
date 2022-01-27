package coin

import (
	"github.com/ksyoon0321/gotrade/util"
)

const (
	TRENDS_HIGH = true
	TRENDS_LOW  = false
)

func latestHighLow(q *util.CircleQueue, isHigh bool, cnt int, field string) float64 {
	list := q.PeekArrayFromTail(cnt)

	if len(list) != cnt {
		return 0.0
	}

	var vl float64
	diff := 0
	for ii := len(list) - 1; ii >= 0; ii-- {
		item := list[ii].(map[string]interface{})
		cur := item[field].(float64)

		if vl == 0.0 {
			vl = cur
		} else {
			var isBigger bool
			if cur-vl > 0 {
				//fmt.Println("cur = ", cur, " vl = ", vl, "BIG")
				isBigger = true
			} else {
				//fmt.Println("cur = ", cur, " vl = ", vl, "SMALL")
				isBigger = false
			}

			if isBigger == isHigh {
				vl = cur
				diff = 0
			} else {
				if diff >= 5 {
					return vl
				}
				diff++
			}
		}
	}
	return vl
}

func LatestHigh(q *util.CircleQueue, cnt int) float64 {
	return latestHighLow(q, true, cnt, "trade_price")
}

func LatestLow(q *util.CircleQueue, cnt int) float64 {
	return latestHighLow(q, false, cnt, "trade_price")
}

func IsStillRaising(q *util.CircleQueue, cnt int) bool {
	list := q.PeekArrayFromTail(cnt + 1)

	if len(list) < cnt+1 {
		cnt = len(list) - 1
	}

	if cnt < 1 {
		return false
	}

	incCnt := 0
	prev := 0.0
	for ii := len(list) - 1; ii >= 0; ii-- {
		item := list[ii].(map[string]interface{})
		cur := item["trade_price"].(float64)

		if prev == 0.0 {
			prev = cur
		} else {
			if cur > prev {
				incCnt++
			} else {
				return false
			}
		}
	}

	return cnt == incCnt
}

func RisingV2(q *util.CircleQueue, bscoin *BuySellCoin) bool {
	low := latestHighLow(q, TRENDS_HIGH, 3, "low_price")
	high := latestHighLow(q, TRENDS_LOW, 3, "high_price")

	if high > low {
		bscoin.Priority += 0.5
		bscoin.TakeProfit += util.GetBidPrice(bscoin.WillBuyPrice) * 4
		bscoin.StopLoss += util.GetBidPrice(bscoin.WillBuyPrice) * 2

		return true
	}
	return false
}
