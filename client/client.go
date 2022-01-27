package client

type IClient interface {
	Call(param interface{}) (interface{}, error)

	GetAllCoins() (interface{}, error)
	GetCoinCandles(id string, min int, cnt int) ([]map[string]interface{}, error)
	GetCoinCurrent(id string) ([]map[string]interface{}, error)
}
