package parser

// TokenStack is a stack of any kind of Token implemented as
// a linked list. This can be used for parenthesis, order of ops, or brackets
type TokenStack struct {
	top    *node
	length int
}

type node struct {
	value *Token
	prev  *node
}

// NewStack - Creates a new, empty, token stack
func NewStack() *TokenStack {
	return &TokenStack{nil, 0}
}

// Len - returns the number of tokens in the stack
func (s *TokenStack) Len() int {
	return s.length
}

// Peek - View the top item on the stack
func (s *TokenStack) Peek() *Token {
	if s.length == 0 {
		return nil
	}
	return s.top.value
}

// Pop the top item of the stack and return it
func (s *TokenStack) Pop() *Token {
	if s.length == 0 {
		return nil
	}
	n := s.top
	s.top = n.prev
	s.length--
	return n.value
}

// Push a value onto the top of the TokenStack
func (s *TokenStack) Push(t *Token) {
	n := &node{t, s.top}
	s.top = n
	s.length++
}
