package converters

type Feed chan []byte

type IConverter interface {
	Convert(locationOffsetMiliseconds int64)
	Payload() Feed // provides feed channel to receive Binance payload
	Terminate()
}
