package exchange

import (
	"bnb/converters"
	"fmt"
	"io"
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

// Client Concentrates websocket information.
type Client struct {
	connection *websocket.Conn
	URL        url.URL

	Stop      chan struct{}
	interrupt chan os.Signal

	retentionSeconds int
	spoolTo          io.Writer
}

func NewExchange(cfg Config, retentionSeconds int, spoolTo io.Writer) (*Client, error) {
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

	return &Client{
		connection:       conn,
		URL:              *url,
		Stop:             make(chan struct{}),
		interrupt:        interrupt,
		retentionSeconds: retentionSeconds,
		spoolTo:          spoolTo,
	}, nil
}

// ReadMessages Method reads websocket feed and pushes it to a converter payload channel.
func (c *Client) ReadMessages(conv converters.IConverter) {
	// payload := make(chan []byte)
	// defer close(payload)

	// stopConversion := make(chan struct{})
	// defer close(stopConversion)

	// converter := conversion.NewTrade(payload, stopConversion, c.retentionSeconds, c.spoolTo)
	// defer func() {
	// 	converter.Stop <- struct{}{}
	// 	c.Stop <- struct{}{}
	// }()

	converterPayload := conv.Payload()
	defer close(converterPayload)

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

func (c *Client) CleanUp() {
	close(c.Stop)
	close(c.interrupt)
	fmt.Println("cleaned up exchange client")
}