# Impact Analysis

## Affected Modules

- `runtime`
- `storage`
- `workflow`
- `reference`

## Module Ownership Impact

The storage module continues to own storage-neutral repository vocabulary and value types under `runtime/internal/modules/storage`.

The storage application service remains under `runtime/internal/app/storage` and is not changed by this proof. The PostgreSQL storage adapter remains under `runtime/internal/platform/persistence/postgres` and is not changed.

The proof adds a test-only in-memory storage repository inside `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go` so the local alpha proof can exercise the route family without creating a new production memory storage backend.

## Public Contract Impact

No public protocol contracts are added or changed.

The proof exercises the already implemented route family:

- query `storage.GetOwnStorageObject`
- query `storage.ListOwnStorageObjects`
- command `storage.PutOwnStorageObject`
- command `storage.DeleteOwnStorageObject`

The payload package remains `vibit.storage.v1`.

## Generated Output Impact

No generated files are added or changed.

The existing generated storage Protobuf output remains traced to `proto/vibit/storage/v1/storage.proto`.

## Data And Migration Impact

No migrations are added or changed.

The proof uses a test-only repository and does not change `storage_objects` SQL shape, PostgreSQL adapter behavior, transaction behavior, startup migrations, or live database requirements.

## Runtime Impact

The local alpha request-loop script now runs two focused Go E2E tests. It still executes inside the Go test runtime and does not start a hosted deployment or publish release artifacts.

The production runtime route registration and storage service behavior are unchanged.

## Authentication And Session Impact

The storage proof uses the existing local onboarding, device credential login, first-message connection binding, and authenticated request wrapper flow. It does not add handshake authentication, token carriers, session persistence changes, bound-session policy changes, or metadata-only identity trust.

## Nakama And Pitaya Alignment

Nakama informs the capability being proven: durable player-owned storage objects with collection/key/value/version semantics and put/get/list/delete operations.

Pitaya informs the layering being proven: transport frame handling, Protobuf adaptation, authenticated route protection, session metadata, application handlers, service behavior, and repository handoff stay separated.

The proof uses vibit route names and payloads. It does not copy Nakama or Pitaya public APIs.

## Test Impact

Adds `TestStorageObjectsProtocolRouteLocalAlphaFlow` to the authenticated gameplay E2E test file and extends the fixture with storage route registration.

The proof covers:

- local onboarding and login;
- first-message connection binding;
- authenticated `PutOwnStorageObject`;
- authenticated `GetOwnStorageObject`;
- authenticated `ListOwnStorageObjects`;
- authenticated `DeleteOwnStorageObject` with expected version;
- post-delete not-found response;
- secret and value redaction expectations for the error path.

## Compatibility Risks

The main risk is that the local proof could drift into a new storage backend. This is bounded by keeping the repository implementation test-only, not changing production startup, and recording the deferral of production memory storage behavior.
