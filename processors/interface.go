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

type IProcessor interface {
	Listen(locationOffsetMiliseconds int64)
	Payload() Feed // provides channel to receive trade payload
	SendBufferTo(io.Writer)
}
