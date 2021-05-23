package converters

type Feed chan []byte

type Streams struct {
	Symbols []string
	Feed    Feed
}

type IConverter interface {
	Convert(locationOffsetMiliseconds int64)
	Payload() Streams // provides feed channel to receive exchange payload
	Terminate()
}
