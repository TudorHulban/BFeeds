package messbinance

import (
	"io"

	"github.com/gorilla/websocket"
)

type MessengerBinance struct {
	connection    *websocket.Conn
	feed          chan string
	stop          chan struct{}
	writeErrorsTo io.Writer
}

func NewMessengerBinance(conn *websocket.Conn, writeErrorsTo io.Writer) *MessengerBinance {
	return &MessengerBinance{
		connection: conn,
		feed:       make(chan string),
	}
}

func (m *MessengerBinance) Terminate() {
	defer m.cleanUp()

	m.stop <- struct{}{}
}

func (m *MessengerBinance) sendMessage() {
	for {
		select {
		case <-m.stop:
			{
				return
			}

		case msg := <-m.feed:
			{
				err := m.connection.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					m.writeErrorsTo.Write([]byte(err.Error()))
					return
				}
			}
		}
	}
}

func (m *MessengerBinance) cleanUp() {
	close(m.feed)
	close(m.stop)
}
