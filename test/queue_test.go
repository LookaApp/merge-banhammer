package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mergeban/mergeban"
)

func TestQueue(t *testing.T) {
	t.Run("enqueue - can enqueue an entry when empty", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		entries := queue.Entries()

		assert.Equal(t, 1, len(entries))
		assert.Equal(t, "1", entries[0])
	})

	t.Run("enqueue - does nothing if the same value is enqueued twice", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		queue.Enqueue("1")
		entries := queue.Entries()

		assert.Equal(t, 1, len(entries))
		assert.Equal(t, "1", entries[0])
	})

	t.Run("enqueue - can enqueue multiple values and return them in FIFO order", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		queue.Enqueue("2")
		entries := queue.Entries()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "1", entries[0])
		assert.Equal(t, "2", entries[1])
	})

	t.Run("dequeue - does nothing if the queue is empty", func(t *testing.T) {
		queue := mergeban.NewQueue()

		dequeuedValue := queue.Dequeue()
		entries := queue.Entries()

		assert.Equal(t, 0, len(entries))
		assert.Nil(t, dequeuedValue)
	})

	t.Run("dequeue - returns the first enqueued value and removes it from the singleton queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		dequeuedValue := queue.Dequeue()
		entries := queue.Entries()

		assert.Equal(t, 0, len(entries))
		assert.Equal(t, "1", *dequeuedValue)
	})

	t.Run("dequeue - returns the first enqueued value and preserves the order of the remaining values", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		queue.Enqueue("2")
		queue.Enqueue("3")
		dequeuedValue := queue.Dequeue()
		entries := queue.Entries()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "1", *dequeuedValue)
		assert.Equal(t, "2", entries[0])
		assert.Equal(t, "3", entries[1])
	})

	t.Run("withdraw - does nothing if the queue is empty", func(t *testing.T) {
		queue := mergeban.NewQueue()

		withdrawnValue := queue.Withdraw("-1")
		entries := queue.Entries()

		assert.Equal(t, 0, len(entries))
		assert.Nil(t, withdrawnValue)
	})

	t.Run("withdraw - does nothing if the provided value is not present in the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		withdrawnValue := queue.Withdraw("999")
		entries := queue.Entries()

		assert.Equal(t, 1, len(entries))
		assert.Equal(t, "1", entries[0])
		assert.Nil(t, withdrawnValue)
	})

	t.Run("withdraw - removes the provided value from a singleton queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		withdrawnValue := queue.Withdraw("1")
		entries := queue.Entries()

		assert.Equal(t, 0, len(entries))
		assert.Equal(t, "1", *withdrawnValue)
	})

	t.Run("withdraw - can withdraw from the head of the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		queue.Enqueue("2")
		queue.Enqueue("3")
		withdrawnValue := queue.Withdraw("1")
		entries := queue.Entries()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "1", *withdrawnValue)
		assert.Equal(t, "2", entries[0])
		assert.Equal(t, "3", entries[1])
	})

	t.Run("withdraw - can withdraw from the middle of the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		queue.Enqueue("2")
		queue.Enqueue("3")
		withdrawnValue := queue.Withdraw("2")
		entries := queue.Entries()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "2", *withdrawnValue)
		assert.Equal(t, "1", entries[0])
		assert.Equal(t, "3", entries[1])
	})

	t.Run("withdraw - can withdraw from the end of the queue", func(t *testing.T) {
		queue := mergeban.NewQueue()

		queue.Enqueue("1")
		queue.Enqueue("2")
		queue.Enqueue("3")
		withdrawnValue := queue.Withdraw("3")
		entries := queue.Entries()

		assert.Equal(t, 2, len(entries))
		assert.Equal(t, "3", *withdrawnValue)
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

		queue.Enqueue("1")
		queue.Enqueue("2")
		queue.Enqueue("3")
		index := queue.FindIndex("2")
		entries := queue.Entries()

		assert.Equal(t, 1, index)
		assert.Equal(t, "2", entries[index])
	})
}
