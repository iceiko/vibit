# Request

## Original Request

The maintainer asked to continue development for up to ten work items. The current next-ready work item is `W-0019: Define disposable PostgreSQL verification environment`.

## Clarified Requirement

Define the repository standard for optional live PostgreSQL verification so agents know how to run or skip database-backed checks without relying on maintainer memory.

The standard must keep live PostgreSQL verification opt-in, avoid making Docker, Podman, cloud PostgreSQL, or any external service manager mandatory, and distinguish source validation from live database verification.

## User-Visible Outcome

Maintainers and agents can read a single PostgreSQL verification environment standard, know which environment variables are required for live checks, know how cleanup should be handled, and run a static repository check that confirms the standard is registered.

## Non-Goals

- Do not make Docker, Podman, or another service manager a required project dependency.
- Do not make live PostgreSQL integration tests mandatory in default checks.
- Do not add cloud-hosted PostgreSQL assumptions.
- Do not implement live database integration tests in this change.
- Do not wire PostgreSQL into normal server startup.

## Unknowns

- The exact live migration and repository integration commands remain future work.
- Persistent runtime composition remains deferred until `W-0020`.

## Acceptance Criteria

- Add an English PostgreSQL verification environment standard and a Simplified Chinese translation.
- Define required and optional environment variables for live PostgreSQL verification.
- Define cleanup expectations and destructive-operation safeguards.
- Add a static `node tools/vibit check postgres-env` command.
- Register the new command and rule metadata in repository guidance.
- Update architecture manifests and module/runtime guidance.
- Mark `W-0019` complete and move `W-0020` to `next_ready`.
