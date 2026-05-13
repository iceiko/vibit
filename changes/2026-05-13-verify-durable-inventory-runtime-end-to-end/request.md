# Request

## Original Request

Continue advancing ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0021`, which verifies the durable inventory runtime end to end. Add a repeatable, opt-in live PostgreSQL verification path for the existing WebSocket Protobuf `GrantItem` and `GetInventory` request loop, record live database unavailability when no disposable DSN exists, and keep domain, application, transport, protocol, persistence, and migration boundaries intact.

## User-Visible Outcome

Maintainers and future agents can run a single Go integration test to verify migration apply/status and the persistent inventory request loop against a disposable PostgreSQL database. Default repository checks still work without PostgreSQL.

## Non-Goals

- Close `M-002` without live PostgreSQL verification.
- Introduce outbox delivery, authentication, session validation, player accounts, item catalog, currency, rewards, or match sessions.
- Make PostgreSQL, Docker, Podman, or any service manager mandatory for default tests.
- Apply migrations automatically during normal server startup.
- Change the WebSocket Protobuf envelope shape.

## Unknowns

- No disposable `VIBIT_POSTGRES_TEST_DSN` is available in the current environment, so the live database branch is expected to skip during local verification.

## Acceptance Criteria

- [ ] Add a live PostgreSQL request-loop verification test that skips clearly when `VIBIT_POSTGRES_TEST_DSN` is unset.
- [ ] Verify migration apply/status through the existing migration runner when the disposable DSN is available.
- [ ] Verify persistent `GrantItem` followed by `GetInventory` through the Protobuf frame handler when the disposable DSN is available.
- [ ] Preserve current default in-memory runtime behavior and the explicit PostgreSQL runtime store path.
- [ ] Update manifests, runtime guidance, and bilingual docs to describe the opt-in live verification path.
- [ ] Run repository verification and record unavailable live PostgreSQL verification honestly.
