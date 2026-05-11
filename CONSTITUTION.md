# Agent-Native Server Constitution

Status: Draft v0.2  
Last updated: 2026-05-12  
Scope: Open-source agent-native game/backend server architecture
Product name: vibit

## 0. Purpose

This document is the constitution of the project.

Its job is not to describe a single implementation detail. Its job is to define the principles, boundaries, standards, and decision rules that every implementation detail must follow.

This project exists to build a server framework that AI coding agents can understand, extend, verify, and maintain with unusually low friction.

In this constitution, "AI-native" does not primarily mean that the server exposes AI features. It means the codebase, architecture, workflow, and extension model are designed from the first day for agent-driven development.

The core thesis:

> A server architecture can be made materially easier for AI agents to work on if the architecture is explicit, machine-readable, strongly bounded, testable, and generated from stable contracts.

## 0.1 Naming Standard

The project needs a name that is memorable to humans and immediately legible to people who already feel this problem but do not yet have a name for it.

The ratified product name is:

> vibit

Meaning:

- `vi`: a short, compact prefix that suggests visibility, vitality, and a clear interface.
- `bit`: the smallest familiar unit of digital systems. The name suggests small, composable, verifiable building blocks.

Recommended descriptive subtitle:

> vibit: an agent-native server framework for AI-maintainable backends.

Long-form positioning:

> vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.

Naming rules:

- The public name must be short, pronounceable, and memorable.
- The subtitle must immediately explain the category.
- The name must not imply that the server merely "has AI features".
- The name must make room for game servers first, but should not prevent backend-general adoption later.
- The repository description should contain the phrase `agent-native server framework`.
- Before finalizing the name, check major public registries and platforms for obvious conflicts.

Candidate names considered:

- `AgentNative`: clear but already too generic and publicly used.
- `AgentFirst`: clear but broad and already used.
- `AgentReady`: useful phrase, but already associated with repository assessment and readiness tools.
- `AgentFrame`: descriptive, but generic and easy to confuse with agent orchestration frameworks.
- `AgentForge`: memorable, but heavily used.
- `Framewright`: distinctive and connected to framework-building, but replaced by the simpler product name `vibit`.

Use `vibit` as the product name and `Agent-Native Server Framework` as the category name.

## 0.2 Documentation Language Standard

The canonical documentation language is English.

Every public-facing project document should have:

- An English source document
- A Simplified Chinese human-readable translation

English files are canonical. Chinese files are translations.

Recommended file naming:

```text
CONSTITUTION.md
CONSTITUTION.zh-CN.md
AGENTS.md
AGENTS.zh-CN.md
docs/<name>.md
docs/<name>.zh-CN.md
```

Rules:

- The English document is the source of truth for architecture, contracts, standards, and governance.
- The Chinese document should be maintained for human readers and contributors who think more clearly in Chinese.
- Translations should preserve meaning rather than translate word-for-word.
- When an English source document changes materially, the paired Chinese translation must be updated in the same change or explicitly marked as out of date.
- Machine-readable manifests and schemas should use English identifiers.
- Code identifiers, module names, commands, events, permissions, and errors should be English unless there is a strong domain reason not to.

This rule exists because the project should be globally legible to agents and tools while still serving the people who are shaping its early ideas.

## 1. Project Definition

This project is an Agent-Native Server Framework.

If focused on games, it is an Agent-Native Game Backend Framework.

Its purpose is to let human developers and AI agents safely evolve a long-lived server codebase through standardized modules, specs, contracts, tests, generated structure, and automated architecture checks.

The project may eventually provide AI gameplay features, such as NPC agents, memory, model routing, or tool-using characters. Those are valid extensions, but they are not the foundation.

The foundation is:

- Agent-friendly architecture
- Agent-readable project context
- Agent-safe change workflows
- Agent-verifiable behavior
- Agent-compatible standards for modules, APIs, events, data, tests, and documentation

## 2. External Standards Position

We should reference existing standards and methodologies, but we should not blindly inherit them.

Known references include:

- Superpowers: useful as a reference for structured agent workflows, skills, planning, test-first development, and verification discipline.
- AGENTS.md: useful as an emerging convention for repository-level instructions that coding agents can read.
- Spec-first development: useful for forcing requirements and contracts to precede implementation.
- DDD, CQRS, Event Sourcing, Actor Model, and Hexagonal Architecture: useful for explicit boundaries, state transitions, and dependency control.
- OpenAPI, AsyncAPI, JSON Schema, Protobuf, GraphQL, and database migration schemas: useful for contract-first implementation.
- Existing game/backend frameworks such as Skynet, Nakama, Colyseus, Agones, Pomelo, ET, and Zinx: useful as architectural references, but not as governing standards.

Our rule:

> We reference outside systems for vocabulary, workflow discipline, and proven patterns. We define our own standard where agent-native server maintainability requires stricter or different constraints.

The project should be compatible with `AGENTS.md` style instructions and may provide generated `AGENTS.md` files, but the constitution is higher-level than `AGENTS.md`.

The project may borrow the "skill" concept from Superpowers-like systems, but our canonical unit is broader: module contracts, change specs, architecture manifests, generated scaffolds, verification gates, and agent operating manuals.

## 3. Constitutional Principles

### 3.1 Architecture Must Be Explicit

No critical architecture rule may live only in a maintainer's memory.

Every module, boundary, dependency, contract, event, permission, error code, test expectation, and migration expectation should be represented in code, schema, manifest, or generated documentation.

### 3.2 Agent Context Is a First-Class Resource

The framework must reduce the amount of context an agent needs to read before making a correct change.

A good module should tell an agent:

- What it owns
- What it does not own
- Which APIs it exposes
- Which events it publishes
- Which events it consumes
- Which dependencies are allowed
- Which invariants must not be broken
- Which tests verify the module
- Which files are extension points
- Which files are generated and should not be hand-edited

### 3.3 No Unbounded Edits

An agent should rarely need to perform unconstrained whole-repository edits.

Every meaningful change should be expressed as a bounded change request with affected modules, expected contracts, generated files, tests, migrations, and verification steps.

### 3.4 Schema Before Code

Public behavior must be specified before implementation.

The preferred order is:

1. Requirement
2. Spec
3. Contract
4. Generated structure
5. Business logic
6. Tests
7. Verification
8. Documentation update

An agent should not invent ad hoc APIs, event payloads, database shapes, or error formats without first updating the relevant schema or manifest.

### 3.5 Generated Shape, Handwritten Logic

The framework should generate repeatable structure.

Agents and humans should focus on behavior, invariants, edge cases, and tests, not on guessing where files go or how boilerplate should look.

### 3.6 Strong Module Boundaries

Modules are bounded contexts.

A module may own data, commands, queries, events, policies, repositories, and tests. Other modules may not reach into its internals directly.

Cross-module communication must happen through approved interfaces:

- Commands
- Queries
- Events
- Public module APIs
- Generated clients

### 3.7 Commands Express Intent; Events Express Facts

Commands describe requested changes.

Events describe facts that have happened.

This distinction is mandatory because it helps agents reason about business flow and makes behavior easier to test, replay, and verify.

### 3.8 The Server Remains Authoritative

Agent-written code must not weaken server authority.

All state-changing actions must pass validation, permission checks, invariants, and transactional consistency rules.

This principle applies both to human users and AI agents.

### 3.9 Architecture Is Testable

Architectural rules must be automatically checkable.

If a rule matters, the framework should eventually provide a command that can verify it.

Examples:

- Modules must not import forbidden modules.
- Events must be versioned.
- Commands must have schemas.
- Public APIs must have contract tests.
- Migrations must have rollback or compatibility tests.
- Permissioned operations must declare permissions.
- Generated files must not be hand-edited.

### 3.10 Documentation Must Be Operational

Documentation should help an agent or human take the next correct action.

The project should avoid decorative documentation that sounds complete but does not guide implementation, testing, or verification.

Every module should eventually have an agent-readable operating guide.

## 4. Required Project Artifacts

The project should evolve toward the following artifact system.

### 4.1 Constitution

File:

```text
CONSTITUTION.md
```

Purpose:

- Defines project principles
- Defines what counts as valid architecture
- Defines how standards are adopted or amended
- Provides the top-level decision framework

### 4.2 Repository Agent Guide

Expected file:

```text
AGENTS.md
```

Purpose:

- Gives coding agents concise operational instructions
- Lists build, test, lint, generation, and verification commands
- Defines "always", "ask first", and "never" rules
- Points agents to module-level guides and architecture manifests

`AGENTS.md` should be generated or reviewed against this constitution.

### 4.3 Architecture Manifests

Expected directory:

```text
.arch/
```

Expected files may include:

```text
.arch/modules.yaml
.arch/dependencies.yaml
.arch/conventions.yaml
.arch/commands.yaml
.arch/events.yaml
.arch/errors.yaml
.arch/permissions.yaml
.arch/test-matrix.yaml
.arch/generation.yaml
```

Purpose:

- Gives agents a machine-readable map of the system
- Enables architecture checks
- Enables code generation
- Enables impact analysis

### 4.4 Module Manifests

Expected file per module:

```text
modules/<module>/module.yaml
```

Minimum expected fields:

```yaml
module: inventory
type: domain
owns:
  entities:
    - inventory
    - inventory_item
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
    - economy
  forbidden:
    - matchmaking
invariants:
  - item_count_must_not_be_negative
  - inventory_capacity_must_not_exceed_limit
tests:
  required:
    - command_tests
    - event_tests
    - invariant_tests
    - migration_tests
```

### 4.5 Module Agent Guides

Expected file per module:

```text
modules/<module>/AGENTS.md
```

Purpose:

- Explains when to use the module
- Explains when not to use the module
- Lists extension points
- Lists forbidden shortcuts
- Lists required tests
- Links to schemas, commands, events, and invariants

### 4.6 Change Specs

Expected directory:

```text
changes/
```

Each non-trivial change should have a directory:

```text
changes/<date>-<change-id>/
  request.md
  spec.yaml
  impact.md
  plan.md
  checklist.md
  verification.md
```

Purpose:

- Gives agents durable context
- Makes decisions traceable
- Supports handoff between agents and humans
- Makes later maintenance easier

Example `spec.yaml`:

```yaml
change_id: add-season-pass
type: feature
affected_modules:
  - player
  - economy
  - reward
new_commands:
  - PurchaseSeasonPass
  - ClaimSeasonPassReward
new_events:
  - SeasonPassPurchased
  - SeasonPassRewardClaimed
data_migrations:
  required: true
compatibility:
  breaking_api: false
  breaking_db: false
acceptance_tests:
  - player_can_purchase_pass
  - player_cannot_claim_locked_reward
  - duplicate_claim_is_rejected
```

### 4.7 Conversation Logs

Expected directory:

```text
conversations/
```

Purpose:

- Preserves maintainer-agent project memory
- Records maintainer narrative with high fidelity
- Summarizes agent responses when exact wording is not required
- Explains how decisions and standards emerged
- Links conversations to change specs and artifacts
- Redacts secrets, tokens, credentials, and unrelated private data before commit

Conversation logs are not a replacement for change specs. A conversation log explains why the project moved in a direction. A change spec explains how a specific change is executed.

Maintainer statements may be preserved in their original language because original wording can be part of product intent.

## 5. Standard Server Shape

The framework should converge toward a predictable structure.

Illustrative structure:

```text
modules/
  player/
    module.yaml
    AGENTS.md
    commands/
    queries/
    events/
    models/
    repositories/
    policies/
    tests/
  inventory/
  economy/
  match/
  chat/
  notification/
schema/
  api/
  commands/
  events/
  errors/
  permissions/
  db/
.arch/
changes/
conversations/
tools/
docs/
```

The exact implementation may change, but the framework should preserve predictability.

## 6. Standard Change Workflow

Every non-trivial feature, bug fix, migration, or refactor should follow this workflow.

### 6.1 Intake

Convert the user request into a clear requirement.

The requirement should identify:

- Desired behavior
- User-visible outcome
- Affected modules
- Unknowns
- Risks
- Acceptance criteria

### 6.2 Impact Analysis

Before code changes, identify:

- Modules affected
- Contracts affected
- Data affected
- Events affected
- Permissions affected
- Tests affected
- Migration requirements
- Compatibility risks

### 6.3 Plan

The plan should be short, bounded, and executable.

It should identify:

- Files to create or edit
- Generated artifacts
- Manual logic
- Required tests
- Verification commands

### 6.4 Contract Update

Update schemas and manifests before implementation where behavior changes public or cross-module contracts.

### 6.5 Code Generation

Use generators for repeatable structure whenever available.

Manual creation of framework-shaped files should be treated as a temporary fallback.

### 6.6 Implementation

Implement only within the defined change boundary.

Do not bypass module public APIs, validators, repositories, migrations, or permission systems.

### 6.7 Verification

A change is not complete until verification is run or explicitly documented as not run.

Verification should include the narrowest relevant tests and any architecture checks.

### 6.8 Documentation Update

If a change modifies architecture, module ownership, public behavior, commands, events, schemas, permissions, or test procedure, update the relevant docs or manifests.

## 7. Module Standard

Each module should have a stable internal shape.

Recommended module directories:

```text
commands/
queries/
events/
models/
repositories/
policies/
services/
tests/
```

Definitions:

- `commands/`: state-changing requests and handlers
- `queries/`: read-only access patterns
- `events/`: facts emitted by the module
- `models/`: domain models and value objects
- `repositories/`: persistence boundaries
- `policies/`: event reactions and cross-module orchestration
- `services/`: local domain services, not generic dumping grounds
- `tests/`: module-local tests

Rules:

- A module owns its internal data.
- Other modules cannot import internal files directly.
- Public behavior must be exposed through declared APIs.
- Events must be versioned once public.
- Commands must validate input and permissions.
- Invariants must be tested.

## 8. Contract Standard

The framework should treat contracts as first-class.

Contract types:

- API contracts
- Command contracts
- Query contracts
- Event contracts
- Error contracts
- Permission contracts
- Database migration contracts
- Generated client contracts

Rules:

- Public contracts must be declared in schema.
- Contracts must be versioned when compatibility matters.
- Breaking changes must be explicit.
- Generated code must be traceable to its source schema.
- Agents must not hand-edit generated contract output unless the generator itself is the target.

## 9. Testing Standard

Testing is not optional process decoration. Testing is an architectural guardrail for agent-driven development.

Required test categories should include:

- Unit tests for handlers, policies, validators, and domain logic
- Contract tests for public APIs, commands, queries, and events
- Invariant tests for business rules
- Integration tests for module collaboration
- Migration tests for database changes
- Replay tests where event history is used
- Architecture tests for dependency and manifest rules

Every module should declare its required test categories in `module.yaml`.

## 10. Architecture Verification Standard

The framework should eventually provide commands similar to:

```bash
server check architecture
server check contracts
server check module <module>
server check change <change-id>
server generate module <module>
server generate command <module> <command>
server generate event <module> <event>
```

These commands should be designed for both humans and agents.

Their output should be concise, deterministic, and actionable.

## 11. Agent Operating Rules

Agents working in this repository should follow these rules.

### 11.1 Always

- Read the relevant constitution, architecture manifest, module manifest, and module `AGENTS.md` before changing a module.
- Prefer schema and manifest updates before implementation.
- Keep edits bounded to the declared change.
- Add or update tests for behavior changes.
- Run relevant verification commands when available.
- Document verification results.

### 11.2 Ask First

Ask before:

- Changing constitutional principles
- Redefining module ownership
- Introducing a new architectural pattern
- Making breaking API or event changes
- Changing generated file conventions
- Removing tests or weakening validation
- Moving data ownership between modules

### 11.3 Never

Never:

- Bypass module boundaries for convenience
- Hide business logic in transport handlers
- Add unregistered public events
- Add unregistered permissions
- Add untyped cross-module payloads
- Hand-edit generated files without marking the reason
- Claim a change is verified when verification was not run

## 12. AI Feature Boundary

The framework may support AI gameplay and AI product features, but those must not be confused with agent-native maintainability.

AI gameplay features may include:

- NPC agents
- Agent memory
- Model routing
- Tool calling
- Agent simulation
- Narrative systems

These are application-layer capabilities.

The constitutional foundation remains:

- Explicit architecture
- Machine-readable standards
- Bounded changes
- Schema-first contracts
- Generated structure
- Automated verification
- Agent-readable context

## 13. Governance

This constitution can change, but changes must be deliberate.

Any constitutional amendment should include:

- Motivation
- Problem being solved
- Alternatives considered
- Compatibility impact
- Required migration
- Updated examples

Constitutional changes should be rare compared with implementation changes.

## 14. Decision Rule

When choosing between two designs, prefer the design that:

1. Gives agents less ambiguous context
2. Creates stronger module boundaries
3. Makes behavior easier to verify
4. Makes contracts more explicit
5. Reduces hidden coupling
6. Supports code generation
7. Supports long-term maintenance
8. Still remains practical for human developers

Agent-native does not mean agent-only.

The best architecture should make both humans and agents more effective.

## 15. Open Questions

These questions remain intentionally open:

- Which language should be the reference implementation?
- Should the first implementation be game-specific or backend-general?
- Should the framework own transport, or integrate with existing transports?
- Should event sourcing be mandatory or optional?
- Should module manifests be YAML, JSON, TOML, or code-native?
- Should generators be built into the framework CLI from day one?
- Which database should be the first-class reference target?
- How much should be enforced by convention versus static analysis?
- How should multi-agent collaboration be coordinated in large changes?
- How should the framework measure whether it is actually easier for agents to maintain?

These questions should be answered through design notes, prototypes, and concrete implementation experience.

## 16. First Implementation Target

The first serious prototype should prove this claim:

> Given a new backend requirement, an AI coding agent can identify the affected module, update the correct contracts, generate the correct structure, implement the behavior, add tests, run verification, and update documentation without damaging unrelated architecture.

The prototype should be small but complete.

Suggested sample domain:

- Player accounts
- Inventory
- Currency
- Rewards
- Tasks or quests
- Match sessions

The demo should emphasize maintainability and agent workflow over feature count.

## 17. Founding Statement

This project is not trying to make agents magically smarter.

It is trying to make server architecture more legible, bounded, verifiable, and evolvable so that agents can do high-quality work reliably.

The standard is simple:

> If an agent cannot understand where a change belongs, what contracts it affects, what tests prove it, and what architecture rules constrain it, then the framework has not yet done its job.
