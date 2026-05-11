# Module Manifest Standard

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: `modules/<module>/module.yaml`

This document defines the first draft of the vibit module manifest.

A module manifest is the local contract between a module and the rest of the system. It tells agents what the module owns, what it exposes, what it depends on, what invariants it protects, and how it must be verified.

## 1. Purpose

Every implementation module must have:

```text
modules/<module>/module.yaml
```

The manifest exists so agents do not need to infer module boundaries from scattered code.

It should answer:

- What is this module?
- What category of module is it?
- What data and concepts does it own?
- What public commands and queries does it expose?
- What events does it publish and subscribe to?
- Which dependencies are allowed or forbidden?
- Which invariants must not be broken?
- Which files are generated?
- Which tests are required?

## 2. Minimal Example

```yaml
schema_version: 0.1
module: inventory
category: domain
status: draft
summary: Owns player inventory state and item capacity rules.

owns:
  entities:
    - inventory
    - inventory_item
  data:
    - inventory_records
  permissions:
    - inventory_read
    - inventory_write

public_api:
  commands:
    - AddItem
    - RemoveItem
  queries:
    - GetInventory

events:
  publishes:
    - ItemAdded
    - ItemRemoved
  subscribes:
    - PlayerCreated

dependencies:
  allowed:
    - player
    - currency
  forbidden:
    - match

invariants:
  - item_count_must_not_be_negative
  - inventory_capacity_must_not_exceed_limit

generated:
  files: []

tests:
  required:
    - command_tests
    - query_tests
    - event_tests
    - invariant_tests
    - migration_tests
```

## 3. Required Fields

### `schema_version`

Manifest schema version.

Initial value:

```yaml
schema_version: 0.1
```

### `module`

Stable module identifier.

Rules:

- Use `snake_case`.
- Match the directory name under `modules/`.
- Do not rename without an explicit migration.

### `category`

Module category.

Allowed initial values:

```text
domain
platform
integration
application
```

Definitions are maintained in `.arch/modules.yaml`.

### `status`

Lifecycle status.

Allowed initial values:

```text
draft
active
deprecated
removed
```

### `summary`

One-sentence explanation of the module's responsibility.

The summary should be useful to an agent deciding whether a change belongs in this module.

## 4. Ownership

`owns` declares the concepts, data, and permissions controlled by the module.

Example:

```yaml
owns:
  entities:
    - inventory
    - inventory_item
  data:
    - inventory_records
  permissions:
    - inventory_read
    - inventory_write
```

Rules:

- Owned data must not be directly modified by other modules.
- Ownership must be explicit before implementation code is added.
- Moving ownership between modules requires a change spec and maintainer approval.

## 5. Public API

`public_api` declares the module's intended external surface.

Example:

```yaml
public_api:
  commands:
    - AddItem
  queries:
    - GetInventory
```

Rules:

- Commands express intent to change state.
- Queries express read-only access.
- Public commands and queries must have schemas before implementation.
- Transport handlers should call module APIs instead of owning business logic.

## 6. Events

`events` declares facts published by the module and facts consumed from other modules.

Example:

```yaml
events:
  publishes:
    - ItemAdded
  subscribes:
    - PlayerCreated
```

Rules:

- Events express facts that already happened.
- Public events must be versioned when compatibility matters.
- Event payloads must be schema-defined.
- Subscriptions should be handled through policies or explicit handlers.

## 7. Dependencies

`dependencies` declares allowed and forbidden module dependencies.

Example:

```yaml
dependencies:
  allowed:
    - player
  forbidden:
    - match
```

Rules:

- A dependency not listed as allowed should be treated as disallowed until declared.
- Forbidden dependencies document especially risky or invalid coupling.
- Agents must not add imports or calls that violate this section.

## 8. Invariants

`invariants` declares business rules that must remain true.

Example:

```yaml
invariants:
  - item_count_must_not_be_negative
```

Rules:

- Use `snake_case`.
- Every invariant should eventually have tests.
- Invariants are stronger than implementation convenience.

## 9. Generated Files

`generated` declares files or directories produced by generators.

Example:

```yaml
generated:
  files:
    - commands/add_item.generated.ts
  directories:
    - clients
```

Rules:

- Agents must not hand-edit generated files unless changing the generator or documenting the reason.
- Generated files must be traceable to source schema.

## 10. Tests

`tests.required` declares required test categories.

Example:

```yaml
tests:
  required:
    - command_tests
    - query_tests
    - event_tests
    - invariant_tests
    - migration_tests
```

Initial allowed test categories:

```text
unit_tests
command_tests
query_tests
event_tests
contract_tests
invariant_tests
integration_tests
migration_tests
replay_tests
architecture_tests
```

Rules:

- Behavior changes should update relevant tests.
- If a required test category cannot run yet, record it as not available.
- Do not remove required test categories to make a change easier.

## 11. Agent Checklist

Before editing a module, an agent should check:

- Does the module manifest exist?
- Does the change belong to this module?
- Does the change affect ownership?
- Does the change add or modify a public command, query, event, error, or permission?
- Does a schema need to change first?
- Are dependencies still allowed?
- Which invariants are at risk?
- Which tests must be added or updated?
- Which generated files should not be hand-edited?

If the manifest is insufficient, improve it before changing implementation code.
