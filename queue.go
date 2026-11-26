package lightdatastructures

type Queue[T any] interface {
	Enqueue(T)
	Dequeue() T
	Peek() T
	Len() int
}
