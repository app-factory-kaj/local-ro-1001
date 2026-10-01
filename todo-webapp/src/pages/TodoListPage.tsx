import { useCallback, useEffect, useState, type JSX, type FormEvent } from "react";
import {
  PageContent,
  PageTitle,
  ListingTable,
  TextField,
  Button,
  Checkbox,
  IconButton,
  Typography,
  Stack,
  formatRelativeTime,
} from "@wso2/oxygen-ui";
import { Trash2 } from "@wso2/oxygen-ui-icons-react";
import { todoApi } from "../api";
import type { components } from "../generated/todo-api";

type Todo = components["schemas"]["Todo"];

export default function TodoListPage(): JSX.Element {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);

  const [title, setTitle] = useState("");
  const [titleError, setTitleError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Ids with a toggle or delete in flight, so their row controls disable
  // rather than let a second click race the first.
  const [pending, setPending] = useState<Set<string>>(new Set());

  const loadTodos = useCallback(async () => {
    setLoading(true);
    setLoadError(null);
    const { data, error } = await todoApi.GET("/todos");
    if (error) {
      setLoadError("Could not load todos. Please try again.");
    } else {
      setTodos(data.data);
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    void loadTodos();
  }, [loadTodos]);

  const handleAdd = async (event: FormEvent): Promise<void> => {
    event.preventDefault();
    const trimmed = title.trim();
    if (!trimmed) {
      setTitleError("Enter a todo title.");
      return;
    }
    setTitleError(null);
    setSubmitting(true);
    const { data, error } = await todoApi.POST("/todos", {
      body: { title: trimmed },
    });
    setSubmitting(false);
    if (error) {
      setTitleError(error.message || "Could not add that todo.");
      return;
    }
    setTodos((current) => [...current, data]);
    setTitle("");
  };

  const withPending = (id: string, run: () => Promise<void>): void => {
    setPending((current) => new Set(current).add(id));
    void run().finally(() => {
      setPending((current) => {
        const next = new Set(current);
        next.delete(id);
        return next;
      });
    });
  };

  const handleToggle = (todo: Todo): void => {
    withPending(todo.id, async () => {
      const { data, error } = await todoApi.POST("/todos/{todoId}/toggle", {
        params: { path: { todoId: todo.id } },
      });
      if (error) {
        setLoadError("Could not update that todo. Please try again.");
        return;
      }
      setTodos((current) => current.map((item) => (item.id === data.id ? data : item)));
    });
  };

  const handleDelete = (todo: Todo): void => {
    withPending(todo.id, async () => {
      const { error } = await todoApi.DELETE("/todos/{todoId}", {
        params: { path: { todoId: todo.id } },
      });
      if (error) {
        setLoadError("Could not delete that todo. Please try again.");
        return;
      }
      setTodos((current) => current.filter((item) => item.id !== todo.id));
    });
  };

  return (
    <PageContent>
      <PageTitle>
        <PageTitle.Header>My Todos</PageTitle.Header>
      </PageTitle>

      <Stack component="form" direction="row" spacing={2} sx={{ mb: 3 }} onSubmit={handleAdd}>
        <TextField
          placeholder="Add a new todo..."
          value={title}
          onChange={(event) => {
            setTitle(event.target.value);
            if (titleError) setTitleError(null);
          }}
          error={Boolean(titleError)}
          helperText={titleError ?? undefined}
          fullWidth
          size="small"
        />
        <Button type="submit" variant="contained" disabled={submitting}>
          Add
        </Button>
      </Stack>

      {loadError && (
        <Typography color="error" sx={{ mb: 2 }}>
          {loadError}
        </Typography>
      )}

      <ListingTable.Container>
        <ListingTable>
          <ListingTable.Head>
            <ListingTable.Row>
              <ListingTable.Cell>Done</ListingTable.Cell>
              <ListingTable.Cell>Title</ListingTable.Cell>
              <ListingTable.Cell>Added</ListingTable.Cell>
              <ListingTable.Cell />
            </ListingTable.Row>
          </ListingTable.Head>
          <ListingTable.Body>
            {loading ? (
              <ListingTable.Row>
                <ListingTable.Cell colSpan={4}>
                  <Typography color="text.secondary">Loading todos…</Typography>
                </ListingTable.Cell>
              </ListingTable.Row>
            ) : todos.length === 0 ? (
              <ListingTable.Row>
                <ListingTable.Cell colSpan={4}>
                  <ListingTable.EmptyState
                    title="No todos yet"
                    description="Add your first todo above."
                  />
                </ListingTable.Cell>
              </ListingTable.Row>
            ) : (
              todos.map((todo) => (
                <ListingTable.Row key={todo.id}>
                  <ListingTable.Cell>
                    <Checkbox
                      checked={todo.done}
                      disabled={pending.has(todo.id)}
                      onChange={() => {
                        handleToggle(todo);
                      }}
                      inputProps={{ "aria-label": `Mark "${todo.title}" as done` }}
                    />
                  </ListingTable.Cell>
                  <ListingTable.Cell
                    sx={
                      todo.done
                        ? { textDecoration: "line-through", color: "text.secondary" }
                        : undefined
                    }
                  >
                    {todo.title}
                  </ListingTable.Cell>
                  <ListingTable.Cell>{formatRelativeTime(new Date(todo.createdAt))}</ListingTable.Cell>
                  <ListingTable.Cell>
                    <ListingTable.RowActions>
                      <IconButton
                        aria-label={`Delete "${todo.title}"`}
                        size="small"
                        disabled={pending.has(todo.id)}
                        onClick={() => {
                          handleDelete(todo);
                        }}
                      >
                        <Trash2 size={18} />
                      </IconButton>
                    </ListingTable.RowActions>
                  </ListingTable.Cell>
                </ListingTable.Row>
              ))
            )}
          </ListingTable.Body>
        </ListingTable>
      </ListingTable.Container>

      <Typography color="text.secondary" sx={{ mt: 2 }}>
        Click a row&apos;s checkbox to toggle done; click Delete to remove a todo.
      </Typography>
    </PageContent>
  );
}
