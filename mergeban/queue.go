package mergeban

type mergeQueue struct {
	queue []string
}

func NewQueue() *mergeQueue {
	return &mergeQueue{queue: make([]string, 0, 12)}
}

func (q *mergeQueue) Enqueue(valueToEnqueue string) {
	for _, enqueuedValue := range q.queue {
		if enqueuedValue == valueToEnqueue {
			return
		}
	}

	q.queue = append(q.queue, valueToEnqueue)
}

func (q *mergeQueue) Entries() []string {
	return q.queue
}

func (q *mergeQueue) Length() int {
	return len(q.queue)
}
