# Verification

## Verified

- `node tools/vibit check architecture`
- `node --check tools/vibit`
- `node tools/vibit check schemas`
- `node tools/vibit check memory`
- `node tools/vibit check contracts`
- `node tools/vibit check generated`
- `node tools/vibit check runtime`
- `node tools/vibit check all --json`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .`
- `git diff --check`

## Not Verified

- Go runtime tests are not verified because no Go runtime implementation exists yet.
- WebSocket behavior is not verified because no WebSocket adapter exists yet.
- Protobuf generation is not verified because proto layout and tooling are intentionally deferred to dependency adoption.

## Notes

`node tools/vibit check runtime` now treats the missing Go runtime as not applicable instead of requiring removed TypeScript tests.
