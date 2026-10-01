// Command todo-api serves the Todo API: an in-memory, process-local store of
// todos with no persistence and no authentication.
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"todo-api/internal/gen"
	"todo-api/internal/handlers"
	"todo-api/internal/store"
)

const port = "9090"

func main() {
	st := store.New()
	srv := handlers.NewTodoServer(st)

	r := chi.NewRouter()
	handler := gen.HandlerWithOptions(gen.NewStrictHandler(srv, nil), gen.ChiServerOptions{BaseRouter: r})

	log.Printf("todo-api listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
