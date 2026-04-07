# Test Strategy Action

> **Part of the do-work skill.** Designs a test strategy for a target scope — identifies what tests should exist, prioritized by risk reduction per effort. Read-only audit.

**Read-only** — this action does NOT modify any files. It produces a structured test strategy report only.

## Philosophy

Good test strategy is risk-driven, not coverage-driven. A codebase with 40% coverage focused on critical paths is more valuable than one with 90% coverage that tests getters and setters. This action answers: "Where should we invest testing effort for maximum bug prevention?"

## Input

`$TARGET` comes from `$ARGUMENTS` — a directory path, prime file reference, or both. If empty, defaults to the current working directory.

### Targeting Modes

Same scoping rules as code-review:
- **Prime file target**: `do work test-strategy prime-auth` — scope is files referenced by the prime
- **Directory target**: `do work test-strategy src/` — scope is all source files in the directory
- **Combined**: `do work test-strategy prime-auth src/utils/` — union of both
- **Default (no args)**: Search for prime files, ask user which scope, or fall back to cwd

## Steps

### Step 1: Resolve Scope and Detect Infrastructure

1. Resolve `$TARGET` to concrete files (same logic as code-review Step 1)
2. Load `prime-*.md` files for project context
3. **Detect test infrastructure**:
   - Test framework: Jest, Vitest, pytest, Go testing, RSpec, PHPUnit, JUnit, etc.
   - Test directories: `__tests__/`, `test/`, `tests/`, `spec/`
   - Test config: `jest.config.*`, `vitest.config.*`, `pytest.ini`, `.mocharc.*`, etc.
   - CI test commands: check `package.json` scripts, `Makefile`, CI config for how tests run
   - Coverage tools: check for coverage config, `.nycrc`, `coverage` in CI
4. **Map existing tests to source files** — for each source file in scope, find its test file(s) if any exist

### Step 2: Risk Assessment

Classify each module/file/component in scope by risk:

| Risk Factor | High Risk | Low Risk |
|-------------|-----------|----------|
| **Data integrity** | Writes to database, processes payments, modifies user accounts | Read-only display, static content |
| **User-facing** | Authentication, checkout, form submission, file upload | Internal tooling, admin dashboards |
| **Complexity** | Complex state machines, multi-step workflows, concurrency, recursive logic | Simple CRUD, pass-through functions |
| **Blast radius** | Shared utilities used by 10+ modules, core middleware, base classes | Leaf components, isolated helpers |
| **Change frequency** | Hot files (modified often) | Stable files (unchanged for months) |

For each module in scope, assign:
- **Risk level**: Critical / High / Medium / Low
- **Primary risk factor**: which factor(s) drove the rating (be specific)

### Step 3: Gap Analysis

For each module, assess current test state:

| State | Meaning |
|-------|---------|
| **Well-tested** | Meaningful tests exist covering happy path + error cases + edge cases |
| **Partially tested** | Some tests exist but gaps in error handling, edge cases, or branches |
| **Smoke only** | Tests exist but only cover the happy path — no error or edge coverage |
| **Untested** | No test file or test coverage at all |

Cross-reference risk level with test state to find the gaps that matter:

- **Critical risk + Untested** = urgent gap
- **Critical risk + Smoke only** = significant gap
- **High risk + Untested** = significant gap
- **Low risk + Untested** = acceptable gap (deprioritize)

### Step 4: Test Pyramid Design

Recommend the right test types for each gap, following the test pyramid:

| Layer | When to use | Characteristics |
|-------|------------|-----------------|
| **Unit tests** | Pure logic, data transformations, validation, calculations, utility functions | Fast, isolated, no I/O. The foundation — most tests should be here |
| **Integration tests** | Database queries, API endpoints, middleware chains, service interactions | Tests real I/O but controls the environment (test DB, mock external APIs) |
| **E2E tests** | Critical user flows (login, checkout, onboarding), multi-page workflows | Slow, brittle — use sparingly. Only for flows where the integration between systems is the risk |

**Anti-patterns to flag:**
- **Inverted pyramid**: More E2E tests than unit tests — slow CI, flaky results
- **Ice cream cone**: Manual testing > E2E > integration > unit — completely inverted
- **Hourglass**: Many unit + many E2E but no integration tests — misses the middle
- **100% coverage target**: Leads to testing implementation details rather than behavior

### Step 5: Flaky Test Prevention

For each recommended test, note flakiness risks and mitigations:

| Flakiness Source | Mitigation |
|-----------------|------------|
| **Time-dependent** | Inject clock/freeze time. Never assert on `Date.now()` or `time.time()` directly |
| **Order-dependent** | Each test sets up and tears down its own state. No shared mutable fixtures |
| **Network-dependent** | Mock external services. Use recorded responses (VCR pattern) for integration tests |
| **Race conditions** | Use deterministic waits (poll for condition), not `sleep(2)`. Set explicit timeouts |
| **Random data** | Use seeded random generators in tests. Log the seed for reproducibility |

### Step 6: Prioritized Recommendations

Rank all recommendations by **risk reduction per effort**:

**Risk reduction** = (risk level of module) x (severity of test gap)
**Effort** = estimated complexity of writing the tests (Trivial / Small / Medium / Large)

Priority order:
1. Critical risk + Untested + Trivial effort = **highest priority**
2. Critical risk + Untested + Small effort
3. High risk + Untested + Trivial effort
4. Critical risk + Smoke only + Small effort
5. ...continue down the matrix

Drop recommendations that are Large effort + Low risk — not worth mentioning.

### Step 7: CI Integration Check

Evaluate the CI pipeline for test quality gates:

| Check | What to verify |
|-------|---------------|
| **Execution time** | Total test suite runs in under 10 minutes? If not, recommend parallelization or splitting |
| **Failure visibility** | Test failures block the PR/merge? Or are they advisory-only? |
| **Coverage reporting** | Coverage measured and visible? Used as a trend, not a gate? |
| **Flaky test handling** | Are flaky tests quarantined? Or do they cause retry-storms in CI? |

## Output Format

```markdown
# Test Strategy Report

**Scope**: {description}
**Files analyzed**: {N} files ({total lines} lines)
**Test infrastructure**: {framework, config, CI integration}
**Date**: {today}

## Summary

{2-3 sentences — overall testing health and the single highest-priority recommendation.}

## Risk Map

| Module | Risk | Primary Factor | Test State | Gap Severity |
|--------|------|----------------|------------|-------------|
| `src/auth/login.ts` | Critical | Authentication, user-facing | Smoke only | Significant |
| `src/utils/format.ts` | Low | Pure utility, stable | Untested | Acceptable |

## Recommended Tests

### Priority 1: {description}

**Module**: `src/auth/login.ts`
**Risk**: Critical | **Current state**: Smoke only
**Test type**: Unit + Integration
**What to test**:
- Invalid credentials return 401 with appropriate error
- Account lockout after N failed attempts
- Session token generation uses crypto-random
- Expired sessions are rejected
**Effort**: Small
**Flakiness risk**: Time-dependent (session expiry) — inject clock

### Priority 2: {description}
...

## Test Pyramid Assessment

**Current shape**: {Pyramid / Hourglass / Inverted / Ice cream cone / None}
**Recommended adjustment**: {what to add or reduce}

| Layer | Current Count | Recommended Direction |
|-------|--------------|----------------------|
| Unit | {N} | {Add more / Sufficient / Reduce} |
| Integration | {N} | {Add more / Sufficient / Reduce} |
| E2E | {N} | {Add more / Sufficient / Reduce} |

## CI Health

| Check | Status | Recommendation |
|-------|--------|----------------|
| Execution time | {N minutes} | {OK / Split suites / Parallelize} |
| Failure blocking | {Blocking / Advisory} | {OK / Make blocking} |
| Coverage trend | {Tracked / Not tracked} | {OK / Add coverage reporting} |

## Strengths

{Give credit — which areas are well-tested and should be maintained as-is.}

## Recommended Next Steps

1. {Highest priority action}
2. {Second priority}
3. {Third priority}

> To act on these findings:
>   do work capture request: [describe the test to write]
>   do work run
```

## Rules

- **Do NOT modify any files.** This action is read-only. Report only.
- **Risk-driven, not coverage-driven.** Never recommend "increase coverage to X%." Recommend specific tests for specific risks.
- **Be honest about what's acceptable.** Low-risk utility functions with no tests is fine. Say so. Don't pad the report.
- **Respect existing infrastructure.** Recommend tests that fit the project's existing framework and patterns. Don't suggest switching from Jest to Vitest just because you prefer it.
- **Account for maintenance cost.** A test that's expensive to maintain (brittle selectors, complex setup, mock-heavy) needs to justify its existence with high risk reduction.
- **Proportional depth.** 10 files get individual risk assessment. 100+ files get module-level grouping. State your approach.
- **Skip vendored and generated files.** Same exclusions as quick-wins.
