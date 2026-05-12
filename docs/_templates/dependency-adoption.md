# Dependency Adoption: <dependency name>

Status: Proposed
Date: YYYY-MM-DD
Decision Makers: Maintainer, Agent
Related changes:

- None yet.

Related decisions:

- `ADR-0010`

## Dependency Identity

- Name:
- Category:
- Package or tool:
- Repository:
- Documentation:
- License:
- Intended owner:

## Problem

Describe the concrete problem this dependency solves.

## Decision

State whether the dependency is proposed, accepted, deferred, rejected, or superseded.

## Evaluation

### Maintenance Activity

Record recent release, commit, issue, or ecosystem signals.

### License Compatibility

Record the license and any known obligations or uncertainty.

### API Stability

Describe stability of the public API or CLI.

### Production Adoption Signals

Record stars, usage, ecosystem position, or known production use when relevant.

### Security And Supply Chain

Record security posture, vulnerability handling, provenance, and update strategy.

### Operational Fit

Explain fit for vibit's Go, WebSocket, Protobuf, PostgreSQL, or S3-compatible object-storage direction.

### Agent Readability

Explain whether agents can understand and use the dependency without hidden framework magic.

### Testability

Explain how behavior using this dependency can be tested.

### Generated-Code Compatibility

Explain whether this dependency affects generated files, templates, or contract alignment.

## Boundary

Allowed direct use:

- Path or layer 1

Forbidden direct use:

- Path or layer 1

Vibit-owned abstraction:

- Interface or package name

## Replacement Path

Describe how vibit could replace this dependency later.

## Verification

Commands or checks required after adopting or updating this dependency:

```bash
node tools/vibit check all
```

## Open Questions

- Question 1
