# ADR-0007: YAML Contract Source Format

Status: Accepted  
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-inventory-contracts/`

Related conversations:

- `conversations/2026-05-12-inventory-contracts.md`

Related artifacts:

- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `contracts/`
- `modules/inventory/module.yaml`

## Context

ADR-0005 established that the first runtime slice should declare commands, queries, events, errors, and permissions before handwritten behavior or generated output. It left the first schema source format open.

The inventory proof slice is now prepared at the module level and needs concrete contract source files for `GrantItem`, `GetInventory`, `ItemGranted`, inventory errors, and inventory permissions.

## Decision

vibit will use YAML as the first source format for runtime contracts.

YAML contract files may contain JSON-Schema-like payload definitions for request, response, event, and metadata shapes. The YAML file is the source contract. Generated TypeScript types, validators, dispatch shapes, clients, and test fixtures must trace back to these YAML files.

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
- Delay the format decision until runtime implementation.

## Rationale

YAML is readable for humans and agents while still being structured enough for future generators. It is less noisy than raw JSON for early business contracts, and it can still embed JSON-Schema-like shapes where strict payload structure matters.

OpenAPI is useful for HTTP APIs, but vibit needs contracts for backend behavior beyond transport routes. Protobuf may become useful later, but choosing it before the first proof slice would add tooling weight before the project has proven the workflow.

## Agent Reasoning Summary

The first contract source should reduce ambiguity without forcing a heavy runtime toolchain. YAML gives agents a compact source of truth and leaves room for generated TypeScript and validation later.

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

- The first inventory contracts are YAML source files.
- Future generators must treat these YAML files as source, not generated output.
- Contract file paths should be registered in `.arch/contracts.yaml`.
- If YAML becomes insufficient, a future ADR may supersede this decision before generated output becomes large.

## Reversal Conditions

Revisit this decision if YAML contracts become hard to validate, hard to generate from, or too ambiguous for agents during the first runtime implementation.

## Follow-Up

- Add a schema or check for contract files after the first contract shape stabilizes.
- Generate TypeScript types and validators from these contracts when runtime scaffolding starts.
