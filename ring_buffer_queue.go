package main

type RingBufferQueue struct {
	cap, head, tail int
	elems         []int
}

func NewRingBufferQueue() *RingBufferQueue {
	return &RingBufferQueue{
		cap:     10,
		elems: make([]int, 10),
		head:  0,
		tail:  0,
	}
}

func (q *RingBufferQueue) empty() bool {
  return q.tail == q.head
}

func (q *RingBufferQueue) next(i int) int {
	return (i + 1) % q.cap
}

func (q *RingBufferQueue) prev(i int) int {
	return (i - 1 + q.cap) % q.cap
}

func (q *RingBufferQueue) grow() {
	newElems := make([]int, 2*q.cap)
	toPos := 0

	for i := q.head; i != q.tail; i = q.next(i) {
		newElems[toPos] = q.elems[i]
		toPos++
	}
  q.cap = 2 * q.cap
	q.head = 0
	q.tail = toPos
  q.elems = newElems
}

func (q *RingBufferQueue) push(a int) {
	if q.next(q.tail) == q.head {
		q.grow()
	}

	q.elems[q.tail] = a
	q.tail = q.next(q.tail)
}

func (q *RingBufferQueue) pop() int {
	if q.head == q.tail {
		panic("Trying to pop from an empty queue. Sorry!")
	}

	ans := q.elems[q.head]
	q.head = q.next(q.head)
	// TODO shrink :)
	return ans
}
