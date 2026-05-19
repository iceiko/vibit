# Verification

Verified:

- buf lint
- buf generate
- gofmt -w runtime/internal/app/connection_binding.go runtime/internal/app/connection_binding_test.go runtime/internal/app/route_authentication.go runtime/internal/platform/protocol/protobuf/connection_binding.go runtime/internal/platform/protocol/protobuf/connection_binding_test.go runtime/internal/platform/protocol/protobuf/frame_handler.go runtime/internal/platform/protocol/protobuf/frame_handler_test.go runtime/internal/platform/protocol/protobuf/inventory_bridge.go runtime/internal/platform/transport/ws/server.go runtime/internal/platform/transport/ws/server_test.go runtime/cmd/vibit-server/main.go runtime/cmd/vibit-server/main_test.go
- cd runtime && go test ./...
- node -c tools/vibit
- node tools/vibit check generated --json
- node tools/vibit check protocol --json
- node tools/vibit check runtime --json
- node tools/vibit check module authentication --json
- node tools/vibit check work --json
- node tools/vibit check memory --json
- node tools/vibit check change implement-first-message-connection-binding --json
- node tools/vibit check all --json
- git diff --check

Notes:

- `node tools/vibit check all --json` passed with one existing warning: `runtime.identity_boundary` reports that `runtime/internal/platform/persistence/postgres/authentication_repository.go` mentions credential dependency and should remain behind the ratified boundary.

Not verified:

- Live PostgreSQL end-to-end connection binding was not run; the default repository checks do not require live PostgreSQL for this slice.

Not applicable:

- Migration verification is not applicable because no migration source was added.
