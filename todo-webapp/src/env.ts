// Typed read of window._env_, the platform's runtime config. This app has no
// OIDC dependency, no `configurations.env` defaults and no external-kind
// dependency, so there is nothing to declare on Env — but /env-config.js is
// still mounted by the platform (and served by mock/plugin.ts in mock mode),
// and its absence is still a real failure worth throwing on.
//
// eslint-disable-next-line @typescript-eslint/no-empty-interface
type Env = Record<string, never>;

declare global {
  interface Window {
    _env_: Env;
  }
}

if (!window._env_) {
  throw new Error(
    "window._env_ not set — /env-config.js failed to load. " +
      "The platform mounts this file; if you see this locally, host " +
      "/env-config.js from your dev server.",
  );
}

export const env: Env = window._env_;
