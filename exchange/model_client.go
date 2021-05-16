package exchange

import (
	"bnb/conversion"
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

// Client Concentrates websocket information.
type Client struct {
	connection *websocket.Conn
	URL        url.URL

	Stop      chan struct{}
	interrupt chan os.Signal

	retentionSeconds int
}

func NewClient(cfg Config, retentionSeconds int) (*Client, error) {
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
	}, nil
}

func (c *Client) ReadMessages() {
	payload := make(chan []byte)
	stopConversion := make(chan struct{})

	converter := conversion.NewTrade(payload, stopConversion, c.retentionSeconds)
	go converter.Convert()

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

				converter.Feed <- message
			}
		}
	}

	converter.Stop <- struct{}{}
	close(payload)
	close(stopConversion)

	c.Stop <- struct{}{}
}

func (c *Client) CleanUp() {
	close(c.Stop)
	close(c.interrupt)
	fmt.Println("cleaned up exchange client")
}
