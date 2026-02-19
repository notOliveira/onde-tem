# My decisions for this project

### Architectural Decisions

- Codebase language: English (identifiers, comments, packages).

- API JSON format: snake_case.

- API envelope: { data, error, meta }.

- Language: Accept-Language header with fallback to en.

- Error codes: UPPER_SNAKE_CASE (stable identifiers).

- IDs: UUID.

- Timestamps: UTC ISO-8601 (RFC 3339, Z) in the API; stored as timestamptz in Postgres.

- Business hours: stored as local time (HH:MM) + establishment timezone (IANA).

### Commits

- test: for tests.

- feat: new features.

- refactor: changes with potentially high impact.

- style: removing comments, indenting code, formatting, etc.

- fix: bug fixes.

- chore: changes that do not involve code.

- docs: documentation changes.

- build: changes that affect the project build.

- perf: performance improvements.

- ci: CI/CD changes.

- revert: reverting to a previous commit/state.