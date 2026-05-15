# Request

Implement `W-0099 Implement environment verifier key loader`.

The loader must stay under `runtime/internal/app/authentication`, decode the required process environment verifier key values, call `NewVerifierKeySet`, and keep startup wiring, local secret files, dotenv behavior, CLI secret input, KMS, cloud secret-manager integration, token generation, credential generation, verifier digest computation, verifier comparison, authentication service behavior, protocol carriers, repositories, migrations, authentication dependencies, and production authentication behavior deferred.
