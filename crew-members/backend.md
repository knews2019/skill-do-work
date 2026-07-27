# The Engineer — Backend Crew Member

<!-- JIT_CONTEXT: This file is loaded by the AI agent only when working on backend-related tasks. Keep rules scoped and concise to minimize token usage. -->

## Security

Authentication, authorization, secrets handling, and rate limiting are owned by `crew-members/security.md` — it loads automatically when the surface touches those categories. Don't restate its checklist here.

## Opinions

- **API versioning:** follow the project's existing convention (URL prefix, header-based, query param). If none exists and you're creating a new API surface, use URL prefix versioning (`/v1/`) and record the choice as a Decision (D-XX).
- **Caching:** for data that changes infrequently (config, feature flags, permissions), don't add caching the REQ didn't request — note the opportunity in Discovered Tasks instead.

## Quality Checks

Before marking UNIFY complete, verify:

| Criterion | What to check |
|-----------|---------------|
| Handles invalid input | Malformed requests return 400, not 500 |
| Auth enforced | Protected routes reject unauthenticated/unauthorized requests |
| No data leaks | Error responses don't expose internal details |
| Idempotency | Safe methods (GET) have no side effects; writes handle retries |
| Existing tests pass | No regressions in adjacent endpoints or services |
