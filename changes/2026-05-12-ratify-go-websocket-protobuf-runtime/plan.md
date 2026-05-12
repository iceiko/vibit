# Plan

## Files To Change

- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `contracts/inventory/**`
- `decisions/ADR-0003-first-reference-runtime-language.md`
- `decisions/ADR-0004-minimal-server-instance-model.md`
- `decisions/ADR-0007-yaml-contract-source-format.md`
- `decisions/ADR-0008-go-server-runtime-language.md`
- `decisions/ADR-0009-websocket-protobuf-client-protocol.md`
- `decisions/ADR-0010-foundational-dependency-adoption.md`
- `rules/check-rules.json`
- `tools/vibit`
- `conversations/2026-05-12-go-websocket-protobuf-direction.md`

## Files To Remove

- `package.json`
- `package-lock.json`
- `tsconfig.json`
- `tools/package.json`
- `modules/inventory/commands/GrantItem.ts`
- `modules/inventory/queries/GetInventory.ts`
- `modules/inventory/repositories/InMemoryInventoryRepository.ts`
- `modules/inventory/policies/inventoryCapacityPolicy.ts`
- `modules/inventory/policies/inventoryPermissionPolicy.ts`
- `modules/inventory/tests/GrantItem.test.ts`
- `modules/inventory/tests/GetInventory.test.ts`
- `modules/inventory/generated/contracts/GrantItem.generated.ts`
- `modules/inventory/generated/contracts/GetInventory.generated.ts`

## Steps

1. Mark `ADR-0003` superseded.
2. Add ADRs for Go runtime, WebSocket/Protobuf protocol, and foundational dependency adoption.
3. Update architecture manifests to reflect Go/WebSocket/Protobuf.
4. Update contract metadata to distinguish semantic source from wire schema.
5. Update docs and Simplified Chinese translations.
6. Remove misleading TypeScript runtime and package baseline artifacts.
7. Update `tools/vibit check runtime` so it passes as not applicable before Go runtime implementation exists.
8. Record conversation memory and verification.

## Verification Plan

Run:

```bash
node tools/vibit check architecture
node tools/vibit check schemas
node tools/vibit check memory
node tools/vibit check contracts
node tools/vibit check generated
node tools/vibit check runtime
node tools/vibit check all --json
rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .
git diff --check
```

## Rollback

If this direction is rejected, restore the removed TypeScript files from git history and supersede `ADR-0008` and `ADR-0009` with the approved alternative. Do not silently change this ADR history.
