# Request

## Original Request

```text
继续
```

## Clarified Requirement

Advance `W-0233 Add friends relationship migration source` by adding the first PostgreSQL migration source for `friend_relationships`, following the accepted friends relationship persistence schema gate.

## User-Visible Outcome

Maintainers and agents can inspect a concrete SQL migration source for future friends relationship current-state persistence.

## Non-Goals

- Adding friends relationship runtime behavior.
- Adding repository interfaces.
- Adding PostgreSQL friends relationship adapters.
- Adding protocol routes.
- Adding Protobuf source files.
- Adding generated output.
- Adding dependencies.
- Adding automatic startup migration behavior.
- Changing authentication/session behavior.
- Adding delivery guarantees.
- Adding stream subscriptions.
- Adding chat rooms.
- Adding groups.
- Adding parties.
- Adding broadcast fanout.
- Adding matchmaking.
- Adding match runtime.
- Adding operations/admin behavior.
- Publishing SDKs or generated client libraries.
- Adding event/audit tables.
- Adding hosted deployments.
- Creating release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, or SDK packages.
- Executing public announcements beyond the GitHub release record.
- Running paid promotion.
- Adding Pitaya-style distributed architecture.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- Storage-neutral repository shape remains deferred.
- PostgreSQL adapter mapping remains deferred.
- Runtime lifecycle behavior remains deferred.
- Protocol route shape remains deferred.
- Event/audit table shape remains deferred.
- Groups, parties, chat targeting, matchmaking filters, and match social context remain future capability work.

## Acceptance Criteria

- [x] `runtime/migrations/postgres/000007_create_friend_relationships.sql` exists.
- [x] The migration declares goose Up and Down markers.
- [x] The migration declares `-- Module: friends`.
- [x] The migration creates only `friend_relationships`.
- [x] Required pair identity, lifecycle state, block timestamp, relationship version, and timestamp columns are present.
- [x] The migration enforces canonical pair ordering.
- [x] The migration constrains lifecycle state to `pending`, `friends`, `rejected`, and `removed`.
- [x] The migration adds canonical pair uniqueness.
- [x] The migration adds pair-member/lifecycle-state indexes and an updated-at index.
- [x] No raw secret, digest, transport metadata, chat, group, party, match, Pitaya, or Nakama compatibility columns are added.
- [x] No Go repository, adapter, runtime behavior, Protobuf, generated output, dependency, event/audit table, hosted deployment, release artifact, public announcement, or paid promotion is added by this slice.
