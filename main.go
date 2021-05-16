package main

import (
	"bnb/exchange"
	"fmt"
	"os"
)

const urlBinance = "wss://stream.binance.com:9443/ws/bnbusdt@trade"

func main() {
	cfg := exchange.Config{
		URI: urlBinance,
	}

	client, errNew := exchange.NewClient(cfg, 300)
	if errNew != nil {
		fmt.Println(errNew)
		os.Exit(1)
	}

	go client.ReadMessages()

	<-client.Stop
	client.CleanUp()
}
