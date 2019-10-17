package mergeban

type mergeQueue struct {
	queue []mergeQueueEntry
}

type mergeQueueEntry struct {
	ResponseURL string
	UserID      string
	UserName    string
}

func NewQueue() *mergeQueue {
	return &mergeQueue{queue: make([]mergeQueueEntry, 0, 12)}
}

func (q *mergeQueue) Enqueue(userIDToEnqueue, userNameToEnqueue, responseURL string) {
	for _, enqueuedValue := range q.queue {
		if enqueuedValue.UserID == userIDToEnqueue {
			return
		}
	}

	q.queue = append(q.queue, mergeQueueEntry{
		ResponseURL: responseURL,
		UserID:      userIDToEnqueue,
		UserName:    userNameToEnqueue,
	})
}

func (q *mergeQueue) Dequeue() *mergeQueueEntry {
	if q.Length() == 0 {
		return nil
	}

	dequeuedValue := q.queue[0]
	q.queue = append(make([]mergeQueueEntry, 0, 12), q.queue[1:]...)

	return &dequeuedValue
}

func (q *mergeQueue) Withdraw(userIDToWithdraw string) int {
	if q.Length() == 0 {
		return -1
	}

	withdrawPosition := int(-1)

	for position, enqueuedValue := range q.queue {
		if enqueuedValue.UserID == userIDToWithdraw {
			withdrawPosition = position
		}
	}

	if withdrawPosition == -1 {
		return -1
	}

	empty := []mergeQueueEntry{}
	var leadingValues, trailingValues []mergeQueueEntry

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

	newQueue := append(make([]mergeQueueEntry, 0, 12), leadingValues...)
	q.queue = append(newQueue, trailingValues...)

	return withdrawPosition
}

func (q *mergeQueue) FindIndex(userID string) int {
	for index, enqueuedValue := range q.queue {
		if enqueuedValue.UserID == userID {
			return index
		}
	}

	return -1
}

func (q *mergeQueue) Peek() *mergeQueueEntry {
	if q.Length() == 0 {
		return nil
	}

	return &q.queue[0]
}

func (q *mergeQueue) UserIDs() []string {
	var acc []string

	for _, entry := range q.queue {
		acc = append(acc, entry.UserID)
	}

	return acc
}

func (q *mergeQueue) UserNames() []string {
	var acc []string

	for _, entry := range q.queue {
		acc = append(acc, entry.UserName)
	}

	return acc
}

func (q *mergeQueue) Length() int {
	return len(q.queue)
}
