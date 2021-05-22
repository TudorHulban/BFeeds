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

type IProcessor interface {
	Listen(locationOffsetMiliseconds int64)
	Payload() chan PayloadTrade // provides channel to receive trade payload
	SendBufferTo(io.Writer)
}
