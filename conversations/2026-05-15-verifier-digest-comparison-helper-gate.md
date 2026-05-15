# Conversation: Verifier Digest Comparison Helper Gate

Date: 2026-05-15  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-15-define-verifier-digest-comparison-helper-gate/`

Related artifacts:

- `docs/verifier-digest-comparison-helper-gate.md`
- `docs/verifier-digest-comparison-helper-gate.zh-CN.md`
- `decisions/ADR-0049-verifier-digest-comparison-helper-gate.md`
- `.arch/work-items.yaml`

## Context

The verifier digest computation helper implementation was completed in `W-0103`. The next ready work item was `W-0104`, defining the verifier digest comparison helper gate before adding comparison code.

## Maintainer Narrative

The maintainer asked the agent to continue autonomous development under the established workflow. Under the project continuation standard, "continue" means advancing one `next_ready` work item unless a real decision boundary requires maintainer confirmation.

The maintainer expects the project to remain self-bootstrapping, controllable, and agent-native, with small bounded steps that prevent future agents from mixing responsibilities.

## Agent Response Summary

The agent advanced `W-0104` as a gate-only change.

The gate defines future comparison helper ownership under `runtime/internal/app/authentication`, names `verifier_comparison.go` and `verifier_comparison_test.go`, keeps `verifier_digest.go` computation-only, requires `crypto/hmac.Equal` as the preferred constant-time primitive, allows `crypto/subtle.ConstantTimeCompare` only as an equivalent standard-library alternative, and forbids raw material, lookup digest, key id, database-only, and protocol metadata comparison.

The gate defines redacted failure posture and future focused tests. It does not add verifier comparison code, authentication service behavior, protocol carriers, repository changes, migrations, startup wiring, new dependencies, or production authentication behavior.

## Decisions

- Use `ADR-0049` for the verifier digest comparison helper gate.
- Keep comparison helpers in new files separate from digest computation helpers.
- Future comparison helpers must compare only computed verifier digest bytes to stored verifier digest bytes.
- Future comparison helpers must use a constant-time comparison primitive.
- The next work item after this gate is `W-0105`, the narrow verifier digest comparison helper implementation slice.

## Artifacts

- `docs/verifier-digest-comparison-helper-gate.md`
- `docs/verifier-digest-comparison-helper-gate.zh-CN.md`
- `decisions/ADR-0049-verifier-digest-comparison-helper-gate.md`
- `changes/2026-05-15-define-verifier-digest-comparison-helper-gate/`
- `conversations/2026-05-15-verifier-digest-comparison-helper-gate.md`

## Open Questions

- None for this gate-only change.

## Follow-Up

- Implement `W-0105` as a helper-only verifier digest comparison slice.
- Keep authentication service behavior, protocol carriers, repositories, migrations, startup wiring, dependencies, and production authentication behavior behind later bounded work items.

## Redaction Notes

No secrets, tokens, generated credentials, account details, or private data are recorded in this conversation log.
