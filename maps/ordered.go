package maps

import "github.com/CzarJoti/lightdatastructures/list"

type Ordered[K comparable, V any] struct {
	list *list.DoublyLinked[K]
	keys map[K]*list.DoublyLinkedNode[K]
	m    map[K]V
}

func NewOrdered[K comparable, V any]() *Ordered[K, V] {
	o := &Ordered[K, V]{
		list: list.New[K](),
		keys: make(map[K]*list.DoublyLinkedNode[K]),
		m:    make(map[K]V),
	}

	return o
}

func (o *Ordered[K, V]) Get(k K) V {
	return o.m[k]
}

func (o *Ordered[K, V]) GetOk(k K) (v V, ok bool) {
	v, ok = o.m[k]
	return
}

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
