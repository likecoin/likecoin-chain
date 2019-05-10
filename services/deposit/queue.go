package deposit

type queue struct {
	Array []uint64
	Head  int
	Tail  int
}

func newQueue(capacity int) *queue {
	return &queue{
		Array: make([]uint64, capacity),
	}
}

func (q *queue) size() int {
	return (q.Tail - q.Head + len(q.Array)) % len(q.Array)
}

func (q *queue) enqueue(n uint64) {
	l := len(q.Array)
	if q.size() == l-1 {
		newArray := make([]uint64, len(q.Array)*2)
		copy(newArray[:], q.Array[q.Head:l])
		copy(newArray[l-q.Head:], q.Array[:q.Head])
		q.Array = newArray
		q.Head = 0
		q.Tail = l - 1
		l = len(newArray)
	}
	q.Array[q.Tail] = n
	q.Tail = (q.Tail + 1) % len(q.Array)
}

func (q *queue) dequeue() {
	if q.size() == 0 {
		return
	}
	q.Head = (q.Head + 1) % len(q.Array)
}

func (q *queue) peek() (uint64, bool) {
	if q.size() == 0 {
		return 0, false
	}
	return q.Array[q.Head], true
}
