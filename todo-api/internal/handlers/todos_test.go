package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"todo-api/internal/gen"
	"todo-api/internal/handlers"
	"todo-api/internal/store"
)

func newTestHandler() http.Handler {
	st := store.New()
	srv := handlers.NewTodoServer(st)
	r := chi.NewRouter()
	return gen.HandlerWithOptions(gen.NewStrictHandler(srv, nil), gen.ChiServerOptions{BaseRouter: r})
}

func doRequest(t *testing.T, h http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreateTodoRejectsBlankTitle(t *testing.T) {
	h := newTestHandler()

	for _, title := range []string{"", "   ", "\t\n"} {
		rec := doRequest(t, h, http.MethodPost, "/todos", gen.NewTodo{Title: title})
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("title %q: expected 400, got %d: %s", title, rec.Code, rec.Body.String())
		}
	}
}

func TestCreateAndListOldestFirst(t *testing.T) {
	h := newTestHandler()

	titles := []string{"first", "second", "third"}
	for _, title := range titles {
		rec := doRequest(t, h, http.MethodPost, "/todos", gen.NewTodo{Title: title})
		if rec.Code != http.StatusCreated {
			t.Fatalf("create %q: expected 201, got %d: %s", title, rec.Code, rec.Body.String())
		}
		var created gen.Todo
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("decode created todo: %v", err)
		}
		if created.Done {
			t.Fatalf("new todo should default done=false")
		}
		if created.ID == "" {
			t.Fatalf("new todo should have a server-generated id")
		}
	}

	rec := doRequest(t, h, http.MethodGet, "/todos", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var page struct {
		Count int        `json:"count"`
		Data  []gen.Todo `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if page.Count != len(titles) {
		t.Fatalf("expected count %d, got %d", len(titles), page.Count)
	}
	for i, title := range titles {
		if page.Data[i].Title != title {
			t.Fatalf("expected oldest-first order, position %d: got %q want %q", i, page.Data[i].Title, title)
		}
	}
}

func TestToggleTodo(t *testing.T) {
	h := newTestHandler()

	rec := doRequest(t, h, http.MethodPost, "/todos", gen.NewTodo{Title: "toggle me"})
	var created gen.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created todo: %v", err)
	}

	rec = doRequest(t, h, http.MethodPost, "/todos/"+created.ID+"/toggle", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("toggle: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var toggled gen.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &toggled); err != nil {
		t.Fatalf("decode toggled todo: %v", err)
	}
	if !toggled.Done {
		t.Fatalf("expected done=true after toggle")
	}

	// Toggle back.
	rec = doRequest(t, h, http.MethodPost, "/todos/"+created.ID+"/toggle", nil)
	var toggledBack gen.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &toggledBack); err != nil {
		t.Fatalf("decode toggled-back todo: %v", err)
	}
	if toggledBack.Done {
		t.Fatalf("expected done=false after second toggle")
	}
}

func TestToggleUnknownTodoIs404(t *testing.T) {
	h := newTestHandler()
	rec := doRequest(t, h, http.MethodPost, "/todos/does-not-exist/toggle", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteTodo(t *testing.T) {
	h := newTestHandler()

	rec := doRequest(t, h, http.MethodPost, "/todos", gen.NewTodo{Title: "delete me"})
	var created gen.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created todo: %v", err)
	}

	rec = doRequest(t, h, http.MethodDelete, "/todos/"+created.ID, nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d: %s", rec.Code, rec.Body.String())
	}

	rec = doRequest(t, h, http.MethodGet, "/todos", nil)
	var page struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if page.Count != 0 {
		t.Fatalf("expected empty store after delete, got count=%d", page.Count)
	}
}

func TestDeleteUnknownTodoIs404(t *testing.T) {
	h := newTestHandler()
	rec := doRequest(t, h, http.MethodDelete, "/todos/does-not-exist", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestListPagination(t *testing.T) {
	h := newTestHandler()

	for i := 0; i < 5; i++ {
		doRequest(t, h, http.MethodPost, "/todos", gen.NewTodo{Title: "item"})
	}

	rec := doRequest(t, h, http.MethodGet, "/todos?limit=2&offset=0", nil)
	var page struct {
		Count    int        `json:"count"`
		Next     string     `json:"next"`
		Previous string     `json:"previous"`
		Data     []gen.Todo `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if page.Count != 5 {
		t.Fatalf("expected count 5, got %d", page.Count)
	}
	if len(page.Data) != 2 {
		t.Fatalf("expected 2 items in page, got %d", len(page.Data))
	}
	if page.Next == "" {
		t.Fatalf("expected a next page URI")
	}
	if page.Previous != "" {
		t.Fatalf("expected no previous page URI on first page")
	}
}
