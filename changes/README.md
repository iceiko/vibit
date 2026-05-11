# Change Specs

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: Durable context for non-trivial changes

This directory stores change specs for non-trivial work.

Use `docs/change-spec.md` as the canonical standard.

Template files live in:

```text
changes/_template/
```

To start a new change:

1. Copy `changes/_template/` to `changes/YYYY-MM-DD-short-change-id/`.
2. Fill in `request.md`.
3. Fill in `spec.yaml`.
4. Complete impact analysis and plan before implementation.
5. Keep `checklist.md` and `verification.md` current as work proceeds.

Small typo fixes and narrow documentation edits may not need a full change spec, but agents must still record verification status in their final handoff.
