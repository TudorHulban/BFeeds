package main

import (
	"bnb/converters/streams"
	"bnb/exchange"
	"bnb/fiber"
	"bnb/processors"
	"bnb/processors/rolling"
	"fmt"
	"os"
	"strings"
	// "github.com/pkg/profile"
)

// const urlBinance = "wss://stream.binance.com:9443/ws/bnbusdt@trade"
// const urlBinance = "wss://stream.binance.com:9443/stream?streams=bnbusdt@trade/btcusdt@trade"

const rootStreamBinance = "wss://stream.binance.com:9443/stream?streams="

func main() {
	// defer profile.Start(profile.MemProfile).Stop()

	// HTTP Server
	serv := fiber.NewHTTPServer()
	go serv.Work()

	// creation of a processor
	processorBNB := rolling.NewLinkedList("bnbusdt@trade", 1, serv)
	// processorBTC := rolling.NewLinkedList("btcusdt@trade", 3, os.Stdout)

	// creation of a trade converter
	// converter := trade.NewTradeConverter(processorTimeList)

	converter := streams.NewStreamsConverter([]processors.IProcessor{processorBNB}...)

	urlStreams := rootStreamBinance + strings.Join(converter.Payload().Symbols, "/")

	// creation of a exchange
	exch, errNew := exchange.NewExchange(exchange.Config{
		URI: urlStreams,
	})
	if errNew != nil {
		fmt.Println(errNew)
		os.Exit(1)
	}

	go exch.ReadMessages(converter)
	exch.Work()

	converter.Terminate()
}
