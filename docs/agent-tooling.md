# Agent Tooling Standard

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Agent-facing inspection, generation, and verification commands

This standard defines the minimum machine-readable tooling surface that lets agents work on vibit without relying on memory or broad source reading.

The paired Simplified Chinese translation is `docs/agent-tooling.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit needs tools that answer narrow questions with structured output.

Agents should be able to inspect the next work item, registered contracts, generated outputs, reference planning context, and verification rules before editing code. This reduces accidental architecture drift and makes continuation predictable.

## 2. Required Commands

The current agent-facing commands are:

```bash
node tools/vibit inspect work --json
node tools/vibit inspect next --json
node tools/vibit inspect contracts --json
node tools/vibit inspect contracts --module inventory --json
node tools/vibit inspect contracts --type command --json
node tools/vibit inspect contracts --status draft --json
node tools/vibit inspect contract --module inventory --type command --id GrantItem
node tools/vibit inspect generated --json
node tools/vibit inspect generated --module inventory --json
node tools/vibit inspect generated --type command --json
node tools/vibit inspect reference --json
node tools/vibit check agent-tooling --json
node tools/vibit generate contract-shapes all
```

These commands are intentionally small. They are not a replacement for reading governing standards before substantial changes.

## 3. Inspection Rules

Inspection commands should:

- Return JSON by default when `--json` is provided.
- Include `schema_version`, `kind`, and `project`.
- Report the source artifact path for every durable fact.
- Normalize repository-relative paths in JSON fields such as `artifact`, `path`, `source`, and `output` to forward-slash form, even on Windows.
- Prefer explicit `exists`, `status`, and `summary` fields over prose.
- Avoid hiding ask-first boundaries.

## 4. Generation Rules

Generation commands should:

- Read source contracts, manifests, or schemas.
- Write only declared generated roots.
- Include generated markers, source traces, and generator traces.
- Be reproducible enough for checks to detect drift.
- Avoid adding runtime behavior unless a separate work item owns that implementation.

`node tools/vibit generate contract-shapes all` generates inspectable Go contract shape files from semantic contract manifests. These files are generated artifacts only. They do not implement handlers, persistence, authentication, routes, or protocol envelope behavior.

## 5. Verification

Current verification:

```bash
node tools/vibit check agent-tooling --json
node tools/vibit check generated --json
node tools/vibit check all --json
```

`check agent-tooling` verifies that this standard and its Simplified Chinese translation exist, that the command surface is documented, and that public docs keep translation pairs.

## 6. Agent Rules

Agents must:

- Run `node tools/vibit inspect next --json` before interpreting a continuation request when work state is unclear.
- Run `node tools/vibit inspect contracts --json` before broad contract or generator work.
- Use `node tools/vibit inspect contracts --module <module> --json`, `node tools/vibit inspect contracts --type <contract-type> --json`, or `node tools/vibit inspect contracts --status <status> --json` when only a narrow contract slice is relevant.
- Run `node tools/vibit inspect generated --json` before editing generated output or generator behavior.
- Use `node tools/vibit inspect generated --module <module> --json` or `node tools/vibit inspect generated --type <contract-type> --json` when only a narrow generated-output slice is relevant.
- Run `node tools/vibit inspect reference --json` before planning new game server capability families.
- Regenerate generated contract shapes through `node tools/vibit generate contract-shapes all` instead of hand-editing those files.

Agents must not:

- Treat generated contract shapes as handwritten runtime implementation.
- Use inspection output as permission to cross ask-first boundaries.
- Replace narrow tooling with broad prose-only standards when a check or JSON inspection can enforce the rule.
