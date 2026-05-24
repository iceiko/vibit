# Impact Analysis

## Nakama Product Capability Impact

This pilot selects `friends_groups_and_parties` as the concrete Nakama-style capability family for the scaffolded intake. The specific future capability is a player friendship relationship lifecycle: request, accept, reject, remove, block, unblock, list, and read relationship status.

The product impact is planning-only in this slice. It records why the friend lifecycle should be defined before broader social behavior, chat targeting, group membership, party invites, matchmaking filters, or match runtime social context.

## Pitaya Impact

Pitaya remains deferred. This change does not add distributed topology, frontend/backend split, RPC, service discovery, groups, cluster routing, distributed sessions, Pitaya-style backend services, or direct Pitaya API compatibility.

## Affected Modules

- `runtime`: only manifests, AGENTS guidance, and repository checks are updated.
- `reference`: Nakama-first roadmap state is advanced from scaffold implementation to a friendship lifecycle gate follow-up.
- `repository_workflow`: the next-ready queue advances from W-0230 to W-0231.
- `social`: future ownership candidate only; no module or runtime package is created in this pilot.

## Module Ownership Impact

No runtime module ownership changes are made. The future W-0231 gate must decide whether friendship lifecycle belongs in a new domain module, an existing player/social boundary, or a staged contract-first package. Storage remains explicitly not the owner.

## Public Contract Impact

No command, query, event, error, permission, route, protocol, or schema is added in this pilot. The future gate is expected to define candidate semantic vocabulary before any Protobuf or generated output work.

Potential future semantic vocabulary for W-0231:

- commands: send friend request, accept friend request, reject friend request, remove friend, block player, unblock player;
- queries: list friend relationships, get relationship status;
- events: friend request created, friend request accepted, friend request rejected, friend removed, player blocked, player unblocked;
- errors: invalid target, self-target forbidden, duplicate request, blocked relationship, invalid transition, relationship not found;
- permissions: validated player identity required.

## Data And Migration Impact

No data ownership, persistence, migration, index, or adapter changes are made. W-0231 must define data and migration expectations before any table, repository, or PostgreSQL adapter work. Future data planning must avoid logging private relationship graph details beyond what is necessary for tests and durable records.

## Test Impact

This pilot requires repository checks only. Future W-0231 test planning should include:

- positive lifecycle behavior tests;
- negative invalid-state tests;
- permission and authenticated identity tests;
- persistence and transaction tests before schema work;
- protocol mapping tests after semantic contract ratification;
- redaction tests for private relationship and authentication data;
- idempotency and concurrency expectations for friend state transitions;
- local alpha E2E proof after route implementation is authorized.

## Documentation And Memory Impact

This change updates:

- scaffolded change artifacts under `changes/2026-05-24-pilot-scaffolded-nakama-feature-request-intake/`;
- `ADR-0138`;
- conversation memory for the pilot;
- `.arch/work-items.yaml`, `.arch/runtime.yaml`, `.arch/reference.yaml`, `.arch/conventions.yaml`, `.arch/contracts.yaml`, and `.arch/modules.yaml`;
- repository AGENTS guides and storage module guide;
- alpha, product maturity, and Nakama roadmap docs;
- `tools/vibit` and `rules/check-rules.json`.

## Compatibility Risks

No API, event, data, protocol, SDK, hosted, distributed runtime, or direct Nakama/Pitaya compatibility changes are made. Direct compatibility remains false unless a later ADR explicitly authorizes it.
