# Impact

## Architecture Impact

This change records a runtime readiness decision, not an implementation.

PostgreSQL becomes the first authoritative durable relational store for runtime state. S3-compatible object storage becomes the planned abstraction for large object artifacts. MinIO is recorded as the preferred local/self-hosted candidate, but it remains subject to dependency adoption before becoming a mandatory runtime dependency.

## Module Impact

`inventory` remains the preferred first proof slice.

The inventory module now records that PostgreSQL is the first authoritative store when runtime persistence begins. It also records that S3-compatible object storage is not required for the first inventory slice.

Domain modules must not import PostgreSQL drivers, S3 SDKs, or MinIO clients directly.

## Public Contract Impact

No command, query, event, error, or permission contract changes.

The persistence direction affects future implementation and migration contracts, not current public behavior.

## Data And Migration Impact

No migration is introduced by this change.

Before persistent repositories are implemented, the project still needs:

- PostgreSQL driver adoption.
- Migration tool or migration file conventions.
- Transaction boundary standard.
- Repository interface layout.
- Storage verification commands.

## Dependency Impact

No dependency is added.

PostgreSQL itself is ratified as a persistence platform. Specific Go packages for connecting to PostgreSQL remain undecided.

S3-compatible object storage is ratified as an abstraction. MinIO remains a candidate because the public `minio/minio` repository is archived as of 2026-04-25 and reports AGPL-3.0 license metadata, so adopting it as required infrastructure needs explicit review.

## Documentation Impact

Architecture manifests, README files, repository AGENTS guides, inventory module metadata, an ADR, and a conversation log are updated.

English public documents and Simplified Chinese translations are updated together where paired files exist.

## Compatibility

No public runtime API exists yet, so there is no external API, event, or data compatibility break.

The decision reduces future ambiguity before runtime persistence code starts.
