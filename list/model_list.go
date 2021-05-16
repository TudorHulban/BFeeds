package timelist

import (
	"fmt"
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

	dropAfterTimeMiliseconds := time.Now().Unix()*1000 - l.timeOffset - int64(l.TimeSpanSeconds*1000)

	l.walkList(dropAfterTimeMiliseconds)
}

func (l *LinkedList) walkList(dropAfter int64) {
	currentNode := l.Head

	var length float64
	var sum float64

	for currentNode.nextNode != nil {
		if dropAfter > currentNode.Payload.UNIXTimeMiliseconds {
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

func (l *LinkedList) printData() {
	next := l.Head

	for next != nil {
		fmt.Println(next.Payload)
		next = next.nextNode
	}

	fmt.Print("\n")
}

func (l *LinkedList) cleanUp() {
	close(l.Trades)
	close(l.Stop)

	l.printData()
}
