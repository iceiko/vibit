# Request

## Original Request

The maintainer asked to continue after the game protocol framework had been defined:

```text
继续推进
```

## Clarified Requirement

Define the first Protobuf protocol source files and generation configuration for the accepted WebSocket Protobuf protocol direction.

The change should create the envelope `.proto`, the first inventory `.proto`, Buf configuration, and checks that let agents verify the protocol source shape without relying on visual inspection.

## User-Visible Outcome

Maintainers and agents can see concrete Protobuf source files for the protocol envelope and inventory wire messages, along with root Buf configuration and protocol checks that validate them.

## Non-Goals

- Do not write WebSocket runtime handlers.
- Do not write Go runtime business logic.
- Do not generate or commit Go Protobuf output.
- Do not add room, match, presence, stream, reconnect, realtime input, or state synchronization behavior.
- Do not introduce authentication or player modules.

## Unknowns

- Exact generated-output trace marker format for Go Protobuf files.
- Whether Buf remote plugin usage should later be pinned or replaced with local plugin invocation.
- Whether `payload_type` needs a stricter naming convention before clients exist.
- Reserved field-number policy for public protocol evolution.

## Acceptance Criteria

- [ ] `buf.yaml` and `buf.gen.yaml` exist and point to the accepted source root and generated output path.
- [ ] `proto/vibit/protocol/v1/envelope.proto` defines the first envelope source.
- [ ] `proto/vibit/inventory/v1/inventory.proto` defines inventory messages aligned with registered contracts.
- [ ] `node tools/vibit check protocol` validates Buf config, envelope shape, enum values, Go package options, source traces, messages, and fields.
- [ ] Architecture manifests and public docs identify the first Protobuf source and generation configuration.
- [ ] ADR, conversation log, and change spec record why generated Go output was not committed.
