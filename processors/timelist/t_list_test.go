package timelist

import (
	"testing"
	"time"
)

func TestList(t *testing.T) {
	payload := make(chan Payload)
	stop := make(chan struct{})

	list := NewLinkedList(1, payload, stop)
	go list.Listen(0)

	payload <- Payload{
		Quantity:            1.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	time.Sleep(300 * time.Millisecond)

	payload <- Payload{
		Quantity:            2.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	time.Sleep(300 * time.Millisecond)

	payload <- Payload{
		Quantity:            3.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	time.Sleep(1400 * time.Millisecond)

	payload <- Payload{
		Quantity:            4.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	payload <- Payload{
		Quantity:            5.,
		UNIXTimeMiliseconds: time.Now().Unix() * 1000,
	}

	stop <- struct{}{}

	list.CleanUp()
}
