# Request

Date: 2026-05-13
Change ID: `define-go-runtime-layout`
Type: standard

## Maintainer Request

The maintainer asked to continue:

```text
继续推进
```

## Clarified Requirement

Continue runtime readiness by defining the package layout and boundary rules needed before creating Go runtime implementation code.

The next implementation step should not have to infer:

- Where `go.mod` lives.
- Which Go packages own transport, persistence, migrations, dispatch, generated protocol code, and domain module logic.
- Where `.proto` sources and Buf configuration live.
- Where generated Go protocol output lives.
- Where SQL migrations live.
- How command handlers enter a transaction boundary and publish events.

## User-Visible Outcome

Future agents should be able to inspect `.arch/runtime.yaml` and the new ADR to know the first Go runtime shape before creating code.

The repository should still avoid implementation code in this change.

## Non-Goals

- Do not create `go.mod`.
- Do not add Go source files.
- Do not generate Protobuf output.
- Do not add migrations.
- Do not add runtime dependencies to a package manager.
- Do not change public command, query, event, error, or permission contracts.

## Acceptance Criteria

- [ ] ADR records the Go runtime package layout and boundary rules.
- [ ] `.arch/runtime.yaml` records the layout in machine-readable form.
- [ ] `.arch/conventions.yaml` links the new runtime layout decision.
- [ ] AGENTS guides explain where future runtime code belongs.
- [ ] README files mention the first runtime package layout.
- [ ] Conversation memory records the continuation and decision.
- [ ] Verification records exactly what was checked.
