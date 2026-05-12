# Request

Date: 2026-05-13
Change ID: `adopt-first-runtime-dependencies`
Type: standard

## Maintainer Request

The maintainer clarified the decision process for technical branch points:

```text
这些你按照专业评估来定。这些问题越专业越好，你自己评估就可以。以后这种问题就是要来跟我确认，我允许你自己评估决定，你才能自己评估决定。
```

## Clarified Requirement

Record that the agent may professionally evaluate and decide technical dependency and runtime preparation questions after the maintainer grants that authority.

Use that authority to adopt the first foundational Go runtime dependencies where the evidence is strong enough and the boundary is narrow enough:

- WebSocket server library.
- Protobuf generation/runtime tooling.
- PostgreSQL driver.
- Migration tooling.

Keep deferred dependency categories deferred when the first runtime slice does not need them yet.

## User-Visible Outcome

The repository should show which foundational dependencies are accepted before runtime implementation begins, why they were selected, where they may be used directly, and where they are forbidden.

Future agents should know that professional technical sub-decisions may be evaluated by the agent after authorization, while product direction, constitutional principles, major architecture patterns, cost or operational commitments, and scope changes still require maintainer confirmation.

## Non-Goals

- Do not add Go implementation code yet.
- Do not create `go.mod` yet.
- Do not install or vendor runtime dependencies yet.
- Do not adopt S3 client tooling, MinIO deployment, observability, or an external Go test framework yet.
- Do not change public inventory contracts.
- Do not weaken the requirement that generated files are immutable to non-system agents.

## Acceptance Criteria

- [ ] Conversation memory records the maintainer's authorization and the boundary around it.
- [ ] An ADR records the technical-decision authorization model.
- [ ] An ADR records accepted first runtime dependencies and rejected or deferred alternatives.
- [ ] `.arch/dependencies.yaml` marks accepted slots only when linked to an adoption record.
- [ ] `.arch/runtime.yaml` records the accepted dependency choices and still marks implementation as not started.
- [ ] Repository agent guides explain the decision authorization and dependency boundaries.
- [ ] Verification records exactly what was checked.
