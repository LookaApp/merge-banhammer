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

func (q *mergeQueue) Dequeue() *string {
	if q.Length() == 0 {
		return nil
	}

	dequeuedValue := q.queue[0]
	q.queue = append(make([]string, 0, 12), q.queue[1:]...)

	return &dequeuedValue
}

func (q *mergeQueue) Entries() []string {
	return q.queue
}

func (q *mergeQueue) Length() int {
	return len(q.queue)
}
