# vibit

vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.

Status: constitutional design phase

## What This Project Means By Agent-Native

Agent-native does not primarily mean that the server has AI features.

It means the server architecture is designed so AI coding agents can work inside it reliably:

- The architecture is explicit instead of tribal knowledge.
- Module ownership is declared instead of guessed.
- Public behavior is contract-first.
- Repeatable structure is generated.
- Business rules are tested as invariants.
- Cross-module communication is bounded.
- Change workflow is documented and verifiable.
- Documentation is written for both humans and agents.

AI gameplay features such as NPC agents, memory, model routing, tool calling, and simulations may become extensions later. They are not the foundation.

## Why This Exists

Many existing server codebases were built for human maintainers with local context, long memory, and implicit team conventions. AI coding agents can help in those codebases, but they often lose force when architecture rules are hidden, module boundaries are weak, tests are incomplete, or public contracts are unclear.

vibit starts from a different premise:

> The next generation of long-lived server software should be designed so agents can safely understand, modify, verify, and extend it.

The goal is not to make agents magically smarter. The goal is to make the codebase more legible, bounded, generated, contract-driven, and testable.

## Current Documents

- `CONSTITUTION.md`: canonical project constitution
- `CONSTITUTION.zh-CN.md`: Simplified Chinese translation
- `AGENTS.md`: repository-level operating guide for coding agents
- `AGENTS.zh-CN.md`: Simplified Chinese translation
- `.arch/README.md`: machine-readable architecture manifest entry point
- `.arch/modules.yaml`: first draft module registry manifest
- `.arch/conventions.yaml`: first draft repository convention manifest
- `.arch/runtime.yaml`: runtime readiness manifest for the first reference implementation
- `.arch/contracts.yaml`: contract registry for public command, query, event, error, and permission source files
- `docs/module-manifest.md`: module manifest standard
- `docs/module-manifest.zh-CN.md`: Simplified Chinese translation
- `docs/change-spec.md`: change spec standard
- `docs/change-spec.zh-CN.md`: Simplified Chinese translation
- `changes/_template/`: reusable change spec template
- `docs/conversation-log.md`: conversation log standard
- `docs/conversation-log.zh-CN.md`: Simplified Chinese translation
- `conversations/`: maintainer-agent project memory
- `docs/agent-decision-record.md`: Agent Decision Record standard
- `docs/agent-decision-record.zh-CN.md`: Simplified Chinese translation
- `decisions/`: durable decision rationale
- `docs/schema-validation.md`: schema validation standard
- `docs/schema-validation.zh-CN.md`: Simplified Chinese translation
- `schema/`: JSON Schema files for machine-checkable standards
- `rules/`: rule catalogs for machine-readable check metadata

English documents are canonical. Simplified Chinese translations are maintained for human readers and early project discussion.

## Intended Direction

vibit should evolve toward:

- Architecture manifests under `.arch/`
- A first TypeScript/Node.js reference runtime governed by `.arch/runtime.yaml` and Agent Decision Records
- Module manifests at `modules/<module>/module.yaml`, following `docs/module-manifest.md`
- Module-level agent guides at `modules/<module>/AGENTS.md`
- Contract-first commands, queries, events, errors, permissions, and migrations
- Contract source files under `contracts/`, registered by `.arch/contracts.yaml`
- Generated scaffolds for repeatable framework structure
- Architecture checks that verify dependency, contract, event, and generated-file rules
- Change specs under `changes/<date>-<change-id>/`, following `docs/change-spec.md`
- Conversation logs under `conversations/`, following `docs/conversation-log.md`
- Agent Decision Records under `decisions/`, following `docs/agent-decision-record.md`
- Schema validation under `schema/`, following `docs/schema-validation.md`
- Rule catalogs under `rules/`, starting with `rules/check-rules.json`

The first serious prototype should prove one claim:

> Given a new backend requirement, an AI coding agent can identify the affected module, update the correct contracts, generate the correct structure, implement the behavior, add tests, run verification, and update documentation without damaging unrelated architecture.

## CLI Prototype

The first executable standard lives at:

```bash
tools/vibit
```

Initial commands:

```bash
node tools/vibit --help
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check schemas
node tools/vibit check schemas --json
node tools/vibit check memory
node tools/vibit check memory --json
node tools/vibit check contracts
node tools/vibit check contracts --json
node tools/vibit check generated
node tools/vibit check generated --json
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit inspect module inventory
node tools/vibit inspect boundary --from inventory --to player
node tools/vibit inspect contract --module inventory --type command --id GrantItem
node tools/vibit inspect change bootstrap-vibit-cli
node tools/vibit inspect memory
node tools/vibit inspect rule check.subcheck
node tools/vibit inspect rules --category check
node tools/vibit generate contract --module inventory --type command --id GrantItem
node tools/vibit check architecture
node tools/vibit check architecture --json
node tools/vibit check change bootstrap-vibit-cli
node tools/vibit check change bootstrap-vibit-cli --json
node tools/vibit check module inventory
node tools/vibit check module inventory --json
node tools/vibit generate module <module>
```

The CLI currently uses Node.js standard-library APIs only. It is a prototype for architecture checks and module skeleton generation, not a server runtime.

Use `--json` when an agent needs machine-readable check results during intake, verification, or handoff. Human-readable text output remains the default.

Each JSON check result item includes a stable `rule_id` and an `artifact` value so agents can route failures without parsing prose. `check all --json` is a compact overview; run the specific failing check with `--json` to get full result details.

Use `node tools/vibit check memory` to verify required conversation log and Agent Decision Record structure.

Use `node tools/vibit check contracts` to verify that `.arch/contracts.yaml` and registered contract source files are consistent.

Use `node tools/vibit check generated` to verify that module-declared generated files exist and include generated, source, and generator trace markers.

Use `node tools/vibit check runtime` to run the current module runtime tests.

Use `node tools/vibit inspect contract --module <module> --type <type> --id <id>` to inspect one registered command, query, event, error catalog, or permission catalog as JSON during agent intake.

Use `node tools/vibit generate contract --module <module> --type <type> --id <id>` to regenerate contract shapes from contract source files.

Use `node tools/vibit inspect change <change-id>` to inspect a change spec directory and its verification metadata without manually opening every file.

Use `node tools/vibit inspect memory` to list change specs, conversation logs, and Agent Decision Records as a machine-readable project memory index.

Rule metadata for check output lives in `rules/check-rules.json`.

Use `node tools/vibit inspect rule <rule-id>` to inspect one rule without parsing the full catalog.

Use `node tools/vibit inspect rules` or `node tools/vibit inspect rules --category <category>` to discover available rules.

The first reference runtime is TypeScript on Node.js, using a modular monolith single-process server model. This is a reference implementation decision, not a permanent restriction on the broader architecture standard. See `.arch/runtime.yaml` and `decisions/ADR-0003-first-reference-runtime-language.md` through `decisions/ADR-0006-first-runtime-proof-slice.md`.

## Early Reference Domain

A small game backend is the recommended first demonstration domain because it naturally contains state, permissions, events, consistency rules, and long-lived modules.

Suggested modules:

- Player accounts
- Inventory
- Currency
- Rewards
- Tasks or quests
- Match sessions

The demo should emphasize maintainability and agent workflow over feature count.

## Governance

Project decisions are governed by `CONSTITUTION.md`.

Before changing constitutional principles, ratifying the name, introducing a major architectural pattern, or making a breaking standard change, read the constitution and record the motivation, alternatives, compatibility impact, and migration path.

## Name

`vibit` is the product name.

The intended category phrase is:

```text
agent-native server framework
```

Before final ratification, the name should be checked against major public registries and platforms for obvious conflicts.
