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

	defer func() {
		stopProcessor <- struct{}{}
		close(stopProcessor)
	}()

	// creation of a trade converter
	stopConverter := make(chan struct{})
	conv := trade.NewTradeConverter(processorTimeList, stopConverter)

	defer func() {
		stopConverter <- struct{}{}
		close(stopConverter)
	}()

	// creation of a exchange
	client, errNew := exchange.NewExchange(exchange.Config{
		URI: urlBinance,
	})
	if errNew != nil {
		fmt.Println(errNew)
		os.Exit(1)
	}
	defer close(client.Stop)

	go client.ReadMessages(conv)

	<-client.Stop
}
