# Dependency Adoption Standard

Status: Draft v0.1
Last updated: 2026-05-12
Scope: Foundational runtime and tooling dependencies

This document defines how vibit evaluates foundational dependencies before they become part of the architecture.

Use this standard together with `ADR-0010`.

## 1. Purpose

vibit should use mature open-source projects when they reduce risk and improve long-term quality. It should not adopt dependencies casually just because they are popular or convenient.

Dependency adoption records make dependency choices inspectable for future agents. They explain why a package is used, where it is allowed, what boundary contains it, and how it can be replaced.

## 2. When Required

A dependency adoption record is required before adding or requiring any dependency that shapes:

- Transport
- Protocol generation
- Persistence
- Migration tooling
- Object storage
- Dispatch
- Module loading
- Lifecycle management
- Observability
- Testing strategy
- Generated code conventions

Small developer-only utilities may use a lighter change-spec note if they do not affect architecture, generated output, public behavior, or runtime shape.

## 3. Location

Machine-readable dependency status lives in:

```text
.arch/dependencies.yaml
```

Reusable adoption template:

```text
docs/_templates/dependency-adoption.md
```

The final adoption record may be:

- A dedicated ADR under `decisions/`.
- A dependency section inside a change spec, when the dependency is small and narrowly scoped.

Foundational dependencies should normally use a dedicated ADR.

## 4. Required Evaluation

An adoption record must evaluate:

- Problem solved
- Package or tool identity
- Ecosystem role
- Maintenance activity
- License compatibility
- API stability
- Production adoption signals
- Security and supply-chain risk
- Operational fit
- Agent readability
- Testability
- Generated-code compatibility
- Abstraction boundary
- Allowed import or invocation locations
- Forbidden import or invocation locations
- Replacement path
- Verification path

High stars, frequent usage, and reputation are useful signals. They are not sufficient by themselves.

## 5. Boundary Rules

Domain modules must not directly import or invoke foundational third-party dependencies for transport, protocol, persistence, object storage, or framework behavior.

Preferred ownership:

- Transport libraries belong in platform transport adapters.
- Protobuf tooling belongs in generation tooling and generated protocol packages.
- PostgreSQL drivers belong in platform persistence adapters.
- Migration tools belong in platform migration tooling.
- S3 SDKs and MinIO clients belong in platform object-storage adapters.
- Test frameworks belong in test infrastructure, not domain logic.

Domain modules should depend on vibit-owned interfaces, generated clients, contract types, repositories, policies, and service abstractions.

## 6. Dependency Status Values

Use these status values in `.arch/dependencies.yaml`:

```text
proposed
accepted
deferred
rejected
superseded
```

Use `proposed` when the dependency slot is known but no implementation has been selected.

Use `accepted` only after the adoption record is complete and linked.

Use `deferred` when the dependency category is real but not needed for the next implementation step.

Use `rejected` when a plausible dependency should not be chosen.

Use `superseded` when a previous dependency choice has been replaced.

## 7. Agent Rules

Agents must:

- Read `.arch/dependencies.yaml` before adding foundational dependencies.
- Create or update a dependency adoption record before changing dependency status to `accepted`.
- Keep dependencies behind the declared abstraction boundary.
- Update AGENTS guides when dependency boundaries affect module work.
- Run the relevant verification commands after changing dependency records.

Agents must not:

- Add a foundational dependency directly to domain module code.
- Treat popularity as a complete evaluation.
- Hide dependency decisions inside implementation code.
- Accept a dependency without a replacement path or verification path.

## 8. Verification Direction

Current verification is documentation and manifest based:

```bash
node tools/vibit check architecture
node tools/vibit check schemas
node tools/vibit check memory
node tools/vibit check all
```

Future checks should verify:

- Every `accepted` foundational dependency has an adoption record.
- Domain modules do not import forbidden dependency packages.
- Platform adapters are the only allowed direct dependency owners.
- Generated files declare their generator dependencies.

## 9. Open Questions

- Should `.arch/dependencies.yaml` eventually be fully schema-validated?
- Should dependency records include machine-readable package coordinates?
- Should replacement path and license review become required structured fields?
- Should dependency checks inspect Go imports once runtime code exists?
