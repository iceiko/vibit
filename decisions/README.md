# Agent Decision Records

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: Durable decision rationale for vibit

This directory stores Agent Decision Records.

Use `docs/agent-decision-record.md` as the canonical standard.

Rules:

- Record durable decisions, not every small implementation detail.
- Keep rationale public, concise, and inspectable.
- Do not store private chain-of-thought.
- Link decisions from modules, change specs, and conversation logs when relevant.
- Use decision records to explain generated-file overrides and long-term architecture choices.

Template:

```text
decisions/_template/adr-agent.md
```
