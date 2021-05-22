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
	head                      *node
	payload                   chan processors.PayloadTrade
	Stop                      chan struct{}
	spoolTo                   []io.Writer
	retentionSeconds          int
	locationOffsetMiliseconds int64
}

func NewLinkedList(retentionSeconds int, stop chan struct{}, spoolTo ...io.Writer) *LinkedList {
	return &LinkedList{
		head:             &node{},
		payload:          make(chan processors.PayloadTrade),
		Stop:             stop,
		retentionSeconds: retentionSeconds,
		spoolTo:          spoolTo,
	}
}

// Listen Method listens for events coming with different offsets.
// For GMT: 3 * 3600 * 1000 = 10800000
func (l *LinkedList) Listen(locationOffsetMiliseconds int64) {
	l.locationOffsetMiliseconds = locationOffsetMiliseconds

loop:
	for {
		select {
		case <-l.Stop:
			{
				log.Println("stopping processor time list")
				break loop
			}

		case payload := <-l.payload:
			{
				l.prepend(&node{
					PayloadTrade: payload,
				})
			}
		}
	}
}

func (l *LinkedList) Payload() processors.Feed {
	return l.payload
}

// SendBufferTo Method would send current boofer to the writer.
func (l *LinkedList) SendBufferTo(w io.Writer) {
	currentNode := l.head

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

func (l *LinkedList) Terminate() {
	l.Stop <- struct{}{}
}

func (l *LinkedList) prepend(n *node) {
	n.nextNode = l.head
	l.head = n

	dropAfterTimeMiliseconds := time.Now().Unix()*1000 - l.locationOffsetMiliseconds - int64(l.retentionSeconds*1000)

	l.walkList(dropAfterTimeMiliseconds) // synchronous to preserve data
}

// walkList Internal method drops data past specified point in time.
func (l *LinkedList) walkList(dropPast int64) {
	currentNode := l.head

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

	if len(l.spoolTo) == 0 {
		return
	}

	for _, writer := range l.spoolTo {
		go writer.Write([]byte(fmt.Sprintf("%.3f --- %.f \n", sum/length, length)))
	}
}
