# Impact

## Summary

This change defines the Pitaya-aligned distributed runtime vocabulary reactivation gate. It moves the repository closer to Pitaya by making the future architecture vocabulary explicit, checkable, and bounded.

## Added

- Gate standard and Simplified Chinese translation.
- ADR-0154.
- Repository check rule registration.
- Change artifacts and conversation memory.
- Manifest updates completing W-0246 and opening W-0247.

## Not Added

- No runtime behavior.
- No protocol shape.
- No Protobuf source.
- No generated output.
- No repository or PostgreSQL adapter changes.
- No migrations.
- No dependencies.
- No startup wiring.
- No distributed runtime implementation.
- No direct Nakama/Pitaya API compatibility.

## Risk

The main risk is vocabulary being mistaken for implementation permission. The gate reduces that risk by recording explicit allowed vocabulary, forbidden use, current single-process mapping, and stop conditions.
