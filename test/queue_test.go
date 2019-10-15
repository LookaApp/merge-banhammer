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
}
