# Verification

Status: Verified

## Verified

- `node --check tools/vibit`
- `node tools/vibit check generated`
- `node tools/vibit check generated --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change add-generated-output-standard --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens outside `.git/`, `.vibit.local.env`, and `node_modules/`

## Not Verified

- `buf lint`: local toolchain availability not checked in this change.
- `buf generate`: generated output is intentionally not produced in this change.
- `go test ./...`: no Go runtime implementation or tests are added in this change.
- `go vet ./...`: no Go runtime implementation is added in this change.

## Notes

This change defines generated-output guardrails before generated Go Protobuf files are committed.
