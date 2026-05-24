# Checklist

## Requirement

- [ ] Original request recorded.
- [ ] Clarified requirement recorded.
- [ ] User-visible outcome recorded.
- [ ] Non-goals recorded.
- [ ] Unknowns recorded.
- [ ] Acceptance criteria recorded before implementation.

## Nakama/Pitaya Alignment

- [ ] Nakama capability family recorded or `no_mapping_applies` justified.
- [ ] Direct Nakama API compatibility remains false unless explicitly authorized.
- [ ] Pitaya remains deferred unless explicitly authorized.

## Architecture

- [ ] Affected modules identified.
- [ ] Ownership impact reviewed.
- [ ] Contracts reviewed.
- [ ] Generated output posture recorded.
- [ ] Migration posture recorded.
- [ ] Dependency posture recorded.
- [ ] Redaction posture recorded.

## Tests

- [ ] Positive behavior tests planned or not-applicable rationale recorded.
- [ ] Negative behavior tests planned or not-applicable rationale recorded.
- [ ] Permission/authentication tests planned or not-applicable rationale recorded.
- [ ] Persistence/protocol/integration tests planned or not-applicable rationale recorded.
- [ ] Repository checks planned.

## Implementation

- [ ] Implementation stayed inside declared boundaries.
- [ ] Forbidden scope was not added.
- [ ] Tests were added before or with implementation when behavior changed.

## Documentation And Memory

- [ ] English docs updated when public-facing behavior or workflow changed.
- [ ] Simplified Chinese translations updated when paired docs changed.
- [ ] ADR added when direction, architecture, boundary, dependency, contract, protocol, generated output, or public behavior changed.
- [ ] Conversation memory added when maintainer intent or product direction matters.
- [ ] Manifests and AGENTS guides updated when continuation or ownership changed.

## Verification

- [ ] Verification commands run or unavailable checks recorded.
- [ ] Results recorded in `verification.md`.
