# Checklist

## Implementation

- [x] Added local onboarding request/result vocabulary.
- [x] Added explicit device credential entropy and id generator dependencies.
- [x] Generated device credential material with `GenerateDeviceCredentialMaterial`.
- [x] Computed credential lookup and verifier digests with existing helpers.
- [x] Created player account and credential record in one unit of work.
- [x] Stored only credential digests and metadata.
- [x] Returned raw credential text only after unit-of-work success.
- [x] Preserved existing login route non-creation behavior.
- [x] Avoided token issuance and runtime session creation from onboarding.
- [x] Avoided protocol, generated output, migration, dependency, and repository interface changes.

## Documentation

- [x] Added ADR.
- [x] Added conversation memory.
- [x] Updated work queue, manifests, guides, README, and alpha goal pointers.
- [x] Added repository check rule metadata.

## Verification

- [x] Focused Go tests passed.
- [x] Full runtime Go tests passed.
- [x] Repository checks passed.
- [x] Results recorded in `verification.md`.
