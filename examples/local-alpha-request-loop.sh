#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
TEST_NAME=TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout

printf '%s\n' "vibit local alpha request loop"
printf '%s\n' "path: local onboarding -> login -> bind connection -> protected inventory -> presence -> logout -> rejected revoked-token request"
printf '%s\n' "redaction: raw credentials, access tokens, digests, verifier keys, DSNs, and transport metadata are not printed"
printf '%s\n' "mode: focused Go E2E proof over existing runtime protocol handlers"
printf '\n'

cd "$ROOT_DIR/runtime"
go test ./internal/platform/protocol/protobuf -run "$TEST_NAME" -v
