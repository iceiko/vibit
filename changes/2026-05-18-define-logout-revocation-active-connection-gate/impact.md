# Impact

This change adds a gate-only standard and repository checks.

It defines:

- Future ownership for logout/revocation policy.
- Future policy questions for token revocation, runtime session revocation, and active connection invalidation.
- A conservative future first posture centered on presented-token logout.
- A requirement that connection registries exist before active sockets can be targeted.
- Redaction and public error boundaries.
- Nakama and Pitaya reference mapping.
- Future test expectations.

It does not add Go runtime behavior, Protobuf source, generated output, migrations, dependencies, WebSocket close behavior, or production logout behavior.
