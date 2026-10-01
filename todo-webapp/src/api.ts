import createClient from "openapi-fetch";
import type { paths } from "./generated/todo-api";

// Same-origin: nginx proxies /api to the todo-api sibling (see
// nginx/15-aep-api-proxy.sh and nginx/default.conf). Never the public gateway
// URL and never a window._env_ key — todo-api is a `component`-kind
// dependency, not `external`.
export const todoApi = createClient<paths>({ baseUrl: "/api" });
