# Impact

## Architecture Impact

This change strengthens the generated-file boundary without changing runtime package layout or public protocol semantics.

It makes `runtime/internal/generated/proto/` an actively checked generated-output root. Future generated Go Protobuf files must trace to `.proto` sources instead of being accepted as arbitrary Go code.

## Contract Impact

No command, query, event, error, or permission contracts change.

## Runtime Impact

No runtime implementation code is added.

The runtime skeleton remains the only Go runtime shape in this change.

## Protocol Impact

No `.proto` schema changes are made.

The generated output path for existing Protobuf sources becomes more strictly documented and checked.

## Documentation Impact

Adds the generated-output standard and Chinese translation. Updates architecture, runtime, protocol, README, and agent guidance so future agents read the standard before touching generated output.

## Compatibility Impact

No public API, event, or data compatibility impact.

Future generated files may fail checks if they are hand-written or do not match the expected `protoc-gen-go` output shape.
