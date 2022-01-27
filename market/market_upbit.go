package market

import (
	"strconv"
	"time"

	"github.com/ksyoon0321/gotrade/client"
	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/config"
	"github.com/ksyoon0321/gotrade/monitor"
	"github.com/ksyoon0321/gotrade/pubsub"
)

type MarketUpbit struct {
	id   int
	conf config.ConfigUpbit

	client client.IClient
	mon    *monitor.Monitor

	monNotify chan pubsub.PubsubTopic

	toDirectorNotify pubsub.IPubsub

	coins map[string]*coin.Coin
}

func NewUpbitMarket(id int, cfg config.ConfigUpbit) IMarket {
	return &MarketUpbit{
		id:               id,
		conf:             cfg,
		client:           client.NewUpbitRestClient(cfg),
		monNotify:        make(chan pubsub.PubsubTopic),
		toDirectorNotify: pubsub.NewPubSubChan(),
		coins:            make(map[string]*coin.Coin)}
}

func (u *MarketUpbit) GetMarketId() int {
	return u.id
}

func (u *MarketUpbit) BuildMonitor() {
	u.mon = monitor.NewMonitor(u.conf, u.client)
}

func (u *MarketUpbit) RegistWatch(id string) {
	u.mon.RegistWatch(id, u.coins[id])
}

func (u *MarketUpbit) RemoveWatch(id string) {
	u.mon.RemoveWatch(id)
}

func (u *MarketUpbit) Subscribe(id string, ch interface{}) {
	u.toDirectorNotify.Subscribe(id, ch.(chan pubsub.PubsubTopic))
}

func (u *MarketUpbit) Run() {
	u.mon.Subscribe(u.conf.Id, u.monNotify)

	go func() {
		for data := range u.monNotify {
			willNotify := u.parseCommand(data)

			if willNotify {
				d, ok := data.GetData().(map[string]interface{})
				if ok {
					id := d["market"].(string)
					coinData := pubsub.NewPubsubTopic(strconv.Itoa(u.id), data.GetCmd(), u.coins[id])
					u.toDirectorNotify.Publish(coinData)
				}
			}
		}
	}()

	u.mon.Run()
}

func (u *MarketUpbit) RequestBuy(id string, bscoin *coin.BuySellCoin) {
	//fmt.Println(" ======== Buy : ", id, bscoin.WillBuyPrice)

	bscoin.IsRequest = true

	bscoin.BoughtTime = time.Now()
	//어차피 시장가로 긁기 때문에 매수가는 다시 조회해야 함
	bscoin.BoughtPrice = bscoin.WillBuyPrice
}

func (u *MarketUpbit) RequestSell(id string, c *coin.Coin, bscoin *coin.BuySellCoin) bool {
	//손절가, 매도가 수정이 필요한가?
	//prevPrefit := bscoin.TakeProfit
	need := c.NeedIncreaseSellPrice(bscoin)

	if need {
		//fmt.Println(" ======== Increase Sell Price : ", fmt.Sprintf("%.4f => %.4f", prevPrefit, bscoin.TakeProfit))
		return false
	} else {
		//시장가로 던지기 때문에 총 손익은 계정에서 조회해야 함.
		//benefit := ((bscoin.TakeProfit - bscoin.BoughtPrice) / c.GetCurrent()) * 100
		//fmt.Println(" ======== Sell : ", id, bscoin.TakeProfit, " Benefit = ", fmt.Sprintf("%.4f", benefit))

		return true
	}
}

func (u *MarketUpbit) RequestExpireSell(id string, c *coin.Coin, bscoin *coin.BuySellCoin) {
	bscoin.BoughtExtend++

	//fmt.Println("========== RequestExpireSell : ", c.GetId(), bscoin.BoughtExtend)
	if bscoin.BoughtExtend > 1 {
		//fmt.Println("++++++++++ RequestExpireSell : ", c.GetId(), bscoin.BoughtExtend)
		//시장가로 던지기 때문에 총 손익은 계정에서 조회해야 함.
		//benefit := ((bscoin.TakeProfit - bscoin.BoughtPrice) / c.GetCurrent()) * 100
		//fmt.Println(" ======== Sell : ", id, bscoin.TakeProfit, " Benefit = ", fmt.Sprintf("%.4f", benefit))
	}
}

func (u *MarketUpbit) GetBoughtCoins() []*coin.BuySellCoin {
	return nil
}

//private
func (u *MarketUpbit) parseCommand(data pubsub.PubsubTopic) bool {
	willNotifyToDirector := false

	d, ok := data.GetData().(map[string]interface{})
	if !ok {
		return false
	}

	//fmt.Println("parseCommand = ", d["market"].(string))
	switch data.GetCmd() {
	case "INIT":
		id := d["market"].(string)
		kornm := d["korean_name"].(string)
		if _, ok := u.coins[id]; ok {
			return false
		}

		newcoin := coin.NewCoin(id, kornm)
		u.coins[id] = newcoin
	case "START_CURRENT":
		//fmt.Println("============ START CURRENT")
		u.mon.GetCurrent(d["id"].(string))
	case "REFRESH_CURRENT":
		//fmt.Println("============ REFRESH CURRENT")
		u.coins[d["market"].(string)].SetCurrent(d)

		if d["acc_trade_price_24h"].(float64) > u.conf.MinTradePrice {
			willNotifyToDirector = true
		}
	}

	return willNotifyToDirector
}
