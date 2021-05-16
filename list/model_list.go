package timelist

import (
	"log"
	"time"
)

type Payload struct {
	Symbol              string
	UNIXTimeMiliseconds int64
	Price               float64
	Quantity            float64
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
	timeOffset      int64
}

func NewLinkedList(seconds int, payload chan Payload, stop chan struct{}) *LinkedList {
	return &LinkedList{
		Head:            &node{},
		Trades:          payload,
		Stop:            stop,
		TimeSpanSeconds: seconds,
		timeOffset:      3 * 3600 * 1000, // offset to GMT in miliseconds
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
	currentNode := l.Head

	length := 1.
	sum := currentNode.Payload.Quantity

	for currentNode.nextNode != nil {
		if currentNode.Payload.UNIXTimeMiliseconds < (time.Now().Unix()*1000 - int64(l.TimeSpanSeconds*60000) - l.timeOffset) {
			log.Println("found aged node")

			currentNode.nextNode = nil
			break
		}

		sum = sum + currentNode.Payload.Quantity
		length++

		currentNode = currentNode.nextNode // advance in list
	}

	log.Printf("%d ---- %.3f --- %.f \n", time.Now().Unix()*1000, sum/length, length)
}
