package maps

import "github.com/CzarJoti/lightdatastructures/list"

// Ordered represents an ordered map
// Internally uses a linked list to mangage the order of insertions
type Ordered[K comparable, V any] struct {
	list *list.DoublyLinked[K]
	keys map[K]*list.DoublyLinkedNode[K]
	m    map[K]V
}

// NewOrdered creates an empty Ordered
func NewOrdered[K comparable, V any]() *Ordered[K, V] {
	o := &Ordered[K, V]{
		list: list.New[K](),
		keys: make(map[K]*list.DoublyLinkedNode[K]),
		m:    make(map[K]V),
	}

	return o
}

// Get returns the value of k for ordered
func (o *Ordered[K, V]) Get(k K) V {
	return o.m[k]
}

// GetOk returns the value of k for ordered
// ok indicates if the value exists, it is true if it does or false otherwise
func (o *Ordered[K, V]) GetOk(k K) (v V, ok bool) {
	v, ok = o.m[k]
	return
}

// Set sets the value of k to v
// It updates the insertion order of the list so that
// k is recorded as the latest updated key
func (o *Ordered[K, V]) Set(k K, v V) {
	_, contains := o.m[k]
	o.m[k] = v
	if !contains {
		n := o.list.Prepend(k)
		o.keys[k] = n
	} else {
		n := o.keys[k]
		o.list.MoveToFront(n)
	}
}
