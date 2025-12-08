# Repository Guidelines

## Project Structure & Module Organization

- Backend Go code lives in `backend/` (`cmd/main.go` entrypoint, business logic under `internal/`, shared helpers in `pkg/`, SQL migrations in `db/migrations/`, black-box integration suites in `tests/`).
- Frontend React (Vite + TS) sits in `frontend/` (`src/` for UI, `public/` for static assets). Shared docs are under `docs/`; demo recordings in `loom/`.

## Build, Test, and Development Commands

- First-time setup: `task copy:env` to copy sample env files, then adjust `backend/.env`.
- Backend: `task back:run` to start the API; `task back:lint` for golangci-lint; `task back:test:unit` (or `...:coverage`) for unit tests; `task back:test:func` or `...:up` for functional API flows; `task docker:backend:up` to run the stack; `task docker:down` to stop.
- Frontend: `task front:dev` for the Vite dev server, `task front:build` for production build. Lint with `bun run lint`.

## Coding Style & Naming Conventions

- Go: keep code `gofmt`/`goimports` clean and lint-passing. Exported identifiers use UpperCamelCase; locals lowerCamelCase; avoid package-level globals. Respect context propagation and explicit error handling; prefer lines under 140 chars.
- Frontend: TypeScript-first. Components in `frontend/src/components/` use PascalCase filenames; hooks/utils use camelCase. Prefer functional components and React hooks. Keep styling consistent with MUI patterns already in use.
- Formatting: format frontend code with Biome (`bunx biome format ...`); run ESLint before shipping.

## Testing Guidelines

- Go unit tests sit beside implementations as `*_test.go`, table-driven with testify assertions. Run via `task back:test:unit`; update coverage with `task back:test:unit:coverage`.
- API changes should be mirrored in `backend/tests/api/` and run with `task back:test:func` (spin dependencies with `task back:test:func:up` if needed).
- Frontend linting is the main gate; add component tests if behavior becomes complex.

## Commit & Pull Request Guidelines

- Use Conventional Commits (e.g., `feat(frontend): add login form`). Scope where practical.
- PRs should describe intent, list key changes, link issues, and note env/schema updates. Include screenshots or Loom clips for UI tweaks and state which tasks/tests ran or why they were skipped. Keep secrets out of history.

## Environment & Configuration Tips

- For local Postgres only: `task docker:postgres:up`; stop containers with `task docker:down`.
- Document any frontend config quirks in `frontend/ENV_SETUP.md` to keep setups reproducible. Ensure `.env` values match your running stack before testing.
