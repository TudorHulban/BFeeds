package trade

import (
	"bnb/converters"
	"bnb/processors"

	"log"

	"github.com/tidwall/gjson"
)

type ConvertorTrade struct {
	processor processors.IProcessor

	payload chan []byte
	stop    chan struct{}
}

func NewTradeConverter(proc processors.IProcessor) *ConvertorTrade {
	return &ConvertorTrade{
		processor: proc,
		payload:   make(chan []byte),
		stop:      make(chan struct{}),
	}
}

// Convert Method converts Binance messages and pushes them further to a processor.
func (t *ConvertorTrade) Convert() {
	go t.processor.Listen(0)
	defer t.processor.Terminate()

	processorPayload := t.processor.Payload()

loop:
	for {
		select {
		case <-t.stop:
			{
				log.Println("stopping converter")
				break loop
			}

		case payload := <-t.payload:
			{
				if len(payload) == 0 {
					continue
				}

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

func (t *ConvertorTrade) Payload() converters.Feed {
	return t.payload
}

func (t *ConvertorTrade) Terminate() {
	defer t.cleanUp()

	t.stop <- struct{}{}
}

func (t *ConvertorTrade) cleanUp() {
	close(t.stop)
}
