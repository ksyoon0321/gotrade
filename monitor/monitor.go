package monitor

import (
	"strings"
	"time"

	"github.com/ksyoon0321/gotrade/client"
	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/config"
	"github.com/ksyoon0321/gotrade/pubsub"
)

type Monitor struct {
	cfg config.ConfigUpbit

	client  client.IClient
	outChan chan interface{}

	toMarketNotify pubsub.IPubsub

	watchCandles map[string]*coin.Coin
}

func NewMonitor(cfg config.ConfigUpbit, client client.IClient) *Monitor {
	mon := &Monitor{cfg: cfg, client: client, outChan: make(chan interface{}), toMarketNotify: pubsub.NewPubSubChan(), watchCandles: make(map[string]*coin.Coin)}

	return mon
}

func (m *Monitor) Subscribe(id string, ch chan pubsub.PubsubTopic) {
	m.toMarketNotify.Subscribe(id, ch)
}

func (m *Monitor) Run() {
	//get all
	m.getAllCoins()
	go m.fetchCandleRoutine()
}

func (m *Monitor) GetCurrent(id string) {
	go func() {
		for {
			time.Sleep(time.Second * 10)
			item := m.getCoinCurrent(id)

			for _, curr := range item {
				//fmt.Println("curr = ", curr)
				msg := pubsub.NewPubsubTopic("", "REFRESH_CURRENT", curr)
				m.toMarketNotify.Publish(msg)
			}
		}
	}()
}

func (m *Monitor) RegistWatch(id string, c *coin.Coin) {
	if _, ok := m.watchCandles[id]; !ok {
		m.watchCandles[id] = c
		m.getInitCandles(id, c)
	}
}

func (m *Monitor) RemoveWatch(id string) {
	delete(m.watchCandles, id)
}

func (m *Monitor) fetchCandleRoutine() {
	go func() {
		for {
			time.Sleep(time.Minute)

			for id, c := range m.watchCandles {
				m.fetchCandleSet(id, 2, c)
				//fmt.Println("routine : ", id)
			}
		}
	}()
}

func (m *Monitor) getAllCoins() {
	list, err := m.client.GetAllCoins()
	if err != nil {
		panic(err)
	}

	maplist := list.([]map[string]interface{})
	var allCoinIds string
	for _, item := range maplist {
		if strings.HasPrefix(item["market"].(string), m.cfg.Currency) {
			msg := pubsub.NewPubsubTopic("", "INIT", item)
			//fmt.Println("init = ", item["market"].(string))
			m.toMarketNotify.Publish(msg)

			if allCoinIds == "" {
				allCoinIds += item["market"].(string)
			} else {
				allCoinIds += "," + item["market"].(string)
			}

		}
	}

	if allCoinIds != "" {
		item := make(map[string]interface{})
		item["id"] = allCoinIds
		msg := pubsub.NewPubsubTopic("", "START_CURRENT", item)
		m.toMarketNotify.Publish(msg)
	}

}

func (m *Monitor) getInitCandles(id string, c *coin.Coin) {
	//fmt.Println("======= INIT CANDLES : ", id)

	c.InitCandles()
	m.fetchCandleSet(id, 60, c)
}

func (m *Monitor) fetchCandleSet(id string, cnt int, c *coin.Coin) {
	//m1
	candles := m.getCoinCandles(id, coin.MIN1, cnt)
	for ii := len(candles) - 1; ii >= 0; ii-- {
		cndl := candles[ii]
		c.AddCandle(coin.MIN1, coin.NewCandle(cndl))
	}
	//m5
	candles5 := m.getCoinCandles(id, coin.MIN5, cnt)
	for ii := len(candles5) - 1; ii >= 0; ii-- {
		cndl := candles5[ii]
		c.AddCandle(coin.MIN5, coin.NewCandle(cndl))
	}

	//m30
	candles30 := m.getCoinCandles(id, coin.MIN30, cnt)
	for ii := len(candles30) - 1; ii >= 0; ii-- {
		cndl := candles30[ii]
		c.AddCandle(coin.MIN30, coin.NewCandle(cndl))
	}

	//m60
	candles60 := m.getCoinCandles(id, coin.MIN60, cnt)
	for ii := len(candles60) - 1; ii >= 0; ii-- {
		cndl := candles60[ii]
		c.AddCandle(coin.MIN60, coin.NewCandle(cndl))
	}

	//m240
	candles240 := m.getCoinCandles(id, coin.MIN240, cnt)
	for ii := len(candles240) - 1; ii >= 0; ii-- {
		cndl := candles240[ii]
		c.AddCandle(coin.MIN240, coin.NewCandle(cndl))
	}
}

func (m *Monitor) getCoinCandles(id string, min int, cnt int) []map[string]interface{} {
	//fmt.Println(id, min, cnt)
	//get candle
	list, err := m.client.GetCoinCandles(id, min, cnt)
	if err != nil {
		panic(err)
	}

	return list
}

func (m *Monitor) getCoinCurrent(id string) []map[string]interface{} {
	curr, err := m.client.GetCoinCurrent(id)
	if err != nil {
		panic(err)
	}

	return curr
}
