# Conversation: Prototype-Ready Local Development Path Package

Date: 2026-05-22
Participants: Maintainer, Agent
Related work item: `W-0200`
Related decision: `ADR-0108`

## Context

The maintainer asked to continue advancing vibit toward a real product and production-useful stage. The agent advanced the next-ready work item, `W-0200 Implement prototype-ready local development path package`.

The work implemented the package inside the `ADR-0107` gate using docs, examples, placeholder configuration, `.gitignore` guardrails, static checks, and continuation metadata. It did not change production runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, hosted deployments, release artifacts, public announcements, paid promotion, authentication/session semantics, broad operations/admin behavior, broad product modules, or direct Nakama/Pitaya API compatibility.

## Maintainer Narrative

The maintainer wants vibit to keep moving beyond a pre-alpha skeleton toward a product that can eventually serve as a real game/backend server development foundation. The immediate instruction was to continue progressing the next bounded work item while preserving the repository's ask-first boundaries.

## Agent Response Summary

The agent packaged the local prototype-ready developer path as a source-first, docs-and-checks slice. The package makes the current local alpha proof easier to find and run, documents the still-manual local setup boundaries, adds placeholder-only local environment guidance, and adds a static rule so future changes preserve the package.

## Decisions

- The prototype-ready local development path package is recorded in `docs/prototype-ready-local-development-path-package.md`.
- The paired Simplified Chinese translation is `docs/prototype-ready-local-development-path-package.zh-CN.md`.
- The package includes `examples/README.md`, `examples/README.zh-CN.md`, and `examples/local.prototype.env.example`.
- Private local env files such as `.vibit.local.env`, `.env.local`, and `.env.*.local` are ignored by `.gitignore`.
- The repository check rule is `runtime.prototype_ready_local_development_path_package`.
- The next bounded direction is `W-0201 Define storage objects behavior gate`.

## Artifacts

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/prototype-ready-local-development-path-package.md`
- `docs/prototype-ready-local-development-path-package.zh-CN.md`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `examples/local.prototype.env.example`
- `decisions/ADR-0108-prototype-ready-local-development-path-package.md`
- `changes/2026-05-22-implement-prototype-ready-local-development-path-package/`
- `rules/check-rules.json`
- `tools/vibit`
- `.gitignore`

## Open Questions

No new product questions were opened by this slice. The local prototype remains intentionally source-first and manual until a later bounded work item authorizes runtime, deployment, packaging, or hosted behavior.

## Follow-Up

The current state remains source-first and local. A developer can clone the repository, run checks, run Go tests, run the redacted authenticated gameplay request-loop proof, inspect local status surfaces, and understand which PostgreSQL and secret setup steps remain manual.

The next product capability family is storage objects beyond the inventory proof slice, but only through a behavior gate first.

## Redaction Notes

Future work must ask first before adding runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependencies, automatic startup migration behavior, hosted deployments, release artifacts, public announcements beyond the GitHub release record, paid promotion, authentication/session semantic changes, broad operations/admin behavior, broad product modules, or direct Nakama/Pitaya API compatibility.

No private local environment file was read or printed. The new example environment file uses placeholders only and contains no real credentials, raw access tokens, verifier keys, DSNs, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.
