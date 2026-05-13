# Impact Analysis

## Affected Modules

- `inventory` remains the first proof slice.
- No module ownership is changed.
- Future module planning now has a broader reference baseline covering account/session, social, storage, matchmaking, realtime, authoritative match, presence, chat, economy, leaderboard, and cluster-related capabilities.

## Module Ownership Impact

No current module ownership changes.

The impact is planning-level: future module proposals should be checked against the reference alignment standard so capability growth is intentional and not only proof-slice driven.

## Public Contract Impact

No command, query, event, error, or permission contracts change.

Future public contracts should be planned with the reference capability matrix in mind, but still follow vibit's contract-first rules.

## Data And Migration Impact

No migrations are added.

No persistence ownership changes are made.

## Test Impact

No runtime tests are added because this is a standards and planning change.

Repository checks should verify documentation, memory, schemas, and architecture consistency.

`tools/vibit check architecture` is extended to require the new reference manifest, standard documents, ADR, and entry-point references. This keeps the reference baseline visible to future agents during normal repository intake.

## Documentation Impact

Adds a new reference alignment standard and its Simplified Chinese translation.

Adds an ADR to make the reference decision durable.

Updates README, AGENTS, and `.arch/README` so future agents discover the alignment during intake.

Updates `.arch/conventions.yaml` so the reference alignment standard is discoverable as a repository convention.

## Compatibility Risks

The main risk is confusing "reference baseline" with "copy the external framework."

The standard explicitly avoids API copying, dependency adoption, and premature distributed architecture. Nakama and Pitaya guide capability planning; vibit remains governed by its own constitution and Agent-Native constraints.
