# Verification

Verified:

- `node tools/vibit inspect contract --module runtime --type command --id ValidateSession`
- `node tools/vibit inspect contract --module runtime --type event --id SessionValidated`
- `node tools/vibit inspect contract --module player --type command --id CreatePlayerAccount`
- `node tools/vibit check contracts --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-runtime-contract-inspection --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `git diff --check`

Not verified:

- None.

Warnings:

- `node tools/vibit check work --json` reported one warning because there was no `next_ready` work item after `W-0035`. That warning was expected before closing the completed milestone and selecting or blocking the next direction.
- `node tools/vibit check all --json` inherited the same `check work` warning.

Not applicable:

- Live PostgreSQL verification; this change does not add database behavior.
- Authentication, token, credential, player account persistence, session persistence, Protobuf envelope, WebSocket handshake, and runtime player handler verification are not applicable because this change does not implement or change those surfaces.
