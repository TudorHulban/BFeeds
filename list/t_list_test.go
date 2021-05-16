package timelist

import (
	"testing"
	"time"
)

func TestList(t *testing.T) {
	payload := make(chan Payload)
	stop := make(chan struct{})

	list := NewLinkedList(1, payload, stop)
	go list.Listen()

	payload <- Payload{
		Quantity:            1.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	// time.Sleep(100 * time.Millisecond)

	time.Sleep(1100 * time.Millisecond)

	payload <- Payload{
		Quantity:            2.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	payload <- Payload{
		Quantity:            3.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	stop <- struct{}{}
}
