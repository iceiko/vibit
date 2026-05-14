# Conversation: Authentication Generated Contract Shape Timing

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-decide-authentication-generated-contract-shape-timing/`

Related artifacts:

- `docs/authentication-generated-contract-shape-timing.md`
- `decisions/ADR-0038-authentication-generated-contract-shape-timing.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue advancing the planned work queue unless a real confirmation point appeared. Earlier context established that routine technical sequencing should be handled by the agent's professional judgment, while product or architecture branch decisions should be surfaced.

The next ready work item was `W-0088 Decide authentication generated contract shape timing`.

## Maintainer Narrative

The maintainer expects the project to stay controlled and self-bootstrapping. Around authentication, the maintainer wants necessary preparation before implementation rather than rushing directly into behavior.

## Agent Response Summary

The agent evaluated the existing semantic authentication contracts, generated-output standard, runtime authentication implementation boundary, and current generated contract shape tooling.

The agent concluded that generated Go authentication contract shapes should be introduced before application authentication service interfaces and before runtime authentication behavior, but that `W-0088` should only record the timing, source, output path, immutability, and check requirements.

## Decisions

- Generated authentication contract shapes should come after semantic contracts and the implementation boundary.
- Generated authentication contract shapes should come before service interfaces and runtime authentication behavior.
- Runtime contract families need a family segment in generated output paths.
- The planned authentication output root is `runtime/internal/generated/contracts/runtime/authentication/`.
- `W-0088` does not generate files.
- A later work item may explicitly authorize generator/check support and generated authentication shape output.
- Token generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, dependencies, repository changes, and migration changes remain deferred.

## Artifacts

- `docs/authentication-generated-contract-shape-timing.md`
- `docs/authentication-generated-contract-shape-timing.zh-CN.md`
- `decisions/ADR-0038-authentication-generated-contract-shape-timing.md`
- `changes/2026-05-14-decide-authentication-generated-contract-shape-timing/`
- `.arch/work-items.yaml`

## Open Questions

- The exact generator implementation details remain for the next bounded work item.
- Whether runtime session contract shapes should later use the same family-aware layout is not decided here.

## Follow-Up

- Advance `W-0089`.
- Extend generator/check support before committing generated authentication contract shape files.
- Keep runtime authentication behavior behind later gates.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
