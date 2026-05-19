# Verification

Status: Verified

Required commands:

```bash
node -c tools/vibit
node tools/vibit check schemas --json
node tools/vibit check protocol --json
node tools/vibit check generated --json
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check change confirm-next-direction-after-first-message-connection-binding-gate --json
node tools/vibit check change define-first-message-connection-binding-implementation-gate --json
node tools/vibit check all --json
git diff --check
cd runtime && go test ./...
```

Notes:

- This change is standards/checking only.
- It intentionally does not add Go runtime behavior, Protobuf source, generated output, migrations, or dependencies.
