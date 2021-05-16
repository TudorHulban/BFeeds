package timelist

import (
	"log"
)

type Payload struct {
	Symbol     string
	UNIXTimeMs int64 // miliseconds
	Price      float64
	Quantity   float64
}

type node struct {
	Payload

	nextNode *node
}

type LinkedList struct {
	Head            *node
	Trades          chan Payload
	Stop            chan struct{}
	TimeSpanSeconds int
}

func NewLinkedList(seconds int, payload chan Payload, stop chan struct{}) *LinkedList {
	return &LinkedList{
		Head:            &node{},
		Trades:          payload,
		Stop:            stop,
		TimeSpanSeconds: seconds,
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
				// go fmt.Println(payload)

				l.prepend(&node{
					Payload: payload,
				})
			}
		}
	}
}

func (l *LinkedList) prepend(n *node) {
	n.nextNode = l.Head
	l.Head = n

	l.walkList()
}

func (l *LinkedList) walkList() {
	next := l.Head

	length := 1.
	sum := l.Head.Payload.Quantity

	for next != nil {
		sum = sum + next.Payload.Quantity
		length++

		next = next.nextNode
	}

	log.Printf("%.3f --- %.f \n", sum/length, length)
}
