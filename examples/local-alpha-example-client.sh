#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
TEST_PATTERN='TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow|TestFriendsRelationshipProtocolRouteLocalAlphaFlow|TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation|TestAuthenticatedGameplayFailurePathsLocalAlphaFlow'

printf '%s\n' "vibit local alpha example client path"
printf '%s\n' "mode: source-first repository-local proof over existing WebSocket/Protobuf runtime handlers"
printf '%s\n' "path: local onboarding -> login -> bind connection -> protected inventory/presence/storage/friends -> logout -> rejected protected request"
printf '%s\n' "failure proof: missing, malformed, unknown, expired, revoked, route-protection, and redacted error behavior"
printf '%s\n' "routes: runtime.authentication.AuthenticateWithDeviceCredential, runtime.connection.BindConnection"
printf '%s\n' "routes: inventory.GrantItem, inventory.GetInventory, presence.GetPlayerPresence"
printf '%s\n' "routes: storage.PutOwnStorageObject, storage.GetOwnStorageObject, storage.ListOwnStorageObjects, storage.DeleteOwnStorageObject"
printf '%s\n' "routes: friends.SendFriendRequest, friends.GetFriendRelationshipStatus, friends.AcceptFriendRequest, friends.ListFriendRelationships"
printf '%s\n' "routes: friends.RemoveFriend, friends.BlockPlayer, friends.UnblockPlayer, friends.RejectFriendRequest"
printf '%s\n' "routes: runtime.authentication.LogoutAccessToken"
printf '%s\n' "redaction: raw credentials, access tokens, digests, verifier keys, DSNs, and transport metadata are not printed"
printf '%s\n' "scope: not an SDK, generated client library, hosted demo, release artifact, or direct compatibility surface"
printf '\n'

cd "$ROOT_DIR/runtime"
go test ./internal/platform/protocol/protobuf -run "$TEST_PATTERN" -v
