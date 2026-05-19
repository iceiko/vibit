# Impact

`LogoutAccessToken` now executes the first ratified logout posture:

- The presented access token must be a Base64URL unpadded 32-byte opaque token.
- Missing and malformed proofs are rejected before unit-of-work creation.
- The service computes the token lookup digest and finds the token record.
- Token kind, active status, issue/expiry time, audience, verifier algorithm, verifier version, and verifier key id are checked before revocation.
- The service computes and compares the verifier digest before revocation.
- The token record is revoked with `logout_presented_access_token`.
- Success is returned only after commit.

Public invalid-token behavior remains collapsed for lookup miss, revoked, expired, wrong kind, wrong audience, unsupported verifier metadata, unknown key id, and verifier mismatch.

The implementation does not call the player repository, session repository, WebSocket transport, protocol adapter, connection registry, or cleanup jobs. It does not add Protobuf logout messages or direct Nakama/Pitaya API compatibility.

Nakama alignment: revoked authentication material no longer authorizes future protected requests when validation is invoked.

Pitaya alignment: logout remains application-owned and separate from connection/session infrastructure side effects.
