# Plan

## Files To Create

- `runtime/internal/generated/contracts/runtime/authentication/commands/AuthenticateWithDeviceCredential.go`
- `runtime/internal/generated/contracts/runtime/authentication/commands/LogoutAccessToken.go`
- `runtime/internal/generated/contracts/runtime/authentication/commands/RefreshAccessToken.go`
- `runtime/internal/generated/contracts/runtime/authentication/commands/ValidateAccessToken.go`
- `runtime/internal/generated/contracts/runtime/authentication/errors/authentication_errors.go`
- `runtime/internal/generated/contracts/runtime/authentication/events/AuthenticationFailed.go`
- `runtime/internal/generated/contracts/runtime/authentication/events/AuthenticationSucceeded.go`
- `runtime/internal/generated/contracts/runtime/authentication/events/LogoutRequested.go`
- `runtime/internal/generated/contracts/runtime/authentication/events/TokenIssued.go`
- `runtime/internal/generated/contracts/runtime/authentication/events/TokenRevoked.go`
- `runtime/internal/generated/contracts/runtime/authentication/events/TokenValidationFailed.go`
- `runtime/internal/generated/contracts/runtime/authentication/permissions/authentication_permissions.go`
- `conversations/2026-05-14-runtime-authentication-contract-shape-generation.md`

## Files To Edit

- `tools/vibit`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `docs/authentication-generated-contract-shape-timing.md`
- `docs/authentication-generated-contract-shape-timing.zh-CN.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/selected-login-token-boundary-checks.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/module.yaml`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`

## Generated Artifacts

Generated through:

```bash
node tools/vibit generate contract-shapes all
```

The runtime authentication files are metadata-only and immutable to non-system agents.

## Handwritten Logic

The only handwritten logic change is in `tools/vibit`, adding runtime-family-aware generated contract shape generation, inspection, drift checks, and boundary validation.

No Go runtime behavior is added.

## Tests

No Go tests are required because this change does not add runtime behavior.

Repository checks verify generator syntax, generated output, runtime guards, contracts, module state, work state, and the change spec.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit generate contract-shapes all`
- `node tools/vibit inspect generated --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-runtime-authentication-contract-shape-generator-support-and-output --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback by removing the generated runtime authentication shape files, reverting runtime-family-aware generator/check support, and returning M-017/W-0089 to the pre-generation state.
