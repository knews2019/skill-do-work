---
id: REQ-185
title: JavaScript behavior probes can all skip while the board suite passes
status: pending
created_at: 2026-08-15T07:13:20Z
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
---

# JavaScript Behavior Probes Can All Skip While the Board Suite Passes

## What

Add an explicit maintainer-strict JavaScript behavior lane so the board suite cannot report success after all four incident-sensitive Node probes skip, and convert remaining state-transition claims from source-token checks to executable behavior.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
