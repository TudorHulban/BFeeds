package main

import (
	"bnb/converters"
	"bnb/converters/trade"
	"bnb/exchange"
	"fmt"
	"os"
)

const urlBinance = "wss://stream.binance.com:9443/ws/bnbusdt@trade"

func main() {
	// creation of a trade converter
	conv := trade.NewTradeConverterr()

	cfg := exchange.Config{
		URI: urlBinance,
	}

	client, errNew := exchange.NewExchange(cfg, 3, os.Stdout)
	if errNew != nil {
		fmt.Println(errNew)
		os.Exit(1)
	}
	defer client.CleanUp()

	go client.ReadMessages()

	<-client.Stop
}
