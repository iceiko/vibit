# Impact

## Affected Areas

- Documentation and product roadmap.
- GitHub issue template configuration.
- Architecture manifests under `.arch/`.
- Repository check rules and `tools/vibit`.
- README and AGENTS continuation pointers.

## Runtime Impact

No runtime behavior changes are introduced.

This change does not modify Go runtime code, WebSocket transport behavior, Protobuf adapters, application dispatch, authentication/session behavior, persistence adapters, migrations, or startup behavior.

## Protocol Impact

No protocol routes, Protobuf source files, or generated output are added or changed.

## Data Impact

No migrations, repository interfaces, storage adapters, indexes, or data compatibility semantics are added or changed.

## Dependency Impact

No dependencies are added.

## Release And Outreach Impact

This change does not execute public announcements beyond the GitHub release record, run paid promotion, add hosted deployment, or create release artifacts.

It adds a feedback issue form only.

## Product Impact

The repository now records a product maturity path:

- source-first alpha: reached;
- prototype-ready foundation: next product stage;
- single-node production-candidate foundation: planned;
- Nakama/Pitaya-class product: long-term target.

This makes early feedback easier to triage without overstating the current alpha.

## Compatibility

No user-facing API, wire protocol, database, generated output, runtime behavior, or release artifact compatibility changes are introduced.
