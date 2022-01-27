package market

import "github.com/ksyoon0321/gotrade/coin"

type IMarket interface {
	GetMarketId() int
	BuildMonitor()

	RegistWatch(id string)
	RemoveWatch(id string)

	Subscribe(string, interface{})

	RequestBuy(id string, bscoin *coin.BuySellCoin)
	RequestSell(id string, c *coin.Coin, bscoin *coin.BuySellCoin) bool
	RequestExpireSell(id string, c *coin.Coin, bscoin *coin.BuySellCoin)

	GetBoughtCoins() []*coin.BuySellCoin
	Run()
}
