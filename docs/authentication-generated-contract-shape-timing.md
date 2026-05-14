# Authentication Generated Contract Shape Timing Standard

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Runtime authentication semantic contracts and generated Go contract shape timing

Use this standard together with `docs/generated-output.md`, `docs/authentication-contract-error-permission-surfaces.md`, `docs/runtime-authentication-implementation-boundary.md`, `ADR-0017`, `ADR-0028`, `ADR-0036`, `ADR-0037`, and `ADR-0038`.

## 1. Purpose

Runtime authentication now has semantic contracts, storage schema boundaries, a storage-neutral repository interface, a PostgreSQL adapter, and an application-owned implementation boundary.

The next risk is that an agent may start service interfaces or runtime behavior directly from prose and repository code, bypassing the machine-readable contract shapes that already help inventory and player work stay predictable.

This standard first decided the timing and boundary for generated Go authentication contract shapes. `W-0089` now completes the first generation slice without adding runtime behavior.

## 2. Timing Decision

Generated Go authentication contract shapes should be introduced before application authentication service interfaces and before runtime authentication behavior.

The recommended order is:

1. Semantic authentication contracts under `contracts/runtime/authentication/`.
2. Runtime authentication implementation boundary.
3. Authentication generated contract shape timing decision.
4. Generator and check support for runtime authentication family shapes.
5. Generated authentication contract shape files.
6. Application authentication service interface boundary.
7. Token generation, verifier comparison, login execution, token validation, logout execution, cleanup, protocol carriers, and runtime behavior through later gated work.

`W-0088` completed step 3 only. `W-0089` completes steps 4 and 5 by adding generator/check support and metadata-only generated files. These steps do not authorize service implementations, handlers, token behavior, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, or migration schema changes.

## 3. Source And Output

The allowed source set is the registered runtime authentication contract family:

```text
contracts/runtime/authentication/commands/*.yaml
contracts/runtime/authentication/events/*.yaml
contracts/runtime/authentication/errors/*.yaml
contracts/runtime/authentication/permissions/*.yaml
```

The registry source is `.arch/contracts.yaml`, under the runtime `authentication` family.

The output root is:

```text
runtime/internal/generated/contracts/runtime/authentication/
```

The file shape is:

```text
runtime/internal/generated/contracts/runtime/authentication/<contract-type>/<ContractID>.go
```

Examples:

```text
runtime/internal/generated/contracts/runtime/authentication/commands/AuthenticateWithDeviceCredential.go
runtime/internal/generated/contracts/runtime/authentication/events/TokenIssued.go
runtime/internal/generated/contracts/runtime/authentication/errors/authentication_errors.go
runtime/internal/generated/contracts/runtime/authentication/permissions/authentication_permissions.go
```

Runtime contract families require the family segment because `runtime` may own more than one semantic family, such as `session` and `authentication`.

The Go package name for the first authentication shape files is:

```text
runtimeauthenticationcontracts
```

## 4. Immutability

Generated authentication contract shape files are immutable to non-system agents.

If generated authentication contract shape output is wrong, agents must change one of these sources instead of patching generated files:

- The semantic contract source under `contracts/runtime/authentication/`.
- The contract registry in `.arch/contracts.yaml`.
- The generator.
- The generated-output standard.
- The relevant change spec or ADR, if a `generated_file_override` is explicitly granted.

No `generated_file_override` is granted by this standard.

## 5. Check Requirements

Generated authentication contract shapes may be committed only when repository tooling supports these checks:

- `node tools/vibit generate contract-shapes all` can generate the runtime authentication family from semantic contracts.
- `node tools/vibit check generated --json` can detect missing, stale, or drifted runtime authentication family shapes.
- `node tools/vibit inspect generated --json` can report runtime authentication generated shape status in machine-readable form.
- `node tools/vibit check contracts --json` still verifies the semantic contract source and registry before generation.
- `node tools/vibit check runtime --json` still proves generated authentication shapes do not add runtime authentication behavior.
- `runtime.selected_login_token_boundary` and `runtime.authentication_implementation_boundary` distinguish metadata-only generated shapes from runtime authentication implementation.
- `node tools/vibit check all --json` includes the generated-shape checks.

`W-0089` updates checks before committing the generated output. Future changes must not weaken selected login/token, authentication implementation, generated-output, protocol, WebSocket, dependency, repository, or migration guards.

## 6. Relationship To Runtime Behavior

Generated authentication contract shapes are metadata-only.

They may guide naming, service interface planning, tests, and future protocol mapping, but they must not implement:

- Login execution.
- Token generation.
- Token verifier comparison.
- Access-token validation.
- Logout execution.
- Refresh behavior.
- Cleanup jobs.
- WebSocket routing.
- Protobuf envelope behavior.
- Persistence.
- Domain dispatch.

Application service interfaces may be designed after generated shapes exist, but service code must still be separately gated.

## 7. Reference Alignment

Nakama and Pitaya remain capability and vocabulary references. They show that production game backends need stable authentication, session, request, and route concepts.

vibit should adapt that lesson into agent-native structure:

- Semantic contracts define the names and intent.
- Generated contract shapes make the contract surface inspectable by agents.
- Handwritten application logic remains separate and explicitly gated.
- Wire protocol and WebSocket carriers remain separate decisions.

This standard does not copy Nakama or Pitaya public APIs.

## 8. Verification

For this timing decision, verification should include:

```bash
node tools/vibit check contracts --json
node tools/vibit inspect generated --json
node tools/vibit check generated --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change decide-authentication-generated-contract-shape-timing --json
node tools/vibit check all --json
git diff --check
```

Runtime Go tests are not required for the timing decision because it does not add or modify Go runtime behavior.

Generation work must run generator, generated-output, runtime, and full repository checks after files are produced.

## 9. Migration Path

The first migration path is:

1. Record this timing decision and output boundary.
2. Add a bounded work item that authorizes generator/check support and generated authentication shape output.
3. Extend generated-output tooling for runtime family-aware paths.
4. Generate files through tooling, not by hand.
5. Verify source trace, drift, stale files, and runtime behavior boundaries.
6. Only then design application authentication service interfaces.
