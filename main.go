package main

import (
	"bnb/converters/trade"
	"bnb/exchange"
	"bnb/processors/rolling"
	"fmt"
	"os"
)

const urlBinance = "wss://stream.binance.com:9443/ws/bnbusdt@trade"

func main() {
	// creation of a processor
	processorTimeList := rolling.NewLinkedList(1, os.Stdout)

	// creation of a trade converter
	converter := trade.NewTradeConverter(processorTimeList)

	// creation of a exchange
	exch, errNew := exchange.NewExchange(exchange.Config{
		URI: urlBinance,
	})
	if errNew != nil {
		fmt.Println(errNew)
		os.Exit(1)
	}

	go exch.ReadMessages(converter)
	exch.Work()

	converter.Terminate()
}
