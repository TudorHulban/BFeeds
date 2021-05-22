package timelist

import (
	"fmt"
	"io"
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
	spoolTo         io.Writer
}

func NewLinkedList(seconds int, spoolTo io.Writer, payload chan Payload, stop chan struct{}) *LinkedList {
	return &LinkedList{
		Head:            &node{},
		Trades:          payload,
		Stop:            stop,
		TimeSpanSeconds: seconds,
		spoolTo:         spoolTo,
	}
}

// Listen Method listens for events coming with different offsets.
// For GMT: 3 * 3600 * 1000 = 10800000
func (l *LinkedList) Listen(locationOffsetMiliseconds int64) {
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
				l.prepend(&node{
					Payload: payload,
				}, locationOffsetMiliseconds)
			}
		}
	}
}

func (l *LinkedList) prepend(n *node, locationOffsetMiliseconds int64) {
	n.nextNode = l.Head
	l.Head = n

	dropAfterTimeMiliseconds := time.Now().Unix()*1000 - locationOffsetMiliseconds - int64(l.TimeSpanSeconds*1000)

	l.walkList(dropAfterTimeMiliseconds)
}

// walkList Internal method drops data past specified point in time.
func (l *LinkedList) walkList(dropPast int64) {
	currentNode := l.Head

	var length float64
	var sum float64

	for currentNode.nextNode != nil {
		if dropPast > currentNode.UNIXTimeMiliseconds {
			currentNode.nextNode = nil
			break
		}

		sum = sum + currentNode.Payload.Quantity
		length++

		currentNode = currentNode.nextNode // advance in list

	}

	// go fmt.Printf("%d ---- %.3f --- %.f \n", time.Now().Unix()*1000, sum/length, length)
	// go fmt.Printf("%.3f --- %.f \n", sum/length, length)

	go l.spoolTo.Write([]byte(fmt.Sprintf("%.3f --- %.f \n", sum/length, length)))
}

// printData Method only for testing.
func (l *LinkedList) printData() {
	currentNode := l.Head

	var length int

	fmt.Print("Printing List\n")

	for currentNode.nextNode != nil {
		fmt.Println(currentNode.Payload)

		length++
		currentNode = currentNode.nextNode
	}

	fmt.Print("\n")
	fmt.Print("Length: ", length)
	fmt.Print("\n")
}
