package processors

import (
	"io"
)

type PayloadTrade struct {
	Symbol              string
	UNIXTimeMiliseconds int64
	Price               float64
	Quantity            float64
}

type Feed chan PayloadTrade

type StreamData struct {
	Stream string
	Feed   Feed
}

type IProcessor interface {
	Listen(locationOffsetMiliseconds int64)

	// provides symbol and channel to receive trade payload
	Payload() StreamData
	SendBufferTo(io.Writer)
	Terminate()
}
