package main

type Queue struct {
	head *Node
	last *Node
	size int
}

func (q *Queue) Enqueue(inf string) {
	newNode := &Node{Inf: inf}
	if q.head == nil {
		q.head = newNode
		q.last = newNode
	} else {
		q.last.Next = newNode
		q.last = newNode
	}
	q.size++
}

func (q *Queue) Dequeue() string {
	if q.head == nil {
		return ""
	}
	inf := q.head.Inf
	q.head = q.head.Next
	q.size--

	if q.head == nil {
		q.last = nil
	}

	return inf
}

func (q *Queue) Size() int {
	return q.size
}

func (q *Queue) Review() {
	current := q.head
	for current != nil {
		println("Inf:", current.Inf)
		current = current.Next
	}
}
