package timelist

import (
	"fmt"
	"log"
)

type Payload struct {
	Symbol     string
	UNIXTimeMs int64 // miliseconds
	Price      float64
	Quantity   float64
}

type Node struct {
	Payload

	nextNode *Node
}

type LinkedList struct {
	Head            *Node
	Trades          chan Payload
	Stop            chan struct{}
	TimeSpanSeconds int

	length          int
	averageQuantity float64
}

func NewLinkedList(seconds int, payload chan Payload, stop chan struct{}) *LinkedList {
	return &LinkedList{
		TimeSpanSeconds: seconds,
		Trades:          payload,
		Stop:            stop,
	}
}

func (l *LinkedList) Listen() {
loop:
	for {
		select {
		case <-l.Stop:
			{
				log.Println("stopping list")
				break loop
			}

		case payload := <-l.Trades:
			{
				go fmt.Println(payload)
			}
		}
	}
}

func (l *LinkedList) prepend(n *Node) {
	n.nextNode = l.Head
	l.Head = n
}

func (l *LinkedList) walkList() {
	next := l.Head

	for next != nil {
		fmt.Printf("%d ", next.Payload)
		next = next.nextNode
	}
}
