package conversion

import (
	"bnb/list"
	"io"
	"log"

	"github.com/tidwall/gjson"
)

type Trade struct {
	Feed             chan []byte
	Stop             chan struct{}
	retentionSeconds int
	spoolTo          io.Writer
}

func NewTrade(feed chan []byte, stop chan struct{}, retentionSeconds int, spoolTo io.Writer) *Trade {
	return &Trade{
		Feed:             feed,
		Stop:             stop,
		retentionSeconds: retentionSeconds,
		spoolTo:          spoolTo,
	}
}

func (t *Trade) Convert() {
	payload := make(chan timelist.Payload)
	defer close(payload)

	stopList := make(chan struct{})
	defer close(stopList)

	list := timelist.NewLinkedList(t.retentionSeconds, t.spoolTo, payload, stopList)
	defer func() {
		list.Stop <- struct{}{}
	}()

	go list.Listen(0)

loop:
	for {
		select {
		case <-t.Stop:
			{
				log.Println("stopping converter")
				break loop
			}

		case payload := <-t.Feed:
			{
				result := gjson.GetManyBytes(payload, "s", "T", "q", "p")

				list.Trades <- timelist.Payload{
					Symbol:              result[0].String(),
					UNIXTimeMiliseconds: result[1].Int(),
					Price:               result[2].Float(),
					Quantity:            result[3].Float(),
				}
			}
		}
	}
}
