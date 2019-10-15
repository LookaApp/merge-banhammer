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
}