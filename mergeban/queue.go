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

func (q *mergeQueue) Withdraw(valueToWithdraw string) *string {
	if q.Length() == 0 {
		return nil
	}

	withdrawPosition := int(-1)

	for position, enqueuedValue := range q.queue {
		if enqueuedValue == valueToWithdraw {
			withdrawPosition = position
		}
	}

	if withdrawPosition == -1 {
		return nil
	}

	withdrawnValue := q.queue[withdrawPosition]
	empty := []string{}
	var leadingValues, trailingValues []string

	if withdrawPosition == 0 {
		leadingValues = empty
		trailingValues = q.queue[1:]
	} else if withdrawPosition == len(q.queue)-1 {
		leadingValues = q.queue[0:withdrawPosition]
		trailingValues = empty
	} else {
		leadingValues = q.queue[0:withdrawPosition]
		trailingValues = q.queue[withdrawPosition+1 : len(q.queue)]
	}

	newQueue := append(make([]string, 0, 12), leadingValues...)
	q.queue = append(newQueue, trailingValues...)

	return &withdrawnValue
}

func (q *mergeQueue) Entries() []string {
	return q.queue
}

func (q *mergeQueue) Length() int {
	return len(q.queue)
}
