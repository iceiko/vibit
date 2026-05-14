# Request

## Original Request

它区域限制呗。按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。

## Clarified Requirement

Advance `W-0081` by defining the authentication PostgreSQL adapter boundary after the storage-neutral authentication repository interface exists.

## User-Visible Outcome

Maintainers and agents can inspect the future authentication PostgreSQL adapter boundary before implementation.

The boundary declares:

- Planned adapter source: `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- Planned focused tests: `runtime/internal/platform/persistence/postgres/authentication_repository_test.go`
- Constructor: `NewAuthenticationRepositoryForUnitOfWork(executor)`
- Interface implemented: `authentication.Repository`
- Executor ownership: caller-supplied application unit-of-work executor
- Transaction rule: no adapter-owned `BEGIN`, `COMMIT`, or `ROLLBACK`
- First SQL scope: credential and token verifier persistence operations only

## Non-Goals

- Do not implement the PostgreSQL authentication adapter.
- Do not add runtime credential lookup, token issuance, token validation, logout, refresh, cleanup, handlers, routes, or production authentication behavior.
- Do not add WebSocket routes, proof carriers, or handshake behavior.
- Do not add Protobuf messages or generated authentication shapes.
- Do not add authentication dependencies.
- Do not change ratified repository interface shape, migration schemas, table ownership, or transaction ownership model.
- Do not make live PostgreSQL verification mandatory in default repository checks.

## Unknowns

- Exact adapter error wrapper names remain future implementation detail.
- Live PostgreSQL adapter integration remains opt-in through `VIBIT_POSTGRES_TEST_DSN`.
- Runtime authentication handler sequencing remains future work after adapter implementation is explicitly authorized.

## Acceptance Criteria

- [x] Define the authentication PostgreSQL adapter boundary in the persistence standard.
- [x] Record planned adapter source, test path, constructor shape, executor ownership, transaction rule, SQL scope, and error mapping expectations.
- [x] Update runtime, contracts, conventions, reference, authentication module, and guides.
- [x] Update runtime checks so the boundary is required after the repository interface exists and future adapter files are restricted to `runtime/internal/platform/persistence/postgres/`.
- [x] Verify without implementing the adapter or runtime authentication behavior.
