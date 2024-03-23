package filebasedcache

type Node struct {
	Value LedgerEntry
	Next  *Node
}

type List struct {
	head     *Node
	tail     *Node
	curr     *Node
	scanning bool
	size     int
	capacity int
}

func NewList(capacity int) *List {
	return &List{
		capacity: capacity,
	}
}

func (l *List) Setup(entries []LedgerEntry) {
	for _, entry := range entries {
		node := &Node{
			Value: entry,
		}
		if l.head == nil {
			l.head = node
		} else {
			l.tail.Next = node
		}
		l.tail = node
		l.size++
	}
}

func (l *List) Size() int {
	return l.size
}

func (l *List) Capacity() int {
	return l.capacity
}

func (l *List) Empty() bool {
	return l.size == 0
}

func (l *List) Full() bool {
	return l.size >= l.capacity
}

func (l *List) Pop() (LedgerEntry, bool) {
	if l.Empty() {
		return LedgerEntry{}, false
	}
	node := l.head
	l.head = l.head.Next
	l.size--
	return node.Value, true
}

func (l *List) Push(entry LedgerEntry) bool {
	node := &Node{
		Value: entry,
	}
	if l.Empty() {
		l.head = node
		l.tail = node
	} else {
		l.tail.Next = node
		l.tail = l.tail.Next
	}
	l.size++
	return true
}

func (l *List) Remove(entry LedgerEntry) {
	if l.head == nil {
		return
	}
	var removed int
	initial := l.head
	node := initial

	for node.Next != nil {
		if node.Next.Value.Filename == entry.Filename {
			removed++
			node.Next = node.Next.Next
		} else {
			node = node.Next
		}
	}
	if initial.Value.Filename == entry.Filename {
		removed++
		l.head = initial.Next
	} else {
		l.head = initial
	}
	l.size -= removed
}

func (l *List) Reset() bool {
	if l.head == nil {
		return false
	}
	l.curr = nil
	l.scanning = false
	return true
}

func (l *List) Scan() bool {
	if l.curr == nil && l.scanning {
		return false
	}
	l.scanning = true
	if l.curr == nil {
		l.curr = l.head
	} else {
		l.curr = l.curr.Next
	}

	return l.curr != nil
}

func (l *List) Entry() LedgerEntry {
	return l.curr.Value
}
