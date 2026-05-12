# ADR-0007: YAML Contract Source Format

Status: Accepted  
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-inventory-contracts/`
- `changes/2026-05-12-ratify-go-websocket-protobuf-runtime/`

Related conversations:

- `conversations/2026-05-12-inventory-contracts.md`

Related artifacts:

- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `contracts/`
- `modules/inventory/module.yaml`
- `decisions/ADR-0009-websocket-protobuf-client-protocol.md`

## Context

ADR-0005 established that the first runtime slice should declare commands, queries, events, errors, and permissions before handwritten behavior or generated output. It left the first schema source format open.

The inventory proof slice is now prepared at the module level and needs concrete contract source files for `GrantItem`, `GetInventory`, `ItemGranted`, inventory errors, and inventory permissions.

## Decision

vibit will use YAML as the first source format for semantic runtime contracts.

YAML contract files may contain JSON-Schema-like payload definitions for request, response, event, and metadata shapes. The YAML file is the source contract for business semantics: command/query/event identity, ownership, permissions, errors, invariants, and module boundaries.

`ADR-0009` adds Protobuf as the first wire message format. Protobuf files own wire-level message shape and compatibility. Tooling must check alignment between YAML semantic contracts and Protobuf wire schemas before the protocol surface grows.

Initial contract files should live under:

```text
contracts/<module>/<contract-type>/*.yaml
```

The initial contract types are:

- `commands`
- `queries`
- `events`
- `errors`
- `permissions`

## Alternatives Considered

- Write raw JSON Schema files as the source format.
- Use OpenAPI as the only source format.
- Use Protobuf as the first contract source.
- Use Protobuf as the first wire schema while keeping YAML as the semantic contract source.
- Delay the format decision until runtime implementation.

## Rationale

YAML is readable for humans and agents while still being structured enough for future generators. It is less noisy than raw JSON for early business contracts, and it can still embed JSON-Schema-like shapes where strict payload structure matters.

OpenAPI is useful for HTTP APIs, but vibit needs contracts for backend behavior beyond transport routes. Protobuf is now selected for wire schema, but it should not replace the semantic contract manifest because wire schemas alone do not express all agent-relevant architecture context.

## Agent Reasoning Summary

The first semantic contract source should reduce ambiguity without forcing all architecture meaning into transport schemas. YAML gives agents a compact source of truth for business semantics. Protobuf complements it at the wire boundary.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: low
  reversibility: medium
  long_term_maintainability: high
confidence: medium
```

## Consequences

- The first inventory semantic contracts are YAML source files.
- Future generators must treat these YAML files as source, not generated output.
- Contract file paths should be registered in `.arch/contracts.yaml`.
- Wire-level Protobuf files must be aligned with these contracts, not used to bypass them.
- If YAML becomes insufficient as a semantic source, a future ADR may supersede this decision before generated output becomes large.

## Reversal Conditions

Revisit this decision if YAML contracts become hard to validate, hard to generate from, or too ambiguous for agents during the first runtime implementation.

## Follow-Up

- Add a schema or check for contract files after the first contract shape stabilizes.
- Define the manifest-to-proto consistency check before broad protocol generation.
- Generate Go types, Protobuf mappings, validators, and fixtures from the appropriate source artifacts when runtime scaffolding starts.
