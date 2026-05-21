# Plan

1. Create the local onboarding/device credential issuance gate standard in English and Simplified Chinese.
2. Record `ADR-0089`.
3. Add conversation memory for the continuation step.
4. Add the complete change spec directory for `W-0181`.
5. Mark `M-109/W-0181` completed and create `M-110/W-0182` as the next ready implementation work item.
6. Update `.arch/runtime.yaml`, `.arch/conventions.yaml`, `.arch/contracts.yaml`, `.arch/reference.yaml`, module manifests, AGENTS guides, README, and alpha goal pointers.
7. Add `runtime.local_onboarding_device_credential_issuance_gate` to `tools/vibit` and `rules/check-rules.json`.
8. Verify the change spec, work queue, runtime checks, memory checks, schema checks, full repository checks, and diff whitespace.

## Generated Artifacts

None.

## Handwritten Runtime Logic

None. This is a gate-only change.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-local-onboarding-device-credential-issuance-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

No data migration or runtime rollback is required because no runtime behavior, Protobuf source, generated output, migration, or dependency is added.
