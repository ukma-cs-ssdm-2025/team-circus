# Repository Guidelines

## Project Structure & Module Organization
- `backend/` houses the Go services: entrypoint in `cmd/main.go`, layered logic under `internal/{config,domain,service,handler}` and persistence in `repo/`; SQL migrations live in `db/migrations`.
- `frontend/` contains the Vite + React client with shared utilities in `src/utils`, typed API helpers under `src/services`, and page-level views in `src/pages`.
- `docs/` and `loom/` store stakeholder artifacts (team charter, demos, specs); keep generated coverage or Swagger output inside the existing backend `docs/` and `coverage/` folders.

## Build, Test & Development Commands
- Bootstrap env via `task copy:env`, `task docker:postgres:up`, `task docker:migrator:up`.
- Backend: `task back:download` + `task back:vendor` resolve dependencies, `task back:run` starts the API locally, and `task back:build` emits `backend/bin/backend`.
- Frontend: `task front:install`, `task front:dev` for hot reload, `task front:build` for production assets, `task front:preview` to smoke-test the build (all run via Bun 1.3+).
- Use `task --list-all` to inspect every helper command.

## Coding Style & Naming Conventions
- Go code must remain `gofmt`/`goimports` clean; GolangCI (`task back:lint`) enforces `errcheck`, `revive`, `staticcheck`, and other rules plus a 140-char max line length. Package names stay lowercase and short; exported symbols use `CamelCase`.
- React/TypeScript follows ESLint config in `frontend/eslint.config.js`: modern ESM syntax, React Hooks best practices, and `typescript-eslint` rules. Components and hooks live in PascalCase files (e.g., `MarkdownCanvas.tsx`), shared utilities in camelCase modules.

## Testing Guidelines
- Unit tests: `task back:test:unit` (or `...coverage`) runs `_test.go` suites and writes coverage to `backend/coverage/coverage.out`; generate HTML with `task back:test:unit:coverage:html`.
- Functional API tests live under `backend/tests/api`; bring up dependencies via `task back:test:func:up` and execute targeted suites with `task back:test:func:run -- ./tests/api/...`.
- Frontend currently relies on manual verification; when automated UI tests arrive, colocate them next to the component as `ComponentName.test.tsx`.

## Commit & Pull Request Guidelines
- Follow the lightweight Conventional Commit style already in history (`docs: add link to video`, `fix: sonar issues`, `feat: auth middleware`).
- Each commit should focus on one concern and include updated docs/tests when behavior changes.
- Pull requests need: a summary of changes, referenced issues (e.g., `Closes #123`), screenshots or curl examples for UI/API shifts, and confirmation that `task back:test:unit` (plus any affected front-end commands) succeeded.

## Environment & Security Notes
- Never commit `.env` files; use `.env.sample` plus `task copy:env:optional` to sync local variables.
- Database credentials flow through Taskfile variables—run the Docker-based test targets so secrets remain containerized, and keep generated artifacts inside tracked `docs/` or `coverage/` paths.
