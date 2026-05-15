# Conversation: Verifier Digest Comparison Helper Implementation

Date: 2026-05-15  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-15-implement-verifier-digest-comparison-helpers/`

Related artifacts:

- `runtime/internal/app/authentication/verifier_comparison.go`
- `runtime/internal/app/authentication/verifier_comparison_test.go`
- `docs/verifier-digest-comparison-helper-gate.md`
- `decisions/ADR-0049-verifier-digest-comparison-helper-gate.md`
- `.arch/work-items.yaml`

## Context

The verifier digest comparison helper gate was completed in `W-0104`. The next ready work item was `W-0105`, implementing helper-only constant-time verifier digest comparison under the already-defined boundary.

## Maintainer Narrative

The maintainer asked the agent to continue development and had previously emphasized that ordinary technical follow-through should proceed without stopping for confirmation, while true architectural or product decisions should be raised explicitly.

The maintainer expects the project to stay self-bootstrapping and controllable through bounded work items, architecture manifests, change specs, checks, and conversation memory.

## Agent Response Summary

The agent advanced `W-0105` as a helper-only implementation change.

The implementation adds application-owned verifier digest comparison helpers under `runtime/internal/app/authentication`. The helpers compare computed credential or token verifier digest bytes to stored verifier digest bytes with `crypto/hmac.Equal`, return a redacted match or mismatch result, and reject lookup digest classes, wrong verifier classes, missing input, and malformed input with redacted typed errors.

The change updates repository checks so this exact helper file is authorized while other files remain blocked from verifier comparison and authentication service behavior.

The change does not add authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf authentication messages, WebSocket proof carriers, repository changes, migration changes, startup wiring, authentication dependencies, or production authentication behavior.

## Decisions

- Implement verifier digest comparison now because `W-0105` authorizes the helper-only slice.
- Use `crypto/hmac.Equal` as the preferred constant-time primitive from `ADR-0049`.
- Treat mismatch as a redacted non-match result instead of public authentication behavior.
- Keep the next ready work item conservative: define the authentication service behavior implementation gate before adding login or token validation behavior.

## Artifacts

- `runtime/internal/app/authentication/verifier_comparison.go`
- `runtime/internal/app/authentication/verifier_comparison_test.go`
- `changes/2026-05-15-implement-verifier-digest-comparison-helpers/`
- `conversations/2026-05-15-verifier-digest-comparison-helper-implementation.md`

## Open Questions

- None for this helper-only implementation.

## Follow-Up

- Define the authentication service behavior implementation gate before wiring login, token validation, protocol carriers, repositories, migrations, startup wiring, dependencies, or production authentication behavior.

## Redaction Notes

No secrets, tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, account details, or private data are recorded in this conversation log.
