screen TodoList "View, add, complete and remove todos"
  navbar "Todo App"
  heading "My Todos"
  row
    input "Add a new todo..."
    button "Add" primary
  table "Done | Title | Added"
    row "☐ | Buy groceries | just now"
    row "☑ | Read a book | yesterday"
  text "Click a row's checkbox to toggle done; click Delete to remove a todo."

flow "Manage todos"
  description "A user views, adds, completes and removes todos"
  TodoList
