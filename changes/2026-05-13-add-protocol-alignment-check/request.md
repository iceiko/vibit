# Request

## Original Request

```text
继续下一步
```

## Clarified Requirement

Add the first machine-checkable manifest-to-Protobuf alignment check before creating `.proto` files or generated protocol output.

The check should make it clear which semantic contracts are expected to map to which Protobuf source files and message names.

## User-Visible Outcome

Agents can run a dedicated protocol check before adding or changing Protobuf files.

## Non-Goals

- Do not create `.proto` message files.
- Do not generate Go Protobuf output.
- Do not add Buf configuration yet.
- Do not add Go runtime business implementation.
- Do not adopt new dependencies.

## Unknowns

- Final Protobuf package naming rules may become stricter when the first `.proto` file is added.
- The exact envelope and service shape for WebSocket messages is still deferred.

## Acceptance Criteria

- `node tools/vibit check protocol` exists.
- `node tools/vibit check protocol --json` exists.
- `node tools/vibit check all --json` includes protocol checks.
- Current planned Protobuf paths pass because no `.proto` files have been created yet.
- The check defines expected Protobuf paths and message names for registered command, query, and event contracts.
- Documentation and Chinese translations are updated.
