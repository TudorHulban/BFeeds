package main

import (
	"bnb/converters/streams"
	"bnb/exchange"
	"bnb/processors"
	"bnb/processors/rolling"
	"fmt"
	"os"
	// "github.com/pkg/profile"
)

const urlBinance = "wss://stream.binance.com:9443/ws/bnbusdt@trade"

func main() {
	// defer profile.Start(profile.MemProfile).Stop()

	// creation of a processor
	processor := rolling.NewLinkedList(1, os.Stdout)

	// creation of a trade converter
	// converter := trade.NewTradeConverter(processorTimeList)

	converter := streams.NewStreamsConverter([]processors.IProcessor{processor}...)

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
