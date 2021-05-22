package trade

import (
	"bnb/converters"
	"bnb/processors"

	"log"

	"github.com/tidwall/gjson"
)

type Trade struct {
	payload chan []byte
	Stop    chan struct{}
}

func NewTradeConverter(stop chan struct{}) *Trade {
	return &Trade{
		payload: make(chan []byte),
		Stop:    stop,
	}
}

// Convert Method converts Binance messages and pushes them further to a processor.
func (t *Trade) Convert(proc processors.IProcessor) {
	processorPayload := proc.Payload()
	defer close(processorPayload)

	// the payload channel should be defined prior to start listening
	go proc.Listen(0)

loop:
	for {
		select {
		case <-t.Stop:
			{
				log.Println("stopping converter")
				break loop
			}

		case payload := <-t.payload:
			{
				result := gjson.GetManyBytes(payload, "s", "T", "q", "p")

				processorPayload <- processors.PayloadTrade{
					Symbol:              result[0].String(),
					UNIXTimeMiliseconds: result[1].Int(),
					Price:               result[2].Float(),
					Quantity:            result[3].Float(),
				}
			}
		}
	}
}

func (t *Trade) Payload() converters.Feed {
	return t.payload
}
