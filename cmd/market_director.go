package cmd

import (
	"container/list"
	"fmt"
	"strconv"
	"time"

	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/config"
	"github.com/ksyoon0321/gotrade/market"
	"github.com/ksyoon0321/gotrade/pubsub"
	"github.com/ksyoon0321/gotrade/strategy"
	"github.com/ksyoon0321/gotrade/txlog"
	"github.com/ksyoon0321/gotrade/util"
)

const (
	EXPIREMIN = 5
)

type MarketDirector struct {
	markets *list.List

	mktNotify chan pubsub.PubsubTopic

	watch         map[string]*coin.BuySellCoin
	strategies    []strategy.IStrategy
	skipCoinCache *util.CacheTimer
	txmgr         *txlog.TxManager
}

func NewMarketDirector(txmanager *txlog.TxManager) *MarketDirector {
	m := &MarketDirector{
		txmgr: txmanager,
	}

	m.markets = list.New()
	m.mktNotify = make(chan pubsub.PubsubTopic)
	m.watch = make(map[string]*coin.BuySellCoin)
	m.skipCoinCache = util.NewCacheTimer(time.Minute)

	return m
}

func (m *MarketDirector) AddMarketConfig(cfg interface{}) {
	m.markets.PushBack(m.marketFactory(cfg))
}

func (m *MarketDirector) RegistStrategy(s strategy.IStrategy) {
	m.strategies = append(m.strategies, s)
}

func (m *MarketDirector) removeWatchAndPutSkip(market market.IMarket, coinid, id string) {
	delete(m.watch, id)

	market.RemoveWatch(coinid)
	m.skipCoinCache.Put(id, true, time.Minute*3)
}
func (m *MarketDirector) IsContainWatch(id string) bool {
	if _, ok := m.watch[id]; ok {
		return true
	}
	return false
}

func (m *MarketDirector) GetMarketById(id int) market.IMarket {
	for mkt := m.markets.Front(); mkt != nil; mkt = mkt.Next() {
		if (mkt.Value.(market.IMarket)).GetMarketId() == id {
			return mkt.Value.(market.IMarket)
		}
	}

	return nil
}

func (m *MarketDirector) Run() {
	go func() {
		for data := range m.mktNotify {
			c := data.GetData().(*coin.Coin)
			watchId := data.GetId() + "_" + c.GetId()

			marketId, err := strconv.Atoi(data.GetId())
			if err != nil {
				continue
			}

			msgMarket := m.GetMarketById(marketId)
			if msgMarket == nil {
				continue
			}

			//fmt.Println("========= CURRENT = ", c.GetCurrent(), watchId)
			//strategy
			if m.IsContainWatch(watchId) {
				bsCoin := m.watch[watchId]
				if bsCoin.BoughtPrice == 0.0 {
					if !bsCoin.IsRequest {
						alreadyOverPrice := bsCoin.WillBuyPrice + util.GetBidPrice(c.GetCurrent())*2
						if alreadyOverPrice <= c.GetCurrent() || bsCoin.TakeProfit <= c.GetCurrent() {

							act := fmt.Sprintf("REMOVE_OVER:CURR=%.4f:WILLBUY=%.4f:OVER=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.WillBuyPrice, alreadyOverPrice, bsCoin.StrategyName)
							m.txmgr.Finish(watchId, txlog.NewTxLog(act))

							m.removeWatchAndPutSkip(msgMarket, c.GetId(), watchId)
						} else if bsCoin.WillBuyPrice >= c.GetCurrent() {
							//buy
							if bsCoin.Priority >= 1 {
								act := fmt.Sprintf("BUY:CURR=%.4f:WILLBUY=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.WillBuyPrice, bsCoin.StrategyName)
								m.txmgr.Push(watchId, txlog.NewTxLog(act))

								bsCoin.WillBuyPrice = c.GetCurrent()
								msgMarket.RequestBuy(c.GetId(), bsCoin)
							}
						}
					}
				} else {
					//expire
					if bsCoin.BoughtTime.Add(time.Minute * EXPIREMIN).Before(time.Now()) {
						act := fmt.Sprintf("TRYEXPIRE:CURR=%.4f:TAKEPROFIT=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.TakeProfit, bsCoin.StrategyName)
						m.txmgr.Push(watchId, txlog.NewTxLog(act))

						bsCoin.BoughtTime = time.Now()
						msgMarket.RequestExpireSell(c.GetId(), c, bsCoin)
						if bsCoin.BoughtExtend > 1 {
							act := fmt.Sprintf("FINEXPIRE:CURR=%.4f:TAKEPROFIT=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.TakeProfit, bsCoin.StrategyName)
							m.txmgr.Finish(watchId, txlog.NewTxLog(act))

							m.removeWatchAndPutSkip(msgMarket, c.GetId(), watchId)
						} else {
							act := fmt.Sprintf("EXPIRED:CURR=%.4f:TAKEPROFIT=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.TakeProfit, bsCoin.StrategyName)
							m.txmgr.Push(watchId, txlog.NewTxLog(act))
						}

					} else {
						if bsCoin.TakeProfit <= c.GetCurrent() {
							//sell
							fin := msgMarket.RequestSell(c.GetId(), c, bsCoin)
							if fin {
								benefit := ((bsCoin.TakeProfit - bsCoin.BoughtPrice) / c.GetCurrent()) * 100
								act := fmt.Sprintf("SELL:CURR=%.4f:TAKEPROFIT=%.4f:STRATEGY=%s:BENEFIT=%.4f", c.GetCurrent(), bsCoin.TakeProfit, bsCoin.StrategyName, benefit)
								m.txmgr.Finish(watchId, txlog.NewTxLog(act))

								m.removeWatchAndPutSkip(msgMarket, c.GetId(), watchId)
							} else {
								act := fmt.Sprintf("INCSELL:CURR=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.StrategyName)
								m.txmgr.Push(watchId, txlog.NewTxLog(act))
							}
						} else if bsCoin.StopLoss >= c.GetCurrent() {
							//stop loss
							act := fmt.Sprintf("LOSSCUT:CURR=%.4f:STOPLOSS=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.StopLoss, bsCoin.StrategyName)
							m.txmgr.Finish(watchId, txlog.NewTxLog(act))

							m.removeWatchAndPutSkip(msgMarket, c.GetId(), watchId)
						} else {
							list := c.GetCurrentHistory().PeekArrayFromTail(2)
							if len(list) >= 2 {
								prev := list[1].(map[string]interface{})
								if c.GetCurrent() >= prev["trade_price"].(float64) {
									bsCoin.BoughtTime = time.Now()

									act := fmt.Sprintf("EXTENDTIME:CURR=%.4f:STRATEGY=%s", c.GetCurrent(), bsCoin.StrategyName)
									m.txmgr.Push(watchId, txlog.NewTxLog(act))
								}
							}
						}
					}
				}
			} else {
				list := c.GetCurrentHistory().PeekArrayFromTail(2)
				if len(list) >= 2 {
					prev := list[1].(map[string]interface{})
					if prev["trade_price"].(float64) <= c.GetCurrent() {
						continue
					}

					switch data.GetCmd() {
					case "REFRESH_CURRENT":
						if m.skipCoinCache.Get(watchId) == nil {
							for _, st := range m.strategies {
								bsCoin := st.Predict(msgMarket, c)

								needInit := m.RegistWatch(watchId, &bsCoin)
								if needInit {
									msgMarket.RegistWatch(c.GetId())
								}
							}
						}
					}
				}
			}
		}
	}()

	//확인용 구매코인 목록 표시
	// go func() {
	// 	for {
	// 		time.Sleep(time.Minute * 5)
	// 		fmt.Println("========================================================================================")
	// 		boughtCount := 0
	// 		for key, item := range m.watch {
	// 			if item.IsRequest {
	// 				fmt.Println(" ============COIN = ", key, " , item = ", item)
	// 				boughtCount++
	// 			}
	// 		}

	// 		if boughtCount == 0 {
	// 			fmt.Println("=====NOT FOUND TARGET")
	// 		}

	// 		fmt.Println("========================================================================================")
	// 	}
	// }()

	for itm := m.markets.Front(); itm != nil; itm = itm.Next() {
		(itm.Value.(market.IMarket)).Subscribe("dc", m.mktNotify)
		(itm.Value.(market.IMarket)).Run()
	}

	for {
		time.Sleep(time.Second)
	}
}

func (m *MarketDirector) RegistWatch(id string, bsCoin *coin.BuySellCoin) bool {
	if bsCoin.Priority <= 0 {
		return false
	}

	if m.IsContainWatch(id) {
		oldBsCoin := m.watch[id]
		if oldBsCoin.Priority < bsCoin.Priority {
			//fmt.Println(" ==== REPLACE BSCOIN = ", id, bsCoin)
			m.watch[id] = bsCoin
		}
		return false
	} else {
		m.watch[id] = bsCoin
		//fmt.Println(" ==== NEW BSCOIN = ", id, bsCoin)

		return true
	}
}

//private
func (m *MarketDirector) marketFactory(cfg interface{}) market.IMarket {
	if cfg == nil {
		panic("nil configuration")
	}

	var mkt market.IMarket
	switch cfg.(type) {
	case config.ConfigUpbit:
		mkt = market.NewUpbitMarket(m.markets.Len(), cfg.(config.ConfigUpbit))
	default:
		panic("not support market type")
	}

	mkt.BuildMonitor()

	return mkt
}
