package conversion

import (
	"fmt"
	"log"
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
	// TODO: add JSON parser

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
				// parsedValues, errParse := p.ParseBytes(payload)
				// if errParse != nil {
				// 	fmt.Println("conversion glitch:", errParse)
				// 	return
				// }

				go fmt.Println(string(payload))
			}
		}
	}
}
