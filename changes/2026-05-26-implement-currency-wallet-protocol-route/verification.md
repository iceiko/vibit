# Verification

Status: Verified

Commands:

```sh
buf generate
buf lint
cd runtime && go test ./internal/platform/protocol/protobuf -run TestCurrencyWalletProtocolRouteLocalAlphaFlow -count=1
cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestCurrency|TestAuthenticated|TestStorage|TestFriends|TestFrameHandler' -count=1
cd runtime && go test ./internal/app/currency ./internal/app/bootstrap ./internal/platform/protocol/protobuf -count=1
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.currency_wallet_protocol_route_implementation
node tools/vibit check change implement-currency-wallet-protocol-route --json
node tools/vibit check generated --json
node tools/vibit check protocol --json
node tools/vibit check module currency --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check all --json
cd runtime && go test ./...
git diff --check
rg -n 'ghp_[A-Za-z0-9_]{20,}|github_pat_[A-Za-z0-9_]+' --glob '!node_modules/**' --glob '!.git/**' --glob '!.vibit.local.env'
```

Focused TDD evidence:

- The currency wallet local alpha E2E test first failed with `AUTHENTICATION_TOKEN_MALFORMED` for a zero-field proto3 request encoded to empty bytes.
- After the authenticated payload handling fix, the test failed with `ROUTE_NOT_FOUND`, confirming route registration was the remaining gap.
- After adding currency service/repository fixture wiring and route registration, the focused E2E test passed.
