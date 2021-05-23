package rolling

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

type RollingList struct {
	head                      *node
	streamData                processors.StreamData
	stop                      chan struct{}
	spoolTo                   []io.Writer
	locationOffsetMiliseconds int64
	retentionSeconds          int
}

var _ processors.IProcessor = (*RollingList)(nil)

func NewLinkedList(symbol string, retentionSeconds int, spoolTo ...io.Writer) *RollingList {
	return &RollingList{
		head: &node{},
		streamData: processors.StreamData{
			Stream: symbol,
			Feed:   make(chan processors.PayloadTrade),
		},
		stop:             make(chan struct{}),
		retentionSeconds: retentionSeconds,
		spoolTo:          spoolTo,
	}
}

// Listen Method listens for events coming with different offsets.
// For GMT: 3 * 3600 * 1000 = 10800000
func (l *RollingList) Listen(locationOffsetMiliseconds int64) {
	l.locationOffsetMiliseconds = locationOffsetMiliseconds

loop:
	for {
		select {
		case <-l.stop:
			{
				log.Println("stopping processor time list")
				break loop
			}

		case payload := <-l.streamData.Feed:
			{
				l.prepend(&node{
					PayloadTrade: payload,
				})
			}
		}
	}
}

func (l *RollingList) Payload() processors.StreamData {
	return l.streamData
}

// SendBufferTo Method would send current boofer to the writer.
func (l *RollingList) SendBufferTo(w io.Writer) {
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

func (l *RollingList) Terminate() {
	defer l.cleanUp()

	l.stop <- struct{}{}
}

func (l *RollingList) prepend(n *node) {
	n.nextNode = l.head
	l.head = n

	dropAfterTimeMiliseconds := time.Now().Unix()*1000 - l.locationOffsetMiliseconds - int64(l.retentionSeconds*1000)

	l.walkList(dropAfterTimeMiliseconds) // synchronous to preserve data
}

// walkList Internal method drops data past specified point in time.
func (l *RollingList) walkList(dropPast int64) {
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
		go writer.Write([]byte(fmt.Sprintf("%s: %.3f --- %.f \n", l.streamData.Stream, sum/length, length)))
	}
}

func (l *RollingList) cleanUp() {
	close(l.stop)
	close(l.streamData.Feed)
}
