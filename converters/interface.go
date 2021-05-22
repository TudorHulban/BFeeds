package converters

type Feed chan []byte

type IConverter interface {
	Convert()
	Payload() Feed // provides feed channel to receive Binance payload
	Terminate()
}
