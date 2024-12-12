package queue

type MonotonicQueue struct {
	queue []int
}

func (q *MonotonicQueue) Empty() bool {
	return q.Len() == 0
}

func (q *MonotonicQueue) Len() int {
	return len(q.queue)
}

func (q *MonotonicQueue) Push(val int) {
	n := len(q.queue)
	for n > 0 && val > q.queue[n-1] {
		q.queue = q.queue[:n-1]
		n--
	}
	q.queue = append(q.queue, val)
}

func (q *MonotonicQueue) Pop(val int) {
	n := len(q.queue)
	if n > 0 && q.queue[0] == val {
		q.queue = q.queue[1:]
	}
}
