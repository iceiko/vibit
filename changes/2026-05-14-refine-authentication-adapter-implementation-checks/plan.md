# Plan

1. Update `tools/vibit` so selected login/token checks exclude declared authentication PostgreSQL adapter files from broad runtime-vocabulary bans.
2. Add adapter-specific static checks that still forbid token generation, token validation, verifier comparison, bearer parsing, transport/protocol behavior, transaction control, and major dependencies.
3. Add manifest markers showing adapter checks are refined while adapter implementation remains absent.
4. Update English and Simplified Chinese standards.
5. Mark `W-0083` completed and make `W-0084` next ready.
6. Run focused and full verification.
