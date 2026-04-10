# Spec: API Endpoint

> Specification template for building API endpoints — route definition, handler logic, validation, and tests.

## Output Structure

- **Route definition** — URL pattern, HTTP method, middleware chain
- **Handler** — request parsing, business logic, response formatting
- **Validation** — input schema (params, query, body), type coercion
- **Error responses** — structured error format, appropriate HTTP status codes
- **Tests** — unit tests for handler logic, integration tests for the full request cycle

## Quality Standards

- All user input validated before reaching business logic
- Auth/authorization checks applied via middleware, not inline in the handler
- Rate limiting awareness — note if the endpoint needs throttling (high-frequency, public-facing, or resource-intensive)
- Proper HTTP status codes: 200/201 for success, 400 for bad input, 401/403 for auth, 404 for missing resources, 409 for conflicts, 500 for server errors
- Error format consistent with existing API endpoints in the project
- Response shape matches existing conventions (envelope pattern, direct data, pagination structure)
- No N+1 queries — if the handler fetches related data, use eager loading or batch queries

## Implementation Checklist

1. Schema/types — define request and response shapes
2. Validation — input validation rules, type coercion
3. Handler logic — core business logic, data access
4. Error handling — catch and format errors consistently
5. Tests — unit tests for logic, integration tests for the full endpoint
6. Documentation — update API docs or OpenAPI spec if the project maintains one

## Evolution Path

- **Simple**: Basic CRUD — single resource, standard HTTP methods, minimal validation
- **Medium**: Pagination, filtering, sorting, query parameter parsing, partial updates
- **Complex**: Real-time/streaming responses, multi-resource transactions, optimistic concurrency, webhook integration

## Common Pitfalls

- Missing input validation — trusting client data leads to injection and corruption
- Inconsistent error format — some endpoints return `{ error: "msg" }`, others return `{ message: "msg" }`. Match the project's existing pattern.
- N+1 queries — fetching related records in a loop instead of a single batch query
- Missing auth middleware — handler works in tests but returns 401 in production because middleware wasn't wired
- Returning 200 for errors — use the right status code; clients depend on it for control flow
- Not handling empty results — return `[]` or `null` explicitly, not an unhandled exception
