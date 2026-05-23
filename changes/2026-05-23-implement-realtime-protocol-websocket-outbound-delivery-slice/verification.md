# Verification

Required commands:

```bash
buf generate
cd runtime && go test ./internal/platform/protocol/protobuf ./internal/platform/transport/ws ./internal/app/realtime
cd runtime && go test ./...
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.realtime_protocol_websocket_outbound_delivery_implementation
node tools/vibit check generated --json
node tools/vibit check protocol --json
node tools/vibit check change implement-realtime-protocol-websocket-outbound-delivery-slice --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Status: Verified.

Results:

- `buf generate` passed.
- `cd runtime && go test ./internal/platform/protocol/protobuf ./internal/platform/transport/ws ./internal/app/realtime` passed.
- `cd runtime && go test ./...` passed.
- `node -c tools/vibit` passed.
- `node tools/vibit inspect next --json` passed and reports `W-0219` as next_ready.
- `node tools/vibit inspect rule runtime.realtime_protocol_websocket_outbound_delivery_implementation` passed.
- `node tools/vibit check generated --json` passed.
- `node tools/vibit check protocol --json` passed.
- `node tools/vibit check change implement-realtime-protocol-websocket-outbound-delivery-slice --json` passed.
- `node tools/vibit check work --json` passed.
- `node tools/vibit check runtime --json` passed with one pre-existing warning.
- `node tools/vibit check memory --json` passed.
- `node tools/vibit check schemas --json` passed after setting this change spec status to `verified`.
- `node tools/vibit check all --json` passed after setting this change spec status to `verified`.
- `git diff --check` passed.

Known notes:

- `runtime.identity_boundary` still reports the pre-existing warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go`; it is unrelated to W-0218.
