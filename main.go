package main

import (
	"bnb/converters/trade"
	"bnb/exchange"
	"bnb/processors/timelist"
	"fmt"
	"os"
)

const urlBinance = "wss://stream.binance.com:9443/ws/bnbusdt@trade"

func main() {
	// creation of a processor
	stopProcessor := make(chan struct{})
	processorTimeList := timelist.NewLinkedList(1, stopProcessor, os.Stdout)

	// creation of a trade converter
	stopConverter := make(chan struct{})
	conv := trade.NewTradeConverter(processorTimeList, stopConverter)

	// creation of a exchange
	exch, errNew := exchange.NewExchange(exchange.Config{
		URI: urlBinance,
	})
	if errNew != nil {
		fmt.Println(errNew)
		os.Exit(1)
	}
	defer close(exch.Stop)

	go exch.ReadMessages(conv)

	<-exch.Stop

	conv.Terminate()
	close(stopConverter)
}
