# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-protocol-logout-route-gate --json`
- `node tools/vibit inspect next --json`

Planned before final close-out:

- `node tools/vibit check memory --json`
- `node tools/vibit check change ratify-nakama-pitaya-product-parity-roadmap --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not applicable:

- Go tests are not required for this gate-only standard, though they are required for the following implementation slice.
- Live PostgreSQL verification is not required by this gate.
- Buf generation is not required because this gate does not add `.proto` sources.
