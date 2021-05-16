package conversion

import (
	"fmt"
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

				go fmt.Println(result)
			}
		}
	}
}
