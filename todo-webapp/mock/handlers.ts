import { http, HttpResponse } from "msw";
import type { components } from "../src/generated/todo-api";

type Todo = components["schemas"]["Todo"];

// State lives in the PAGE, not a server: setupWorker resolves every request in
// the page's own JS context, so a reload (or any full page load) re-runs this
// module and resets the seed below. Only in-app navigation — adding, toggling,
// deleting without a reload — carries a change forward. That is also what
// makes a verification run repeatable.
//
// Seeded oldest first, matching the contract's `listTodos` ordering and the
// titles/done-state wireframes.dsl draws for TodoList.
let todos: Todo[] = [
  {
    id: "1",
    title: "Read a book",
    done: true,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 26).toISOString(), // yesterday
  },
  {
    id: "2",
    title: "Buy groceries",
    done: false,
    createdAt: new Date(Date.now() - 1000 * 30).toISOString(), // just now
  },
];
let nextId = 3;

export const handlers = [
  // Most specific first: /todos/:id/toggle before /todos/:id.
  http.get("/api/todos", ({ request }) => {
    const url = new URL(request.url);
    const limit = Number(url.searchParams.get("limit") ?? 20);
    const offset = Number(url.searchParams.get("offset") ?? 0);
    const page = todos.slice(offset, offset + limit);
    return HttpResponse.json({
      count: todos.length,
      next: offset + limit < todos.length ? `/todos?limit=${String(limit)}&offset=${String(offset + limit)}` : null,
      previous: offset > 0 ? `/todos?limit=${String(limit)}&offset=${String(Math.max(0, offset - limit))}` : null,
      data: page,
    });
  }),

  http.post("/api/todos", async ({ request }) => {
    const input = (await request.json()) as { title?: string };
    const title = input?.title?.trim();
    if (!title) {
      return HttpResponse.json(
        { code: 400, message: "Blank or missing title" },
        { status: 400 },
      );
    }
    const created: Todo = {
      id: String(nextId++),
      title,
      done: false,
      createdAt: new Date().toISOString(),
    };
    todos = [...todos, created];
    return HttpResponse.json(created, { status: 201 });
  }),

  http.post("/api/todos/:todoId/toggle", ({ params }) => {
    const todo = todos.find((t) => t.id === params.todoId);
    if (!todo) {
      return HttpResponse.json(
        { code: 404, message: "No todo with that id" },
        { status: 404 },
      );
    }
    const updated: Todo = { ...todo, done: !todo.done };
    todos = todos.map((t) => (t.id === updated.id ? updated : t));
    return HttpResponse.json(updated);
  }),

  http.delete("/api/todos/:todoId", ({ params }) => {
    const before = todos.length;
    todos = todos.filter((t) => t.id !== params.todoId);
    return before === todos.length
      ? HttpResponse.json({ code: 404, message: "No todo with that id" }, { status: 404 })
      : new HttpResponse(null, { status: 204 });
  }),
];
