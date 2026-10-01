// Package handlers implements the todo-api operations against the in-memory
// store, behind the generated strict server interface.
package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"todo-api/internal/gen"
	"todo-api/internal/store"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// TodoServer implements gen.StrictServerInterface against an in-memory Store.
type TodoServer struct {
	store *store.Store
}

// NewTodoServer builds a TodoServer backed by the given Store.
func NewTodoServer(s *store.Store) *TodoServer {
	return &TodoServer{store: s}
}

// ListTodos returns a page of todos, oldest first.
func (s *TodoServer) ListTodos(ctx context.Context, request gen.ListTodosRequestObject) (gen.ListTodosResponseObject, error) {
	limit := request.Params.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	offset := request.Params.Offset
	if offset < 0 {
		offset = 0
	}

	items, count := s.store.List(limit, offset)

	resp := gen.ListTodos200JSONResponse{
		Count: count,
		Data:  items,
	}
	if offset+limit < count {
		resp.Next = pageURI(offset+limit, limit)
	}
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		resp.Previous = pageURI(prevOffset, limit)
	}

	return resp, nil
}

// CreateTodo adds a new todo. A blank or whitespace-only title is rejected.
func (s *TodoServer) CreateTodo(ctx context.Context, request gen.CreateTodoRequestObject) (gen.CreateTodoResponseObject, error) {
	if request.Body == nil || strings.TrimSpace(request.Body.Title) == "" {
		return gen.CreateTodo400JSONResponse{
			Code:        http.StatusBadRequest,
			Message:     "blank title",
			Description: "title is required and must not be blank or whitespace-only",
		}, nil
	}

	todo := s.store.Create(request.Body.Title)
	return gen.CreateTodo201JSONResponse(todo), nil
}

// ToggleTodo flips a todo's done status.
func (s *TodoServer) ToggleTodo(ctx context.Context, request gen.ToggleTodoRequestObject) (gen.ToggleTodoResponseObject, error) {
	todo, ok := s.store.Toggle(request.TodoID)
	if !ok {
		return gen.ToggleTodo404JSONResponse{
			Code:        http.StatusNotFound,
			Message:     "todo not found",
			Description: "no todo with id " + request.TodoID,
		}, nil
	}
	return gen.ToggleTodo200JSONResponse(todo), nil
}

// DeleteTodo removes a todo.
func (s *TodoServer) DeleteTodo(ctx context.Context, request gen.DeleteTodoRequestObject) (gen.DeleteTodoResponseObject, error) {
	if !s.store.Delete(request.TodoID) {
		return gen.DeleteTodo404JSONResponse{
			Code:        http.StatusNotFound,
			Message:     "todo not found",
			Description: "no todo with id " + request.TodoID,
		}, nil
	}
	return gen.DeleteTodo204Response{}, nil
}

// pageURI builds a relative /todos URI for the given page parameters.
func pageURI(offset, limit int) string {
	v := url.Values{}
	v.Set("limit", strconv.Itoa(limit))
	v.Set("offset", strconv.Itoa(offset))
	return "/todos?" + v.Encode()
}
