package list

import (
	"fmt"
)

// DoublyLinked represents a doublely linked list
// It can be traversed either forwards or backwards
// When computing the index, negative values are calculated backwards
// With -1 representing the last value of the list
type DoublyLinked[T any] struct {
	len        int
	head, tail *doublyLinkedNode[T]
}

// NewDoublyLinked creates and returns a DoublyLinked list with an initial length set by initLen
// All values in this list are initialized to their zero value
func NewDoublyLinked[T any](initLen int) *DoublyLinked[T] {
	if initLen < 0 {
		panic(fmt.Sprintf("invalid list length %d", initLen))
	}

	list := &DoublyLinked[T]{
		len: initLen,
	}

	for range initLen {
		node := &doublyLinkedNode[T]{
			prev: list.tail,
		}
		if list.tail == nil {
			list.head = node
			list.tail = node
		} else {
			list.tail.next = node
			list.tail = node
		}

	}

	return list
}

type doublyLinkedNode[T any] struct {
	prev, next *doublyLinkedNode[T]
	value      T
}

// Len returns the length of the list
func (l *DoublyLinked[T]) Len() int {
	return l.len
}

// Get returns the value of the list at i
func (l *DoublyLinked[T]) Get(i int) T {
	l.checkBounds(i)

	return l.nodeAt(i).value
}

// Peek returns the value of the last value of the list
func (l *DoublyLinked[T]) Peek() T {
	return l.Get(-1)
}

// Set sets the value of the list at i to t
// returns the original value of the list at i
func (l *DoublyLinked[T]) Set(i int, t T) T {
	l.checkBounds(i)
	n := l.nodeAt(i)
	v := n.value
	n.value = t
	return v
}

// Append appends t to the end of the list
func (l *DoublyLinked[T]) Append(t T) {
	l.Add(t, -1)
}

// Prepend prepends t to the beginning of the list
func (l *DoublyLinked[T]) Prepend(t T) {
	l.Add(t, 0)
}

// Push is a wrapper for Append
func (l *DoublyLinked[T]) Push(t T) {
	l.Append(t)
}

// Enqueue is a wrapper for Prepend
func (l *DoublyLinked[T]) Enqueue(t T) {
	l.Prepend(t)
}

// Swap swaps the values at k and j in the list
func (l *DoublyLinked[T]) Swap(k, j int) {
	l.checkBounds(k)
	l.checkBounds(j)

	if k == j {
		return
	}

	kN := l.nodeAt(k)
	jN := l.nodeAt(j)
	kN.value, jN.value = jN.value, kN.value
}

// Add inserts the value t before the value at i if i is nonnegative
// if i is negative, Add inserts the value t after the value i
// if i == l.Len(), t is inserted at the end of the list
// if i == -(l.Len() + 1), t is inserted at the beginning of the list
func (l *DoublyLinked[T]) Add(t T, i int) {
	if l.len == 0 && (i == 0 || i == -1) {
		n := &doublyLinkedNode[T]{value: t}
		l.head = n
		l.tail = n
		l.len = 1
		return
	}

	if i == l.len || i == -1 {
		n := &doublyLinkedNode[T]{
			value: t,
			prev:  l.tail,
		}

		l.tail.next = n
		l.tail = n
		l.len++
		return
	}

	if i == 0 || i == (-(l.len + 1)) {
		n := &doublyLinkedNode[T]{
			value: t,
			next:  l.head,
		}

		l.head.prev = n
		l.head = n
		l.len++
		return
	}

	l.checkBounds(i)

	at := l.nodeAt(i)

	var next, prev *doublyLinkedNode[T]
	if i >= 0 {
		next = at
		prev = at.prev
	} else {
		prev = at
		next = at.next
	}

	n := &doublyLinkedNode[T]{
		value: t,
		prev:  prev,
		next:  next,
	}

	prev.next = n
	next.prev = n

	l.len++
}

// Remove removes the value at i
// joining the list back together such that L[i] = L[i  + 1]
// returns the original value at i
func (l *DoublyLinked[T]) Remove(i int) T {
	l.checkBounds(i)
	node := l.nodeAt(i)
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}

	val := node.value
	var zero T
	node.value = zero
	node.prev = nil
	node.next = nil

	l.len -= 1
	return val
}

// Pop removes the final element of the list
func (l *DoublyLinked[T]) Pop() T {
	return l.Remove(-1)
}

// Dequeue is a wrapper for Pop
func (l *DoublyLinked[T]) Dequeue() T {
	return l.Pop()
}

// ForEach iterates through the list starting at the head, calling f on each value
func (l *DoublyLinked[T]) ForEach(f func(i int, t T)) {
	node := l.head
	for i := 0; i < l.len; i, node = i+1, node.next {
		f(i, node.value)
	}
}

// ReverseForEach iterates through the list in reverse order, starting at the tail
// calling f on each value
func (l *DoublyLinked[T]) ReverseForEach(f func(i int, t T)) {
	node := l.tail
	for i := -1; i >= -l.len; i, node = i-1, node.prev {
		f(i, node.value)
	}
}

func (l *DoublyLinked[T]) nodeAt(i int) *doublyLinkedNode[T] {
	if i < 0 {
		i += l.len
	}

	if i <= l.len/2 {
		node := l.head
		for range i {
			node = node.next
		}

		return node
	}

	node := l.tail
	for range l.len - 1 - i {
		node = node.prev
	}

	return node
}

func (l *DoublyLinked[T]) checkBounds(i int) {
	totalLen := l.len
	if i >= totalLen || i < -totalLen {
		panic(fmt.Sprintf("list index %d out of bounds", i))
	}
}
