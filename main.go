package main

import (
	"fmt"

	"github.com/ksyoon0321/gotrade/cmd"
	"github.com/ksyoon0321/gotrade/coin"
	"github.com/ksyoon0321/gotrade/config"
	"github.com/ksyoon0321/gotrade/strategy"
	"github.com/ksyoon0321/gotrade/txlog"
)

func main() {
	fmt.Println("init main")

	trans := txlog.NewConsoleTransfer()
	txmgr := txlog.NewTxManager(trans)

	director := cmd.NewMarketDirector(txmgr)

	director.AddMarketConfig(config.CfgUpbit)

	director.RegistStrategy(strategy.NewStrategyTrends())
	director.RegistStrategy(strategy.NewStrategyAvg(coin.MIN1))
	director.RegistStrategy(strategy.NewStrategyAvg(coin.MIN30))
	director.RegistStrategy(strategy.NewStrategyAvg(coin.MIN60))
	director.RegistStrategy(strategy.NewStrategyAvg(coin.MIN240))
	director.RegistStrategy(strategy.NewStrategyJump())
	director.RegistStrategy(strategy.NewStrategyCross60())
	director.RegistStrategy(strategy.NewStrategyCrossHigh())

	director.Run()
}
