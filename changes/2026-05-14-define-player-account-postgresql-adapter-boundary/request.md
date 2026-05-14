# Request

## Original Request

继续推进

## Clarified Requirement

Advance `W-0053` by defining the PostgreSQL adapter implementation boundary for the player account repository after the module-owned repository interface exists.

## User-Visible Outcome

Maintainers and agents can inspect the future player account PostgreSQL adapter boundary before implementation.

The boundary declares:

- Planned adapter source: `runtime/internal/platform/persistence/postgres/player_account_repository.go`
- Planned focused tests: `runtime/internal/platform/persistence/postgres/player_account_repository_test.go`
- Constructor: `NewPlayerAccountRepositoryForUnitOfWork(executor)`
- Interface implemented: `player.Repository`
- Executor ownership: caller-supplied application unit-of-work executor
- Transaction rule: no adapter-owned `BEGIN`, `COMMIT`, or `ROLLBACK`
- First SQL scope: `player_accounts`, `player_account_events`, and current lifecycle reads only

## Non-Goals

- Do not implement the PostgreSQL player account repository adapter.
- Do not add runtime player account command or query handlers.
- Do not add WebSocket routes or handshake behavior.
- Do not add authentication providers.
- Do not add token behavior, token storage, refresh tokens, or signing metadata.
- Do not add credential or password storage.
- Do not add external identity linking.
- Do not add runtime session persistence.
- Do not change the Protobuf envelope.
- Do not change the ratified player account migration schema, table ownership, repository interface shape, or transaction ownership model.
- Do not make live PostgreSQL verification mandatory in default repository checks.
- Do not copy Nakama or Pitaya public API shapes.

## Unknowns

- Exact application-facing player account error wrappers remain future implementation detail, but adapter error categories are now required.
- Runtime player account handler and protocol bridge sequencing remains future work after adapter implementation.
- Live PostgreSQL adapter integration remains opt-in through `VIBIT_POSTGRES_TEST_DSN`.

## Acceptance Criteria

- [x] Define the player account PostgreSQL adapter boundary in the persistence standard.
- [x] Record planned adapter source, test path, constructor shape, executor ownership, transaction rule, SQL scope, and error mapping expectations.
- [x] Update runtime, player, contracts, module, and conventions manifests.
- [x] Update English and Simplified Chinese runtime/player guidance.
- [x] Update runtime identity checks so the boundary is required and future adapter files are restricted to `runtime/internal/platform/persistence/postgres/`.
- [x] Verify without implementing the adapter, player handlers, WebSocket routes, authentication, tokens, credentials, session persistence, Protobuf envelope changes, or WebSocket handshake changes.
