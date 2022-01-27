package config

type ConfigUpbit struct {
	Accesskey     string
	Secretkey     string
	Apiurl        string
	Apiver        string
	Id            string
	MinTradePrice float64
	Currency      string

	MinimumBuy      int
	MaximumBuy      int
	MaximumBuyCount int
}

var CfgUpbit ConfigUpbit

func init() {

	CfgUpbit.Accesskey = "api"
	CfgUpbit.Secretkey = "sec"
	CfgUpbit.Apiurl = "https://api.upbit.com"
	CfgUpbit.Apiver = "v1"
	CfgUpbit.Id = "UPBIT_1"
	CfgUpbit.MinTradePrice = 3000000000
	CfgUpbit.MinTradePrice = 100000
	CfgUpbit.Currency = "KRW"
	//CfgUpbit.Currency = "KRW-WEMIX"

	CfgUpbit.MinimumBuy = 1000
	CfgUpbit.MaximumBuy = 5000
	CfgUpbit.MaximumBuyCount = 3
}
