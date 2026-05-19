# Plan

1. Add `runtime/migrations/postgres/000005_create_runtime_sessions.sql`.
2. Include goose Up and Down markers plus `runtime.session` ownership.
3. Add only `runtime_sessions`; do not add connection registry tables.
4. Include required lifecycle columns and optional revocation/access-token record linkage columns.
5. Add constraints and indexes for lifecycle inspection.
6. Update manifests, module metadata, AGENTS guides, rules, tools, ADR, and conversation memory.
7. Verify migration, runtime, module, work, memory, change, all checks, and diff whitespace.

## Rollback

Reversal means removing the migration source, reverting `ADR-0060`, and reopening `M-059/W-0131` before any repository or adapter work depends on `runtime_sessions`.
