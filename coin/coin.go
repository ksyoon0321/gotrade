package coin

import (
	"time"

	"github.com/ksyoon0321/gotrade/util"
)

const (
	MIN1   = 1
	MIN5   = 5
	MIN30  = 30
	MIN60  = 60
	MIN240 = 240
)

type BuySellCoin struct {
	TakeProfit   float64
	StopLoss     float64
	WillBuyPrice float64
	BoughtTime   time.Time
	BoughtExtend int
	BoughtPrice  float64
	BoughtQty    float64
	Priority     float64
	IsRequest    bool
	StrategyName string
}

func NewBuySellCoin() BuySellCoin {
	return BuySellCoin{}
}

func (b *BuySellCoin) PrintInfo(title string) {
	// fmt.Println("++++++++", title, " : Strategy = ", b.StrategyName,
	// 	", WillBuy = ", util.GetPrintFloat64(b.WillBuyPrice),
	// 	", BoughtPrice = ", util.GetPrintFloat64(b.BoughtPrice),
	// 	", Stoploss = ", util.GetPrintFloat64(b.StopLoss),
	// 	", Profit = ", util.GetPrintFloat64(b.TakeProfit))
}

type Candle struct {
	time_kst      string
	opening_price float64 //시가
	high_price    float64 //고가
	low_price     float64 //저가
	trade_price   float64 //종가

	acc_trade_price  float64
	acc_trade_volumn float64
}

func (c *Candle) GetTradePrice() float64 {
	return c.trade_price
}

type CoinMinuteSet struct {
	M1   float64
	M5   float64
	M30  float64
	M60  float64
	M240 float64
}

type Coin struct {
	id       string
	name_kor string

	m1Candle   *util.CircleQueue
	m5Candle   *util.CircleQueue
	m30Candle  *util.CircleQueue
	m60Candle  *util.CircleQueue
	m240Candle *util.CircleQueue

	currentHistory *util.CircleQueue

	current float64
	avgset  CoinMinuteSet
	sumset  CoinMinuteSet
	minset  CoinMinuteSet
	maxset  CoinMinuteSet

	isProcessing bool
}

func NewCandle(item map[string]interface{}) Candle {
	return Candle{
		time_kst:      item["candle_date_time_kst"].(string),
		opening_price: item["opening_price"].(float64),
		high_price:    item["high_price"].(float64),
		low_price:     item["low_price"].(float64),
		trade_price:   item["trade_price"].(float64),

		acc_trade_price:  item["candle_acc_trade_price"].(float64),
		acc_trade_volumn: item["candle_acc_trade_volume"].(float64),
	}
}

func NewCoin(id, name_kor string) *Coin {
	return &Coin{
		id:       id,
		name_kor: name_kor,

		m1Candle:   util.NewCircleQueue(120),
		m5Candle:   util.NewCircleQueue(120),
		m30Candle:  util.NewCircleQueue(120),
		m60Candle:  util.NewCircleQueue(120),
		m240Candle: util.NewCircleQueue(120),

		currentHistory: util.NewCircleQueue(50),
		isProcessing:   false,
	}
}

func (c *Coin) GetId() string {
	return c.id
}

func (c *Coin) IsDoingSomething() bool {
	return c.isProcessing
}

func (c *Coin) startDoSomething() {
	for c.isProcessing {
		time.Sleep(time.Second / 10)
	}

	c.isProcessing = true
}

func (c *Coin) finishDoSomething() {
	c.isProcessing = false
}

func (c *Coin) GetCandleQueue(min int) *util.CircleQueue {
	switch min {
	case MIN1:
		return c.m1Candle
	case MIN5:
		return c.m5Candle
	case MIN30:
		return c.m30Candle
	case MIN60:
		return c.m60Candle
	case MIN240:
		return c.m240Candle
	}
	return nil
}

func (c *Coin) AddCandle(min int, candle Candle) {
	//check tail candle
	c.startDoSomething()
	defer c.finishDoSomething()

	isUpdated := c.isUpdatedTailCandle(min, candle)

	if isUpdated {
		return
	}

	var cndl interface{}
	switch min {
	case MIN1:
		c.sumset.M1 += candle.trade_price
		cndl = c.m1Candle.Enqueue(candle)
	case MIN5:
		c.sumset.M5 += candle.trade_price
		cndl = c.m5Candle.Enqueue(candle)
	case MIN30:
		c.sumset.M30 += candle.trade_price
		cndl = c.m30Candle.Enqueue(candle)
	case MIN60:
		c.sumset.M60 += candle.trade_price
		cndl = c.m60Candle.Enqueue(candle)
	case MIN240:
		c.sumset.M240 += candle.trade_price
		cndl = c.m240Candle.Enqueue(candle)
	}
	c.minusSum(min, cndl)
	c.calcAvg(min)

	c.setMinSet(min, candle)
	c.setMaxSet(min, candle)
}

func (c *Coin) isUpdatedTailCandle(min int, candle Candle) bool {
	var list *util.CircleQueue

	switch min {
	case MIN1:
		list = c.m1Candle
	case MIN5:
		list = c.m5Candle
	case MIN30:
		list = c.m30Candle
	case MIN60:
		list = c.m60Candle
	case MIN240:
		list = c.m240Candle
	}

	if list == nil {
		return false
	}

	if list.Count() == 0 {
		return false
	}

	v := list.PeekTail()
	if v == nil {
		return false
	}

	tail := v.(Candle)

	if tail.time_kst < candle.time_kst {
		return false
	}

	if tail.time_kst == candle.time_kst {
		list.OverwriteTail(candle)

		c.minusSum(min, tail)
		c.calcAvg(min)
		return true
	}

	return false
}

func (c *Coin) PeekArrayFromTail(min, cnt int) []Candle {
	c.startDoSomething()
	defer c.finishDoSomething()

	var arrPeek []interface{}
	switch min {
	case MIN1:
		arrPeek = c.m1Candle.PeekArrayFromTail(cnt)
	case MIN5:
		arrPeek = c.m5Candle.PeekArrayFromTail(cnt)
	case MIN30:
		arrPeek = c.m30Candle.PeekArrayFromTail(cnt)
	case MIN60:
		arrPeek = c.m60Candle.PeekArrayFromTail(cnt)
	case MIN240:
		arrPeek = c.m240Candle.PeekArrayFromTail(cnt)
	}

	if arrPeek == nil {
		return nil
	}

	//fmt.Println("arr Peek = ", arrPeek)
	retPeek := make([]Candle, len(arrPeek))
	for idx, item := range arrPeek {
		retPeek[idx] = item.(Candle)
	}

	return retPeek
}

func (c *Coin) InitCandles() {
	c.startDoSomething()
	defer c.finishDoSomething()

	c.m1Candle.InitQueue()
	c.m5Candle.InitQueue()
	c.m30Candle.InitQueue()
	c.m60Candle.InitQueue()
	c.m240Candle.InitQueue()

	c.avgset = CoinMinuteSet{}
	c.sumset = CoinMinuteSet{}
	c.minset = CoinMinuteSet{}
	c.maxset = CoinMinuteSet{}
}

func (c *Coin) SetCurrent(d map[string]interface{}) {
	c.current = d["trade_price"].(float64)
	c.currentHistory.Enqueue(d)
}

func (c *Coin) GetAvgSet() CoinMinuteSet {
	return c.avgset
}

func (c *Coin) GetMinSet() CoinMinuteSet {
	return c.minset
}

func (c *Coin) GetMaxSet() CoinMinuteSet {
	return c.maxset
}

func (c *Coin) GetCurrent() float64 {
	return c.current
}

func (c *Coin) GetCurrentHistory() *util.CircleQueue {
	return c.currentHistory
}

//private
func (c *Coin) setMinSet(min int, cndl Candle) {
	v := cndl.trade_price

	switch min {
	case MIN1:
		if v < c.minset.M1 || c.minset.M1 == 0.0 {
			c.minset.M1 = v
		}
	case MIN5:
		if v < c.minset.M5 || c.minset.M5 == 0.0 {
			c.minset.M5 = v
		}
	case MIN30:
		if v < c.minset.M30 || c.minset.M30 == 0.0 {
			c.minset.M30 = v
		}
	case MIN60:
		if v < c.minset.M60 || c.minset.M60 == 0.0 {
			c.minset.M60 = v
		}
	case MIN240:
		if v < c.minset.M240 || c.minset.M240 == 0.0 {
			c.minset.M240 = v
		}
	}
}

func (c *Coin) setMaxSet(min int, cndl Candle) {
	v := cndl.trade_price

	switch min {
	case MIN1:
		if v > c.maxset.M1 {
			c.maxset.M1 = v
		}
	case MIN5:
		if v > c.maxset.M5 {
			c.maxset.M5 = v
		}
	case MIN30:
		if v > c.maxset.M30 {
			c.maxset.M30 = v
		}
	case MIN60:
		if v > c.maxset.M60 {
			c.maxset.M60 = v
		}
	case MIN240:
		if v > c.maxset.M240 {
			c.maxset.M240 = v
		}
	}
}

func (c *Coin) minusSum(min int, cndl interface{}) {
	if cndl == nil {
		return
	}

	switch min {
	case MIN1:
		c.sumset.M1 -= cndl.(Candle).trade_price
	case MIN5:
		c.sumset.M5 -= cndl.(Candle).trade_price
	case MIN30:
		c.sumset.M30 -= cndl.(Candle).trade_price
	case MIN60:
		c.sumset.M60 -= cndl.(Candle).trade_price
	case MIN240:
		c.sumset.M240 -= cndl.(Candle).trade_price
	}
}

func (c *Coin) calcAvg(min int) {
	var count int
	var sum float64
	switch min {
	case MIN1:
		count = c.m1Candle.Count()
		sum = c.sumset.M1
	case MIN5:
		count = c.m5Candle.Count()
		sum = c.sumset.M5
	case MIN30:
		count = c.m30Candle.Count()
		sum = c.sumset.M30
	case MIN60:
		count = c.m60Candle.Count()
		sum = c.sumset.M60
	case MIN240:
		count = c.m240Candle.Count()
		sum = c.sumset.M240
	}

	switch min {
	case MIN1:
		c.avgset.M1 = sum / float64(count)
	case MIN5:
		c.avgset.M5 = sum / float64(count)
	case MIN30:
		c.avgset.M30 = sum / float64(count)
	case MIN60:
		c.avgset.M60 = sum / float64(count)
	case MIN240:
		c.avgset.M240 = sum / float64(count)
	}
}

func (c *Coin) NeedIncreaseSellPrice(bscoin *BuySellCoin) bool {
	return RisingV2(c.currentHistory, bscoin)
}
