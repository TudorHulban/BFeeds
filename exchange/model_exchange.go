package exchange

import (
	"bnb/converters"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type Config struct {
	RequestHeader       http.Header
	URI                 string
	PongIntervalSeconds uint
}

// Exchange Concentrates websocket connection information.
type Exchange struct {
	connection *websocket.Conn
	URL        url.URL

	Stop      chan struct{}
	interrupt chan os.Signal
}

func NewExchange(cfg Config) (*Exchange, error) {
	url, errParse := url.Parse(cfg.URI)
	if errParse != nil {
		return nil, errParse
	}

	conn, _, errConn := websocket.DefaultDialer.Dial(url.String(), cfg.RequestHeader)
	if errConn != nil {
		return nil, errConn
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	return &Exchange{
		connection: conn,
		URL:        *url,
		Stop:       make(chan struct{}),
		interrupt:  interrupt,
	}, nil
}

// ReadMessages Method reads websocket feed and pushes it to a converter payload channel.
func (c *Exchange) ReadMessages(conv converters.IConverter) {
	converterPayload := conv.Payload()
	defer close(converterPayload)

	defer close(c.interrupt)

	go conv.Convert()

loop:
	for {
		select {
		case <-c.interrupt:
			{
				fmt.Println("interrupt")
				break loop
			}
		default:
			{
				_, message, errRead := c.connection.ReadMessage()
				if errRead != nil {
					fmt.Println("read glitch:", errRead)
					return
				}

				converterPayload <- message
			}
		}
	}
}
