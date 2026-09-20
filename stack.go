package main

type Stack struct {
	top  *Node
	size int
}

func (s *Stack) Push(inf string) {
	newNode := &Node{Inf: inf}
	if s.top == nil {
		s.top = newNode
	} else {
		newNode.Next = s.top
		s.top = newNode
	}
	s.size++
}

func (s *Stack) Pop() string {
	if s.top == nil {
		return ""
	}
	inf := s.top.Inf
	s.top = s.top.Next
	s.size--
	return inf
}

func (s *Stack) Size() int {
	return s.size
}

func (s *Stack) Review() {
	current := s.top
	for current != nil {
		println("Inf:", current.Inf)
		current = current.Next
	}
}

func (s *Stack) Peek() (*Node, bool) {
	return s.top, s.top != nil
}
