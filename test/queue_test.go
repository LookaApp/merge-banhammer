package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mergeban/pkg"
)

func TestQueue(t *testing.T) {
	nullResponseURL := ""
	t.Run("enqueue - can enqueue an entry when empty", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		entries := queue.UserIDs()

		assert.Equal(t, 1, len(entries))
		assert.Equal(t, "1", entries[0])
	})

	t.Run("enqueue - does nothing if the same value is enqueued twice", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		queue.Enqueue("1", "John Doe", nullResponseURL)
		entries := queue.UserIDs()

		assert.Equal(t, 1, len(entries))
		assert.Equal(t, "1", entries[0])
	})

	t.Run("enqueue - can enqueue multiple values and return them in FIFO order", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		queue.Enqueue("2", "John Doe", nullResponseURL)
		entries := queue.UserIDs()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "1", entries[0])
		assert.Equal(t, "2", entries[1])
	})

	t.Run("enqueue - stores the name of the enqueued user", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "Alice", nullResponseURL)
		queue.Enqueue("2", "Bob", nullResponseURL)
		entries := queue.UserNames()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "Alice", entries[0])
		assert.Equal(t, "Bob", entries[1])
	})

	t.Run("dequeue - does nothing if the queue is empty", func(t *testing.T) {
		queue := mergeban.NewQueue()

		dequeuedValue := queue.Dequeue()
		entries := queue.UserIDs()

		assert.Equal(t, 0, len(entries))
		assert.Nil(t, dequeuedValue)
	})

	t.Run("dequeue - returns the first enqueued value and removes it from the singleton queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		dequeuedValue := queue.Dequeue()
		entries := queue.UserIDs()

		assert.Equal(t, 0, len(entries))
		assert.Equal(t, "1", dequeuedValue.UserID)
	})

	t.Run("dequeue - returns the first enqueued value and preserves the order of the remaining values", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		queue.Enqueue("2", "John Doe", nullResponseURL)
		queue.Enqueue("3", "John Doe", nullResponseURL)
		dequeuedValue := queue.Dequeue()
		entries := queue.UserIDs()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "1", dequeuedValue.UserID)
		assert.Equal(t, "2", entries[0])
		assert.Equal(t, "3", entries[1])
	})

	t.Run("withdraw - does nothing if the queue is empty", func(t *testing.T) {
		queue := mergeban.NewQueue()

		withdrawnIndex := queue.Withdraw("-1")
		entries := queue.UserIDs()

		assert.Equal(t, 0, len(entries))
		assert.Equal(t, -1, withdrawnIndex)
	})

	t.Run("withdraw - does nothing if the provided value is not present in the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		withdrawnIndex := queue.Withdraw("999")
		entries := queue.UserIDs()

		assert.Equal(t, 1, len(entries))
		assert.Equal(t, "1", entries[0])
		assert.Equal(t, -1, withdrawnIndex)
	})

	t.Run("withdraw - removes the provided value from a singleton queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		withdrawnIndex := queue.Withdraw("1")
		entries := queue.UserIDs()

		assert.Equal(t, 0, len(entries))
		assert.Equal(t, 0, withdrawnIndex)
	})

	t.Run("withdraw - can withdraw from the head of the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		queue.Enqueue("2", "John Doe", nullResponseURL)
		queue.Enqueue("3", "John Doe", nullResponseURL)
		withdrawnIndex := queue.Withdraw("1")
		entries := queue.UserIDs()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, 0, withdrawnIndex)
		assert.Equal(t, "2", entries[0])
		assert.Equal(t, "3", entries[1])
	})

	t.Run("withdraw - can withdraw from the middle of the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		queue.Enqueue("2", "John Doe", nullResponseURL)
		queue.Enqueue("3", "John Doe", nullResponseURL)
		withdrawnIndex := queue.Withdraw("2")
		entries := queue.UserIDs()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, 1, withdrawnIndex)
		assert.Equal(t, "1", entries[0])
		assert.Equal(t, "3", entries[1])
	})

	t.Run("withdraw - can withdraw from the end of the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		queue.Enqueue("2", "John Doe", nullResponseURL)
		queue.Enqueue("3", "John Doe", nullResponseURL)
		withdrawnIndex := queue.Withdraw("3")
		entries := queue.UserIDs()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, 2, withdrawnIndex)
		assert.Equal(t, "1", entries[0])
		assert.Equal(t, "2", entries[1])
	})

	t.Run("findIndex - returns -1 if the provided entry is not present in the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		index := queue.FindIndex("42")

		assert.Equal(t, -1, index)
	})

	t.Run("findIndex - the index of the provided entry in the queue, if it exists", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", nullResponseURL)
		queue.Enqueue("2", "John Doe", nullResponseURL)
		queue.Enqueue("3", "John Doe", nullResponseURL)
		index := queue.FindIndex("2")
		entries := queue.UserIDs()

		assert.Equal(t, 1, index)
		assert.Equal(t, "2", entries[index])
	})

	t.Run("peek - returns nil if the queue is empty", func(t *testing.T) {
		queue := mergeban.NewQueue()

		head := queue.Peek()

		assert.Nil(t, head)
	})

	t.Run("peek - returns the head of the queue otherwise", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1", "John Doe", "http://example.com")
		head := queue.Peek()

		assert.Equal(t, "1", head.UserID)
		assert.Equal(t, "http://example.com", head.ResponseURL)
	})
}
