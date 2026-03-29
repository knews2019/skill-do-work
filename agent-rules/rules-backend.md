# Backend Rules

<!-- JIT_CONTEXT: This file is loaded by the AI agent only when working on backend-related tasks. Keep rules scoped and concise to minimize token usage. -->

## Implementation Patterns

### API Design
- Follow the project's existing API conventions (REST, GraphQL, RPC).
- Consistent error response format across endpoints.
- Input validation at the boundary — never trust client data.
- Use appropriate HTTP status codes (don't 200 everything).

### Data Layer
- Follow the project's existing ORM/query patterns.
- Migrations for schema changes — never modify production schemas directly.
- Transactions for multi-step writes that must succeed or fail together.
- Parameterized queries — no string interpolation for SQL.

### Security Baseline
- Authentication checks on every protected endpoint.
- Authorization: verify the user can access the specific resource, not just that they're logged in.
- No secrets in code, logs, or error responses.
- Rate limiting awareness — note if an endpoint needs it and doesn't have it.

### Error Handling
- Catch errors at the boundary (route handler / controller), not deep in business logic.
- Log with enough context to debug (request ID, user ID, action) but never log sensitive data (passwords, tokens, PII).
- Return structured errors to clients; keep stack traces server-side.

## Quality Checks

Before marking UNIFY complete, verify:

| Criterion | What to check |
|-----------|---------------|
| Handles invalid input | Malformed requests return 400, not 500 |
| Auth enforced | Protected routes reject unauthenticated/unauthorized requests |
| No data leaks | Error responses don't expose internal details |
| Idempotency | Safe methods (GET) have no side effects; writes handle retries |
| Existing tests pass | No regressions in adjacent endpoints or services |

## Scope Discipline

- Do not refactor unrelated endpoints while fixing a bug in one.
- Do not change database schemas beyond what the REQ requires.
- Do not introduce new dependencies for functionality the existing stack already provides.
