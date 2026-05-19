# Logout Access Token Behavior Gate

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Gate-only boundary for future presented opaque access-token logout behavior
Depends on: `docs/authentication-service-behavior-implementation-gate.md`, `docs/access-token-validation-service-behavior-gate.md`, `docs/logout-revocation-active-connection-gate.md`, `docs/session-creation-composition-gate.md`, `docs/bound-identity-route-policy-gate.md`, `decisions/ADR-0071-logout-revocation-active-connection-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0072`

The paired Simplified Chinese translation is `docs/logout-access-token-behavior-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

`LogoutAccessToken` exists as a semantic contract and as a fail-closed service method, but vibit has not yet selected the behavior boundary for executing logout.

The previous gate established that logout, runtime session revocation, and active WebSocket connection invalidation are separate lifecycle decisions. This gate narrows the next future behavior to the first safe logout posture: revoke only the presented opaque access token, through the application authentication service, without changing runtime session state or active WebSocket connections.

Nakama provides the important product pressure: authenticated session material needs a lifecycle that includes logout, refresh, expiration, and rejection of revoked material for future gameplay requests. Pitaya provides the layering pressure: handlers should receive context from session/connection infrastructure and should not parse credentials or own connection lifecycle side effects.

vibit adapts those lessons by making first logout behavior token-record scoped, transactionally explicit, redacted, and separated from socket close, connection registry, reconnect, protocol route exposure, refresh, and logout-all behavior.

```yaml
logout_access_token_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0156
decision: ADR-0072
check_rule: runtime.logout_access_token_behavior_gate
future_service_owner: runtime/internal/app/authentication
future_repository_capability_owner: runtime/internal/modules/authentication
future_transaction_owner: runtime/internal/app
existing_contract: contracts/runtime/authentication/commands/LogoutAccessToken.yaml
existing_service_method: runtime/internal/app/authentication.Service.LogoutAccessToken
existing_service_behavior: fail_closed_not_implemented
first_logout_scope: presented_access_token_only
proof_shape: opaque_base64url_unpadded_32_byte_access_token
proof_carrier_status: already_decoded_service_request_only
token_lookup_before_revocation_required: true
token_verifier_comparison_before_revocation_required: true
token_status_must_be_active_before_revocation: true
revoked_token_public_behavior: invalid_token
already_revoked_token_public_behavior: invalid_token
expired_token_public_behavior: invalid_token
runtime_session_revocation_added: false
active_connection_invalidation_added: false
connection_registry_added: false
websocket_close_policy_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
refresh_behavior_added: false
logout_all_sessions_added: false
admin_revocation_added: false
cleanup_jobs_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This standard does not implement logout.

## 2. Ownership

Future logout execution must stay application-owned:

```yaml
future_logout_service_owner: runtime/internal/app/authentication
future_unit_of_work_owner: runtime/internal/app
repository_interface_owner: runtime/internal/modules/authentication
postgres_adapter_owner: runtime/internal/platform/persistence/postgres
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
domain_handler_owner: runtime/internal/modules/*
```

Rules:

- The authentication service owns proof validation, token lookup, verifier comparison, public error collapse, and the request to revoke the token.
- The authentication repository owns storage-neutral token record mutation vocabulary only. It must not parse raw access-token text, compute digests, compare verifiers, decide logout proof validity, or close connections.
- The unit-of-work boundary must wrap token lookup, verifier comparison decision inputs, token revocation mutation, and commit outcome.
- WebSocket transport must not parse logout credentials or decide logout side effects.
- Protobuf adapters may expose logout only after a later protocol route gate.
- Domain modules must receive an already-authenticated request identity or a logout result; they must not parse access tokens or call token repositories directly.

## 3. Future Behavior Sequence

A later implementation may execute logout only in this order:

```yaml
future_logout_sequence:
  - reject_missing_or_malformed_access_token_before_unit_of_work
  - compute_access_token_lookup_digest
  - begin_application_unit_of_work
  - obtain_authentication_repository_from_unit_of_work
  - find_token_record_by_lookup_digest
  - require_token_kind_access_token
  - require_token_status_active
  - require_not_expired_at_service_clock_now
  - require_expected_audience
  - require_supported_verifier_algorithm_and_version
  - compute_access_token_verifier_digest_using_record_key_id
  - compare_verifier_digest_constant_time
  - revoke_presented_token_record_with_reason_logout_presented_access_token
  - commit_unit_of_work
  - return_revoked_result_after_commit_only
```

Rules:

- Missing or malformed token proof must be rejected before opening a unit of work.
- Token lookup must use the existing lookup digest helper; raw token text must not reach the repository.
- Verifier digest comparison must happen before revocation so a lookup collision or wrong raw token cannot revoke a token record.
- The first implementation posture must revoke only the verified presented token record.
- A revoked, expired, wrong-audience, wrong-kind, unsupported, unknown-key, lookup-missing, or mismatched token must collapse to public invalid-token behavior.
- Repository unavailability must collapse to public token-store-unavailable behavior.
- No raw token text may be returned or stored in results, errors, logs, events, or test names.

## 4. Transaction Boundary

The future first logout behavior must be transactional for token state:

```yaml
transaction_boundary:
  includes:
    - token_record_lookup_by_lookup_digest
    - token_record_status_and_expiry_check
    - verifier_digest_comparison_decision_inputs
    - token_record_revocation_mutation
  excludes:
    - runtime_session_revocation
    - active_connection_invalidation
    - websocket_close
    - protocol_response_mapping
    - cleanup_jobs
```

Rules:

- The service must return `revoked` only after the unit of work commits.
- If commit fails, the public result must not claim logout success.
- Logout must be idempotence-explicit: this gate selects fail-closed invalid-token behavior for already revoked, expired, or inactive tokens. A later ADR may choose idempotent success, but must change the standard and tests deliberately.
- Runtime session revocation is not part of the first transaction. The `runtime_sessions` row may remain active until a later session revocation policy changes it.
- Active WebSocket connection invalidation is outside SQL transaction control and is not authorized here.

## 5. Public Result And Error Boundary

Future public logout behavior must be minimal and redacted:

```yaml
candidate_public_success:
  status: revoked
  revoked: true
  logout_scope: presented_access_token
  token_type: opaque_access_token

candidate_public_errors:
  missing: AUTHENTICATION_TOKEN_MISSING
  malformed: AUTHENTICATION_TOKEN_MALFORMED
  invalid: AUTHENTICATION_TOKEN_INVALID
  unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
```

Rules:

- Public invalid-token behavior must not distinguish lookup miss, already revoked, expired, wrong audience, wrong token kind, unsupported verifier metadata, unknown verifier key id, verifier mismatch, or missing player account.
- Success must not include raw token text, lookup digest, verifier digest, verifier key id, Authorization header, cookie, query string, WebSocket subprotocol value, session id, connection id, or remote address.
- Internal errors must remain typed enough for tests and redacted observability, but must not leak secrets in `Error()`.
- If the future result includes `TokenRecordID`, it must be treated as an internal/audit-safe identifier only after an observability standard classifies it.

## 6. Relationship To Runtime Session And Active Connections

This gate intentionally keeps runtime session and active connection behavior separate:

```yaml
runtime_session_revocation_added: false
active_connection_invalidation_added: false
connection_registry_added: false
websocket_close_policy_added: false
session_last_seen_update_added: false
duplicate_connection_replacement_added: false
reconnect_epoch_behavior_added: false
```

Rules:

- Logging out the presented access token must not implicitly revoke all runtime sessions for the player.
- Logging out the presented access token must not implicitly revoke all tokens for the player, credential, device, or account.
- Logging out the presented access token must not implicitly close active WebSocket connections.
- Existing request-level access-token validation will reject the revoked token on later protected requests when validation is invoked.
- Bound-connection and bound-session route behavior after token revocation remains a later policy question.
- Reconnect, resume, duplicate connection replacement, connection epoch, and targeted kick/disconnect remain deferred.

## 7. Relationship To Protocol

This gate does not expose logout through Protobuf or WebSocket:

```yaml
protobuf_logout_route_added: false
protobuf_authentication_message_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
generated_output_added: false
```

Rules:

- `proto/vibit/authentication/v1/authentication.proto` must not add logout request or response messages under this gate.
- Generated Go Protobuf output must not be edited by hand.
- WebSocket transport must not parse logout credentials from headers, cookies, query strings, or subprotocol values.
- Protocol adapters must not call `LogoutAccessToken` until a later protocol logout route gate defines carrier, route name, response mapping, and error mapping.

## 8. Relationship To Existing Service Skeleton

The current service may remain fail-closed:

```yaml
current_logout_service_method: LogoutAccessToken
current_status: not_implemented
current_public_error: AUTHENTICATION_NOT_IMPLEMENTED
behavior_change_authorized_by_this_gate: false
```

Rules:

- This gate may guide a future implementation slice, but it must not change `runtime/internal/app/authentication/service.go`.
- `RefreshAccessToken` remains unsupported.
- Access-token validation behavior remains unchanged.
- Device credential login behavior remains unchanged.
- Session creation behavior remains unchanged.
- Route policy remains unchanged.

## 9. Future Test Expectations

A later implementation slice must add focused tests for:

- Missing token proof rejected before unit of work.
- Malformed token proof rejected before unit of work.
- Lookup miss collapses to invalid-token public behavior.
- Already revoked token collapses to invalid-token public behavior.
- Expired token collapses to invalid-token public behavior.
- Wrong token kind, audience, algorithm, version, or key id collapses to invalid-token public behavior.
- Verifier mismatch does not revoke the token record.
- Repository lookup/revocation/commit failures do not claim success.
- Successful logout calls `RevokeToken` exactly once with `logout_presented_access_token`.
- Success is returned only after commit.
- Raw token text, lookup digest, verifier digest, verifier key id, session id, connection id, Authorization header, cookie, query string, and WebSocket subprotocol value are absent from public result and error strings.
- Runtime session repository, connection registry, WebSocket transport, and Protobuf adapter are not called.

## 10. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

- Adopt the lesson that session material has lifecycle state and revoked material must not authorize future gameplay requests.
- Adapt logout to vibit's opaque server-issued access-token record model rather than Nakama-compatible JWT/session APIs.
- Keep refresh, logout-all, session management APIs, dashboard/admin revocation, and realtime socket invalidation as later explicit surfaces.

Pitaya reference mapping:

- Adopt the separation between connection/session infrastructure and handler logic.
- Adapt logout as application service behavior that may later coordinate with connection/session infrastructure through narrow interfaces.
- Keep frontend/backend cluster routing, distributed kick/disconnect, and server-to-server invalidation as later gates.

Nakama and Pitaya guide capability pressure only. This gate does not add direct public API compatibility with either project.

## 11. Non-Goals

This gate does not authorize:

- `LogoutAccessToken` implementation.
- Token revocation execution.
- Runtime session revocation.
- Active WebSocket connection invalidation.
- Connection registry.
- WebSocket close policy.
- Kick/disconnect behavior.
- Reconnect, resume, duplicate replacement, or epoch behavior.
- Protobuf logout route.
- Protocol session carrier.
- Existing Protobuf envelope changes.
- Refresh token behavior.
- Logout-all-sessions behavior.
- Admin revocation.
- Cleanup jobs.
- New dependencies.
- Memory durable session behavior.
- Direct Nakama or Pitaya API compatibility.

## 12. Required Follow-Up

The next implementation slice, if selected later, should be:

```text
implement_logout_access_token_behavior
```

That slice must keep the behavior inside `runtime/internal/app/authentication`, use existing repository and verifier helper boundaries, add focused tests, and preserve all protocol, transport, session revocation, active connection, reconnect, dependency, and direct compatibility deferrals unless separate ADRs authorize them.
