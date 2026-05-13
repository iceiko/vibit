# Verification

Verified:

- `node --check tools/vibit`
- `node tools/vibit check protocol`
- `node tools/vibit check protocol --json`
- `node tools/vibit check architecture`
- `node tools/vibit check schemas`
- `node tools/vibit check memory`
- `node tools/vibit check contracts`
- `node tools/vibit check generated`
- `node tools/vibit check runtime`
- `node tools/vibit check all --json`
- `node tools/vibit check change define-protobuf-envelope-and-generation`
- `git diff --check`
- Secret scan excluding `.git`, `.vibit.local.env`, and `node_modules`

Not verified:

- `buf lint`: Buf CLI is not available in the local environment.
- `buf generate`: Buf CLI is not available, and generated Go Protobuf output is intentionally outside this change.
- `protoc`: `protoc` is not available in the local environment.
- `go test ./...`: Go is not available in the local environment, and no Go runtime source files were added.

Not applicable:

- PostgreSQL migrations: no persistence schema or migration was changed.
- Runtime integration tests: no WebSocket, Protobuf adapter, or server runtime behavior was implemented.
