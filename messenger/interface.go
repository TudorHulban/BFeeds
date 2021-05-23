package messengers

type IMessenger interface {
	Subscribe(symbols ...string) error
	Unsubsribe(symbols ...string) error
	Terminate()
}
