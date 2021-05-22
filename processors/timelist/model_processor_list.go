package timelist

import (
	"bnb/processors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

type node struct {
	processors.PayloadTrade

	nextNode *node
}

type LinkedList struct {
	Head            *node
	Payload         chan processors.PayloadTrade
	Stop            chan struct{}
	TimeSpanSeconds int
	spoolTo         io.Writer
}

func NewLinkedList(seconds int, spoolTo io.Writer, payload chan processors.PayloadTrade, stop chan struct{}) *LinkedList {
	return &LinkedList{
		Head:            &node{},
		Payload:         payload,
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

		case payload := <-l.Payload:
			{
				l.prepend(&node{
					PayloadTrade: payload,
				}, locationOffsetMiliseconds)
			}
		}
	}
}

// SendBufferTo Method would send current boofer to the writer.
func (l *LinkedList) SendBufferTo(w io.Writer) {
	currentNode := l.Head

	var length int

	w.Write([]byte("Printing Buffer\n"))

	for currentNode.nextNode != nil {
		w.Write([]byte(fmt.Sprintf("%v", currentNode.PayloadTrade)))

		length++
		currentNode = currentNode.nextNode
	}

	w.Write([]byte("\n"))
	w.Write([]byte("Length: " + strconv.Itoa(length)))
	w.Write([]byte("\n"))
}

func (l *LinkedList) prepend(n *node, locationOffsetMiliseconds int64) {
	n.nextNode = l.Head
	l.Head = n

	dropAfterTimeMiliseconds := time.Now().Unix()*1000 - locationOffsetMiliseconds - int64(l.TimeSpanSeconds*1000)

	l.walkList(dropAfterTimeMiliseconds) // synchronous to preserve data
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

		sum = sum + currentNode.PayloadTrade.Quantity
		length++

		currentNode = currentNode.nextNode // advance in list

	}

	// go fmt.Printf("%d ---- %.3f --- %.f \n", time.Now().Unix()*1000, sum/length, length)
	// go fmt.Printf("%.3f --- %.f \n", sum/length, length)

	go l.spoolTo.Write([]byte(fmt.Sprintf("%.3f --- %.f \n", sum/length, length)))
}
