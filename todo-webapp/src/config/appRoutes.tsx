import { type RouteProps, Navigate } from "react-router-dom";
import AppLayout from "../layouts/AppLayout";
import TodoListPage from "../pages/TodoListPage";

export interface AppRoute extends Omit<RouteProps, "children"> {
  children?: AppRoute[];
  label?: string;
}

// One screen — TodoList — per wireframes.dsl. No sign-in, no routing beyond
// this single screen: every other path redirects here.
const appRoutes: AppRoute[] = [
  {
    element: <AppLayout />,
    children: [
      { path: "/", element: <TodoListPage />, label: "My Todos" },
      { path: "*", element: <Navigate to="/" replace /> },
    ],
  },
];

export default appRoutes;
