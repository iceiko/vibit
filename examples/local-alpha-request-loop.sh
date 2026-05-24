#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)

printf '%s\n' "vibit local alpha request loop"
printf '%s\n' "compatibility entrypoint: delegating to examples/local-alpha-example-client.sh"
printf '%s\n' "path: local onboarding -> login -> bind connection -> protected inventory/presence/storage objects -> logout -> rejected revoked-token request"
printf '%s\n' "storage: authenticated own-player put/get/list/delete through existing WebSocket/Protobuf route flow"
printf '%s\n' "storage routes: storage.GetOwnStorageObject, storage.ListOwnStorageObjects, storage.PutOwnStorageObject, storage.DeleteOwnStorageObject"
printf '%s\n' "tests: TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout, TestStorageObjectsProtocolRouteLocalAlphaFlow"
printf '%s\n' "command: go test ./internal/platform/protocol/protobuf"
printf '%s\n' "redaction: raw credentials, access tokens, digests, verifier keys, DSNs, and transport metadata are not printed"
printf '%s\n' "mode: focused Go E2E proof over existing runtime protocol handlers"
printf '\n'

exec "$ROOT_DIR/examples/local-alpha-example-client.sh"
