# Plan

1. Add `proto/vibit/authentication/v1/authentication.proto` with the gate-authorized request, response, and account creation intent enum.
2. Run `buf generate` and keep generated Go Protobuf output immutable.
3. Add a Protobuf authentication bridge for request and response payload mapping.
4. Register the generated authentication Protobuf package in the payload registry.
5. Delegate authentication route bridge handling before inventory route handling.
6. Add application bootstrap handlers for the public login route.
7. Map service public errors to application errors without leaking secret material.
8. Add a `TransactionalDispatcher` bypass list and configure it for `AuthenticateWithDeviceCredential`.
9. Register the public login route in PostgreSQL startup composition only when the authentication service is composed.
10. Add focused tests and update repository checks, manifests, guides, change spec, and conversation memory.
