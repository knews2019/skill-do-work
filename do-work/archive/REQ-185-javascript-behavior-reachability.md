---
id: REQ-185
title: JavaScript behavior probes can all skip while the board suite passes
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T10:53:45Z
completed_at: 2026-08-15T11:22:22Z
user_request: UR-041
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-182, REQ-183, REQ-184, REQ-186, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [skills/do-work-board/tools/queue-kanban/generate_test.go, skills/do-work-board/tools/queue-kanban/web/board.js]
route: B
kb_status: pending
kb_entry:
---

# JavaScript Behavior Probes Can All Skip While the Board Suite Passes

## What

Add an explicit maintainer-strict JavaScript behavior lane so the board suite cannot report success after all four incident-sensitive Node probes skip, and convert remaining state-transition claims from source-token checks to executable behavior.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Centralize Node discovery/execution behind one helper that counts attempted JavaScript probes and preserves ordinary skips. Add a subprocess-backed maintainer entrypoint whose strict child fails after the run when zero probes execute, then convert recent-window, By-UR, empty-state, and confirmed testing-update transition claims into executable Node behavior while retaining narrow token checks only for generated call-site wiring.
- [x] **[APPLY]:** Added the zero-probe RED regression, centralized counted Node probe execution, a maintainer-strict subprocess lane, stable behavior-test naming, executable state-transition probes, and the two narrow production transition helpers within the declared two-file scope.
- [x] **[UNIFY]:** Reviewed the full `generate_test.go` and `web/board.js` diff. Verified the strict and ordinary no-Node boundaries, all seven Node probes, production helper wiring, full Go tests, `go vet ./...`, both repository contract suites, and `git diff --check`; no debug artifacts remain.

## Why

Four incident-sensitive board JavaScript tests independently call `t.Skip` when Node is absent. `go test` can therefore pass after executing zero JavaScript behavior probes, while some transition claims rely on source-token slicing and string checks against the third-ranked hotspot.

## Context

- Audit priority: P2; impact 3; effort normal.
- Root-cause key: `maintainer-js-behavior-reachability`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 5.
- Reproduce: `cd skills/do-work-board/tools/queue-kanban && probe_bin="$(mktemp -d /tmp/go-without-node.XXXXXX)" && ln -s "$(command -v go)" "$probe_bin/go" && PATH="$probe_bin:/usr/bin:/bin" "$probe_bin/go" test -v -run 'TestDrawerHeadingDeduplicationBehavior|TestByUserRequestLensEmptyStateBehavior|TestByUserRequestLensDefaultScopeUsesScopeOnlyEmptyState|TestTestingDoneWindowIsViewSpecific' .`

## Detailed Requirements

- Centralize Node capability detection for the board's JavaScript behavior probes.
- Add a maintainer-strict mode that exits nonzero if the behavior lane executes zero probes.
- Convert state-transition claims currently proven by `sliceBalancedBlockAfter` or `strings.Contains(indexHtml...)` into executable Node probes.
- Retain token checks only for generated-asset wiring, where token presence is the behavior being asserted.
- Preserve documented skip behavior for ordinary consumer/package tests when Node is unavailable.
- Reuse the no-Node replay as the regression case for the strict lane.

## Constraints

- Do not require Node for normal board use or ordinary package tests.
- Do not add a browser framework.
- Lock-in limit: zero successful maintainer runs with zero JavaScript behavior probes.

## Dependencies

None. REQ-187 may later invoke this strict lane, but this REQ must define its own complete behavior first. The `generate_test.go` overlap with REQ-183 is declared.

## Builder Guidance

Firm intent. Keep the strictness boundary explicit and maintainer-only; executable state behavior matters more than source-shape assertions.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Put only the Go binary plus `/usr/bin:/bin` on PATH and run `TestDrawerHeadingDeduplicationBehavior`, both user-request-lens empty-state tests, and `TestTestingDoneWindowIsViewSpecific`.
**Why RED now:** All four tests skip when Node is absent, but the selected `go test` command exits zero after running no JavaScript behavior.
**GREEN when:** The same no-Node replay exits nonzero in maintainer-strict mode, ordinary package tests retain documented skip behavior, and state transitions run as executable Node probes when Node is present.
**Validation:** Confirmed by the user during verification on 2026-08-15.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 05, labeled P2, impact 3, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 5 in the canonical audit for the complete batch constraints and validated evidence record.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*

## Triage

**Route B** — the behavior and two-file boundary are clear, but the existing Node helpers, four skip sites, source-token assertions, and maintainer entrypoint shape must be traced before the builder can declare the exact strict-lane scope.

## Exploration

- Four incident-sensitive tests independently call `exec.LookPath("node")` and `t.Skip`: drawer heading de-duplication, pure By-UR empty-state copy, default-scope By-UR caller behavior, and testing done-window copy.
- With Node hidden, the selected four-test replay exits 0 after four skips, so ordinary `go test` has no evidence that JavaScript behavior ran.
- `TestByUserRequestLensCountsRecentlyDoneAsActive`, `TestRecentlyDoneWindowHandlerRefreshesUserRequestLens`, `TestByUserRequestLensEmptyStateNamesWindow`, and `TestTestingStatusUpdateInvalidatesUserRequestLens` still use source slicing/token presence to claim runtime predicate or transition behavior.
- `sliceBalancedBlockAfter` remains useful for extracting real production functions into Node probes and for CSS/generated wiring checks; the defect is using source shape as the behavior proof.
- `board.js` has two narrow refactoring seams: apply a recent-window state change from the click handler, and apply a server-confirmed testing-status transition from the fetch success callback. Both can be executed under Node with small state/render stubs while keeping one call-site token check as asset wiring.
- A `TestMain` post-run count under an internal strict marker proves the lane executed at least one probe. A canonical parent test re-execs the current test binary with a stable JavaScript-behavior prefix; ordinary runs keep skip semantics and never inherit the strict marker.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — centralized Node probe runner/count, strict subprocess entrypoint/no-Node regression, stable behavior-test prefix, and executable transition probes
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modify) — extract narrow recent-window and confirmed-testing transition helpers used by existing production call sites

**Files I will NOT touch:** ordinary board runtime requirements, browser frameworks, Go production code, package manifests, or maintainer orchestration owned by later REQs.

**Acceptance criteria (restated from REQ):**
- [ ] Maintainer-strict JavaScript behavior mode exits nonzero when zero probes execute.
- [ ] Ordinary package tests still skip JavaScript probes and pass when Node is unavailable.
- [ ] Node capability detection/execution is centralized rather than repeated in four tests.
- [ ] Runtime state-transition claims execute production JavaScript; token checks remain only where wiring/presence is itself the contract.
- [ ] No Node dependency is added to normal board use and no browser framework is introduced.

## Pre-Flight

**Git:** ⚠ Four pre-existing edits under `do-work/queue/` (REQ-189–192) belong to other work and must remain unstaged; both scoped files are clean.
**Tests baseline:** ✓ Node v22.23.2 is available and `go test -count=1 ./...` passes (5.557s); the separate no-Node replay was confirmed to exit 0 after four skips.
**Dependencies:** ✓ Go and Node are available for the executable behavior lane.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — centralizes counted Node execution, enforces a strict zero-probe failure in `TestMain`, exposes the canonical maintainer lane, preserves ordinary skips, and replaces source-shape transition claims with seven executable JavaScript probes
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified) — extracts and wires narrow helpers for recent-window selection and server-confirmed testing-status transitions so production state changes can be executed directly in Node

**Behavior:** Maintainers can select one stable strict test entrypoint that fails if no JavaScript behavior probe actually starts. Ordinary package tests still skip when Node is unavailable, while Node-capable runs execute the production predicates, empty-state decisions, recent-window refresh, testing view copy, and confirmed testing transition.

## Testing

**RED:** With Node hidden, the new strict regression initially observed the behavior-test child report `testing: warning: no tests to run` and `PASS`; zero JavaScript probes still produced exit 0.

**GREEN:**
- `go test -count=1 -run '^TestMaintainerStrictJavaScriptBehaviorLaneRejectsZeroProbes$' .` — PASS; the nested no-Node strict child exits nonzero with the stable zero-probe diagnostic
- `go test -count=1 -run '^TestMaintainerStrictJavaScriptBehaviorLane$' -v .` — PASS
- `go test -count=1 -run '^TestJavaScriptBehavior' -v .` — PASS; all seven probes execute under Node
- Ordinary no-Node behavior-prefix replay — PASS with all seven probes skipped
- `go test -count=1 ./...` — PASS
- `go vet ./...` — PASS
- `bash _dev/tests/contract-regressions.sh` — PASS
- `bash _dev/tests/shipped-package-reference-contract.sh` — PASS
- `git diff --check` — PASS

## Qualification

- **Scope:** PASS — `scope-drift.sh` reports that the two-file Implementation Summary exactly matches the declared Scope; unrelated REQ-189–192 edits remain excluded.
- **Mechanical checks:** PASS — `qualify.sh` found both implementation files in the diff, all P-A-U phases complete, and no debug artifacts.
- **Substance and traceability:** PASS — the stable maintainer entrypoint re-executes only `TestJavaScriptBehavior*` tests under a strict marker, and `TestMain` rejects an otherwise successful child that attempted zero Node probes.
- **Wiring/data flow:** PASS — all seven behavior tests use the centralized counted runner; the two extracted production helpers remain wired at the original click and fetch-success call sites and are executed directly in Node with state/render observations.

## Review

**Result:** Approve — Acceptance: Pass
**Overall score:** 99%

- **Requirements (100%):** Centralized Node execution, strict zero-probe failure, ordinary skips, executable state claims, and the no-framework/runtime constraints are all delivered.
- **Code quality (98%):** The maintainer boundary is explicit and the two production extractions preserve the original sequencing without widening the runtime surface.
- **Test adequacy (100%):** All seven probes run through one counted lane; the no-Node strict/ordinary split, production recent-window caller, and visible/hidden confirmed-testing branches are covered.
- **Scope (100%):** Exactly the two declared files changed.

**Important findings:** None.
**Minor findings:** None.
**Explicit remediation:** The initial review found that recent-terminal lens population and hidden-lens invalidation could survive mutations. The single remediation added production-caller and hidden-branch observations; both targeted deletions failed before restoration, and re-review confirmed the finding fully closed.

## Lessons Learned

- An optional-tool test lane needs two separate contracts: ordinary consumers may skip unavailable probes, while the maintainer entrypoint must count attempted behavior and reject an otherwise green zero-probe run.
- Executing pure helpers is not enough when the regression lives in caller composition or a hidden-state branch. Mutation-resistant coverage must observe the production caller and each cache/render branch whose transition is part of the claim.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.
