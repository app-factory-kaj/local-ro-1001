# Todo App — PRD

## Problem Statement

People juggling small, everyday tasks often reach for heavyweight todo tools when all they need is a quick place to jot a task down, see what's outstanding, and check it off. The overhead of accounts, sync, and configuration gets in the way of just tracking a short list of things to do.

## Solution

A minimal todo app: a small Go REST API that keeps todos (title and done status) in memory, and a single-page React web app that calls it to list, add, complete, and remove todos. No accounts, no database — just a fast, disposable list.

## Actors

- **User** — anyone using the app. There is no sign-in, so every user has the same, full access to the shared todo list.

## User Stories

1. As a User, I want to view the list of todos, so that I can see everything I need to do.
2. As a User, I want to add a new todo by entering a title, so that I can track a new task.
3. As a User, I want to toggle a todo's done status, so that I can mark a task complete or incomplete.
4. As a User, I want to delete a todo, so that I can remove a task I no longer need.

## Product Decisions

- No sign-in: the app is open and unauthenticated. This overrides the organization's default of SSO via Thunder for web apps, per the explicit brief.
- Storage: todos are stored in memory within the Go API process; no database is used, and all data is lost when the API restarts.
- Todo fields are limited to title (string) and done (boolean) — no due dates, priorities, categories, or tags.
- New todos default to done = false. *assumed*
- The API rejects creating a todo with a blank/whitespace-only title, returning a validation error. *assumed*
- The todo list is returned in creation order, oldest first. *assumed*

## Out of Scope

- User accounts, sign-in, or per-user data separation.
- Persistent storage (database or file) of any kind.
- Editing a todo's title after creation.
- Due dates, priorities, categories, tags, sorting, or filtering options.
- Multi-user collaboration or sharing of lists.

## Open Questions

None at this time.