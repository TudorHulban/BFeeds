package timelist

type Payload struct {
	UNIXTimeMs int64 // miliseconds
	Price      float64
	Quantity   float32
}

type Node struct {
	Payload

	nextNode *Node
}

type LinkedList struct {
	Head            *Node
	TimeSpanMinutes int

	length          int
	averageQuantity float64
}

func NewLinkedList(n *Node, minutes int) *LinkedList {
	return &LinkedList{
		Head:            n,
		TimeSpanMinutes: minutes,
	}
}

func (l *LinkedList) Prepend(n *Node) {
	n.nextNode = l.Head
	l.Head = n
}

func (l *LinkedList) WalkList() {
	next := l.Head

	for next != nil {
		fmt.Printf("%d ", next.Data)
		next = next.nextNode
	}
}
