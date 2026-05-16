# Plan

## Steps

1. Read the active work item, related authentication boundaries, helper gates, and current runtime manifests.
2. Add the English authentication service behavior implementation gate standard.
3. Add the Simplified Chinese translation.
4. Add `ADR-0050`.
5. Add a conversation log preserving maintainer intent and redaction.
6. Update architecture manifests and agent guides with the new gate.
7. Add a runtime check rule for the gate.
8. Mark `W-0106` completed and open the next conservative work item.
9. Run repository verification.

## Boundaries

This change is documentation, manifest, and tooling only.

Do not add:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- login execution
- token validation
- logout, refresh, or cleanup behavior
- protocol carriers
- repository changes
- migrations
- startup wiring
- authentication dependencies
- production authentication behavior

## Verification Plan

Run:

- `node -c tools/vibit`
- `cd runtime && go test ./...`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-authentication-service-behavior-implementation-gate --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan excluding `.git` and `.vibit.local.env`
