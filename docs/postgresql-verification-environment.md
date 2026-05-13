# PostgreSQL Verification Environment Standard

Status: Draft v0.1  
Last updated: 2026-05-13  
Scope: Disposable PostgreSQL environment rules for live persistence and migration verification

This standard defines how agents and maintainers verify PostgreSQL-backed behavior without relying on maintainer memory.

The paired Simplified Chinese translation is `docs/postgresql-verification-environment.zh-CN.md`. The English file is authoritative.

Use this standard together with `docs/postgresql-persistence-boundary.md`, `.arch/runtime.yaml`, `.arch/conventions.yaml`, `ADR-0011`, `ADR-0013`, and `ADR-0020`.

## 1. Purpose

Source checks are not live database verification.

`node tools/vibit check migrations` validates migration files, ownership traces, and source conventions. It does not prove that migrations apply cleanly to PostgreSQL, that repository SQL runs against a real server, or that transaction behavior survives real database constraints.

This standard defines the smallest explicit environment contract for those live checks while keeping normal repository verification fast and local.

## 2. Dependency Position

Disposable PostgreSQL verification is optional by default.

Rules:

- Normal unit tests and `node tools/vibit check all` must not require a running PostgreSQL server.
- Live PostgreSQL checks must be opt-in through explicit environment variables.
- Docker, Podman, system PostgreSQL, or another service manager may be used by the human or agent, but none is a required project dependency in this standard.
- Cloud-hosted PostgreSQL is not assumed by this standard.
- Agents must not hide credentials in tracked files.

## 3. Environment Variables

The live verification environment is described by these variables:

```text
VIBIT_POSTGRES_TEST_DSN
VIBIT_POSTGRES_TEST_DATABASE
VIBIT_POSTGRES_TEST_CLEANUP
VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE
```

### `VIBIT_POSTGRES_TEST_DSN`

Required for live PostgreSQL verification.

This DSN points to a disposable database or database server that the agent is allowed to mutate during verification.

Rules:

- The value must come from explicit local environment input.
- The value must not be stored in tracked files.
- If unset, live PostgreSQL checks must skip and record that no disposable environment was available.

### `VIBIT_POSTGRES_TEST_DATABASE`

Optional database name used by scripts or manual setup instructions when the DSN does not already fully identify the target database.

Recommended default:

```text
vibit_test
```

### `VIBIT_POSTGRES_TEST_CLEANUP`

Optional cleanup mode.

Allowed values:

```text
drop_schema
truncate
keep
```

Default:

```text
drop_schema
```

Rules:

- `drop_schema` is the preferred disposable mode because it leaves the database ready for a clean rerun.
- `truncate` may be used when dropping schema objects is unavailable.
- `keep` must be recorded in verification output because it leaves state behind for inspection.

### `VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE`

Required for destructive setup or cleanup.

Allowed values:

```text
0
1
```

Default:

```text
0
```

Rules:

- Destructive setup or cleanup must not run unless this value is `1`.
- Destructive means dropping schemas, dropping databases, truncating module-owned tables, or applying rollback checks that remove schema objects.
- If a check needs destructive behavior and the value is not `1`, it must skip or fail with a clear explanation rather than guessing intent.

## 4. Verification Categories

PostgreSQL verification has three categories.

### Source Verification

Source verification is always safe to run by default:

```bash
node tools/vibit check migrations
node tools/vibit check postgres-env
```

It validates repository artifacts and standards. It does not open a database connection.

### Unit Verification

Unit verification is safe to run by default:

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime
```

Unit tests may use fake executors or fake database handles. They must not require `VIBIT_POSTGRES_TEST_DSN`.

### Live PostgreSQL Verification

Live PostgreSQL verification is opt-in.

It may include:

```text
migration status against VIBIT_POSTGRES_TEST_DSN
migration apply against VIBIT_POSTGRES_TEST_DSN
repository integration tests against VIBIT_POSTGRES_TEST_DSN
transaction runner integration tests against VIBIT_POSTGRES_TEST_DSN
rollback or cleanup checks when destructive verification is explicitly allowed
```

Until live verification commands exist, agents must record:

```text
Not verified: live PostgreSQL verification is unavailable because no repository command exists yet.
```

or, if commands exist but `VIBIT_POSTGRES_TEST_DSN` is unset:

```text
Not verified: live PostgreSQL verification skipped because VIBIT_POSTGRES_TEST_DSN was not set.
```

## 5. Setup Expectations

This standard intentionally does not prescribe one service manager.

Any setup path is acceptable if it provides:

- PostgreSQL compatible with the first runtime's migration and repository code.
- A disposable database or schema that can be mutated.
- A DSN supplied through `VIBIT_POSTGRES_TEST_DSN`.
- Cleanup behavior matching `VIBIT_POSTGRES_TEST_CLEANUP`.
- Explicit destructive permission through `VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1` when cleanup can remove schema or data.

Agents may document the exact local command they used, but they must not commit local service credentials, socket paths, passwords, or tokens.

## 6. Cleanup Rules

Agents must leave the environment in a known state.

Rules:

- Prefer isolated databases or isolated schemas.
- Prefer cleaning created schema objects after a successful live verification run.
- If cleanup is skipped intentionally, record why.
- If cleanup fails, record it as a verification risk.
- Never run destructive cleanup against a DSN that was not explicitly provided for disposable verification.

## 7. Recording Verification

Every change that touches PostgreSQL migrations, repositories, transaction boundaries, or persistent runtime composition must record one of:

```text
Verified: live PostgreSQL verification ran with VIBIT_POSTGRES_TEST_DSN set.
Not verified: live PostgreSQL verification skipped because VIBIT_POSTGRES_TEST_DSN was not set.
Not verified: live PostgreSQL verification command is not implemented yet.
Not applicable: change did not touch PostgreSQL-backed behavior.
```

When live verification runs, record:

- Commands run.
- Whether cleanup ran.
- Whether destructive cleanup was allowed.
- Any skipped integration tests.

Do not record the DSN itself.

## 8. Current Repository Tooling

The current static environment-standard check is:

```bash
node tools/vibit check postgres-env
```

This check verifies that the disposable PostgreSQL verification standard, runtime manifest references, and guidance artifacts exist. It does not connect to PostgreSQL.

Live migration and repository integration commands remain future work.
