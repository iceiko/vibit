# Conversation: Generated Output Standard

Date: 2026-05-13

Related changes:

- `changes/2026-05-13-add-generated-output-standard/`

Related artifacts:

- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `decisions/ADR-0017-generated-output-standard.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `tools/vibit`

## Context

The maintainer asked to continue development after the Protobuf envelope and generation configuration were committed and pushed.

The active project direction is to prepare carefully before runtime implementation, keep the project self-bootstrapping and controllable, and avoid drifting into disposable demo code.

## Maintainer Narrative

The maintainer's current instruction was to continue.

Prior context from the maintainer emphasized that the project should not rush into a minimal example too early. The framework is targeted and long-lived, so necessary preparation is appropriate when it directly protects the project goal.

## Agent Response Summary

The agent committed and pushed the Protobuf envelope and generation configuration change, then chose the next bounded step: define and verify generated-output rules before committing generated Go Protobuf files or adding broad runtime code.

The agent treated generated output control as a practical risk reducer for future agent work.

## Decisions

- Add a generated output standard before running or committing generated Go Protobuf output.
- Keep generated Protobuf output planned until the accepted toolchain is available.
- Extend `node tools/vibit check generated` so generated Go Protobuf files must look like `protoc-gen-go` output and trace back to existing `.proto` sources.

## Artifacts

- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `decisions/ADR-0017-generated-output-standard.md`
- `changes/2026-05-13-add-generated-output-standard/`
- `tools/vibit`

## Open Questions

- When should the first actual `buf generate` run happen in an environment with Buf and Go available?
- Should generated contract files use a custom vibit generator before or after the first protocol adapter slice?

## Follow-Up

- Verify the new generated output check.
- Commit and push the generated output standard.
- After this standard is stable, prepare the first bounded runtime protocol adapter or generation run.

## Redaction Notes

No secrets were included.
