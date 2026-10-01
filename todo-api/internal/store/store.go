// Package store holds the process-local, in-memory todo list. Data lives for
// the lifetime of the process only — there is no persistence layer.
package store

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"todo-api/internal/gen"
)

// Store is a mutex-guarded, insertion-ordered list of todos. Items are always
// appended on create, so the backing slice is already oldest-first.
type Store struct {
	mu    sync.Mutex
	items []gen.Todo
}

// New returns an empty Store.
func New() *Store {
	return &Store{}
}

// List returns up to limit items starting at offset, oldest first, plus the
// total count of items in the store.
func (s *Store) List(limit, offset int) ([]gen.Todo, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := len(s.items)
	if offset < 0 {
		offset = 0
	}
	if offset > count {
		offset = count
	}
	end := offset + limit
	if end > count {
		end = count
	}
	if end < offset {
		end = offset
	}

	page := make([]gen.Todo, end-offset)
	copy(page, s.items[offset:end])
	return page, count
}

// Create appends a new todo with the given title and returns it.
func (s *Store) Create(title string) gen.Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := gen.Todo{
		ID:        newID(),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now().UTC(),
	}
	s.items = append(s.items, todo)
	return todo
}

// Toggle flips the done flag of the todo with the given id. The second return
// value is false when no such todo exists.
func (s *Store) Toggle(id string) (gen.Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Done = !s.items[i].Done
			return s.items[i], true
		}
	}
	return gen.Todo{}, false
}

// Delete removes the todo with the given id. It returns false when no such
// todo exists.
func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.items {
		if s.items[i].ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

// newID returns a random, URL-safe identifier.
func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand failing is effectively unrecoverable; fall back to a
		// timestamp so the service degrades rather than panics.
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(b)
}
