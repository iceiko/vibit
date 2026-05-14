# Conversation: Runtime Authentication Contract Shape Generation

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-add-runtime-authentication-contract-shape-generator-support-and-output/`

Related artifacts:

- `tools/vibit`
- `runtime/internal/generated/contracts/runtime/authentication/`
- `.arch/work-items.yaml`
- `docs/generated-output.md`
- `docs/authentication-generated-contract-shape-timing.md`

## Context

The maintainer asked the agent to continue through the queued work items unless a real confirmation point appeared. `W-0088` had already decided that generated Go authentication contract shapes should exist before application service interfaces and runtime authentication behavior.

The next ready work item was `W-0089 Add runtime authentication contract shape generator support and output`.

## Maintainer Narrative

The maintainer continues to prioritize a controlled, self-bootstrapping architecture over rushing runtime behavior. Routine technical sequencing can proceed by agent judgment, while product or architecture branch decisions should be surfaced.

## Agent Response Summary

The agent extended the contract shape generator and generated-output checks for the runtime authentication family-aware output path. It generated metadata-only Go authentication contract shape files from semantic authentication contract sources and updated runtime boundary checks so generated metadata is allowed while runtime authentication behavior remains deferred.

## Decisions

- Runtime authentication generated contract shapes are now committed as metadata-only generated output.
- The accepted output root is `runtime/internal/generated/contracts/runtime/authentication/`.
- The generated package name is `runtimeauthenticationcontracts`.
- Runtime authentication generated files must include source trace, generator trace, contract trace, and `Family: authentication`.
- Generated authentication shapes do not authorize application service implementations, token generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, dependencies, repository changes, or migration schema changes.
- `M-017` is completed and `M-018 Application Authentication Service Interface Boundary` is opened.

## Artifacts

- `runtime/internal/generated/contracts/runtime/authentication/**`
- `tools/vibit`
- `changes/2026-05-14-add-runtime-authentication-contract-shape-generator-support-and-output/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `docs/authentication-generated-contract-shape-timing.md`
- `docs/authentication-generated-contract-shape-timing.zh-CN.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/selected-login-token-boundary-checks.zh-CN.md`

## Open Questions

- The exact application authentication service interface boundary remains for `W-0090`.
- Runtime authentication behavior remains for later gated work.
- Whether other runtime families, such as runtime session contracts, should use the same generated path shape remains a future decision.

## Follow-Up

- Advance `W-0090`.
- Define the application-owned service interface boundary before implementing authentication behavior.
- Keep generated output immutable and source-driven.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
