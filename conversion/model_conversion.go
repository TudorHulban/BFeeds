package conversion

import (
	"bnb/list"
	"log"

	"github.com/tidwall/gjson"
)

type Trade struct {
	Feed chan []byte
	Stop chan struct{}
}

func NewTrade(feed chan []byte, stop chan struct{}) *Trade {
	return &Trade{
		Feed: feed,
		Stop: stop,
	}
}

func (t *Trade) Convert() {
	payload := make(chan timelist.Payload)
	stop := make(chan struct{})

	list := timelist.NewLinkedList(1, payload, stop)
	go list.Listen()

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

	list.Stop <- struct{}{}
	close(payload)
	close(stop)
}
