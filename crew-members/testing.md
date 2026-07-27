# The Verifier — Testing Crew Member

<!-- JIT_CONTEXT: This file is loaded by the AI agent when working on test-heavy tasks (domain: testing), when the REQ has tdd: true, or when the test failure loop (Step 6.5) exceeds 1 attempt alongside debugging.md. Provides structured guidance on test strategy, framework selection, and common pitfalls. -->

## Core Principle: Tests Prove Behavior, Not Implementation

A good test breaks when the behavior changes, not when the code is refactored. Test the contract (inputs → outputs, side effects, error conditions), not the internal wiring.

## Test Framework Detection

Before writing tests, identify the project's existing test setup. Don't guess — check:

| Indicator | Framework | Runner |
|-----------|-----------|--------|
| `pytest.ini`, `pyproject.toml [tool.pytest]`, `conftest.py` | pytest | `pytest` / `uv run pytest` |
| `jest.config.*`, `"jest"` in package.json | Jest | `npx jest` / `npm test` |
| `vitest.config.*`, `"vitest"` in package.json | Vitest | `npx vitest` |
| `playwright.config.*` | Playwright | `npx playwright test` |
| `cypress.config.*`, `cypress/` dir | Cypress | `npx cypress run` |
| `Cargo.toml` | Rust built-in + cargo-nextest | `cargo test` / `cargo nextest run` |
| `*_test.go` files | Go built-in | `go test ./...` |
| `spec/*_spec.rb`, `.rspec`, `Gemfile` with rspec | RSpec | `bundle exec rspec` |

**Rule:** Use the project's existing framework. Never introduce a second test framework unless the REQ explicitly requests it.

## Opinions

- **Default test type:** if the REQ doesn't specify, write unit tests. Escalate to integration only when the behavior involves a real I/O boundary. Reserve E2E for critical user journeys — too expensive to maintain for edge cases.
- **Test at the caller seam.** A test that exercises a unit only through a harness production code never uses can pass while the real integration is broken. At least one test per behavior should call through the seam the real caller hits — the public API, the route handler, the CLI entry point.
- **Fixtures must be production-faithful.** Wrong field names, simplified nesting, or fabricated enum values let a green suite certify broken code. Derive fixtures from real captured payloads or the schema/type definitions, and validate hand-written ones against them.
- **More than 3 mocks needed for one test** is a design-smell signal (too many dependencies) — note it, don't refactor unless the REQ asks for it.
- **Encounter a flaky test during implementation? Fix the flakiness before proceeding.** Don't re-run and hope.

## Red-Green Workflow (TDD Requests)

When the REQ has `tdd: true`:

1. **Red:** Write a test that captures the expected behavior. Run it. Verify it fails for the right reason (not a syntax error or import failure — the assertion itself must fail).
2. **Green:** Write the minimum code to make the test pass. No more.
3. **Refactor:** Clean up the implementation without changing behavior. Tests must still pass.

**Evidence:** After the Green step, record the test output showing the transition from red to green. This is the Red-Green Proof referenced in the REQ.

## Anti-Patterns

- **Asserting on error messages:** Brittle — messages change. Assert on error types, status codes, or structured error fields instead.
- **Test-per-method symmetry:** Mirroring the source file 1:1 (one test per method) tests method inventory, not behavior. Some methods need 5 tests; some need zero.
- **Catch-all assertions:** `expect(result).toBeTruthy()` or `assert response` without checking specific values passes on wrong results and catches nothing.
- **Ignoring test output:** Re-running a failing test without reading the failure message. The assertion diff usually points directly at the bug.
