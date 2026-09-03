---
id: REQ-483
title: '[impact-critical] Review fix: Bound the architecture bundle-claim loop and restore --commit'
status: pending-heavy-testing
priority: now
created_at: 2026-09-01T11:51:27Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-418]
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
related: [REQ-418, REQ-420]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-418
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-03T21:44:31Z
  basis:
    - trivial short-circuit
claimed_at: 2026-09-03T21:43:37Z
route: A
dispatch_at: 2026-09-03T21:48:13Z
implementation_at: 2026-09-03T21:53:38Z
builder_handback_at: 2026-09-03T21:55:13Z
integration_at: 2026-09-03T21:57:28Z
testing_at: 2026-09-03T22:00:11Z
status_changed_at: 2026-09-03T22:00:11Z
commit: 3eb87519df19a14103f407159b9f6e753b51ca7b
write_set:
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go
---

# Bound the Architecture Bundle-Claim Loop and Restore --commit

## What

Make `architecture-report-preflight --publish` terminate with a typed finding when its
bundle claim fails for any reason other than "already exists", instead of spinning an
unbounded CPU-burning retry loop, and make its `--commit` mode work again. Both are
reproduced regressions introduced by REQ-418's sole remediation against the
pre-remediation build `a7c975c5`.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in
any UR sharing either root cause (the claim-retry loop or the hardcoded dry-run
passthrough). REQ-420 is adjacent for the `--commit` parity half, but an
`impact-critical` finding is never deferred behind a gate, and the loop is the same
code path, so both land here as one fix.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A direct implementation: pin both reproduced regressions, then restrict retry to collision identity and separate the dry-run preflight from the caller's commit mutation.
- [x] **[APPLY]:** Added deterministic RED tests and implemented the two control-flow corrections in exactly the declared architecture source and test files.
- [x] **[UNIFY]:** Reviewed both declared files and the complete branch diff; verified collision-only retry, typed path/cause evidence, commit passthrough, deterministic test seams, formatting, vet, focused tests, full module tests, and no debug artifacts.

## Requirements

- `architecture.go:156-164`: the sequence-increment retry must advance only on the
  claim's already-exists refusal; any other claim error (permissions, read-only
  filesystem, I/O) terminates with a typed finding naming the path and cause, matching
  the pre-remediation behavior of exiting with the underlying error.
- `architecture.go:144`: stop passing `dryRun=true` alongside the caller's `commit` —
  `--commit` currently fails in every argument position with
  `--dry-run and --commit cannot be combined`, naming an option the user never passed.
- Add committed tests for both: a non-existence claim failure (e.g. unwritable
  `reports/`) terminating with the typed finding, and a `--commit` publication
  committing exactly as the pre-remediation build does. Neither path has any coverage
  today.
- Preserve every rooted-publication and confinement protection the remediation added.

## Red-Green Proof

**RED prompt/case:** In a real-Git fixture, `chmod 555 reports` then run
`architecture-report-preflight --publish draft.html reports/run`; separately run the
same command with `--commit` on a writable fixture.
**Why RED now:** The publish loop was still running after 15s at 98.6% CPU and needed
SIGKILL (the `a7c975c5` build exits 3 with "permission denied" on the identical
fixture); `--commit` is rejected in every argument position.
**GREEN when:** The claim failure exits promptly with a typed finding, `--commit`
publishes and commits, and both behaviors are pinned by committed tests.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/runs/work-2026-08-31-165510/REQ-418-rereview.md` (findings N1 and N2)
while the run scratch survives; the durable record is the `## Review` section of
`do-work/archive/REQ-418-migrate-toolbox-absorb-audit-metrics.md`.

---
*Source: REQ-418 fresh re-review findings N1 (`impact-critical`) and N2 (same file, same fix surface).*

## Triage

**Route: A** - Simple

**Reasoning:** Both reproduced regressions name the exact architecture command sites, expected branch behavior, and caller-seam tests. The change is a small control-flow correction with no architectural decision left open.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go`

**Acceptance criteria:** Retry only exclusive-create collisions; return a typed path-and-cause finding for any other bundle-claim error; make caller-authored `--commit` publish and commit exactly the generated index; preserve rooted publication protections.

## Pre-Flight

**Git:** The wave baseline was clean at `c27d349a` after the three claims, estimates, run manifest, and briefs were committed.

**Tests:** Direct canonical fast gate passed and was recorded at the shared wave baseline before dispatch.

**Dependencies:** REQ-418 is completed; REQ-420 is adjacent but not a prerequisite.

## Implementation Summary

**Builder commit:** `f2ae022a49189acec1f9c934d0c85e83c42a52f1`

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go` (modified)

**What was done:** The validation-only transaction always uses dry-run without commit; the real transaction retains the caller's commit choice. Exclusive bundle claims now retry only `fs.ErrExist` and otherwise return a typed `ARCHITECTURE-BUNDLE-CLAIM-FAILED` finding naming the candidate and cause. A package-local claim seam makes the permission refusal deterministic in tests without weakening production rooting.

**Builder verification:** Both RED regressions turned GREEN; focused toolbox tests, full module tests, `go vet ./...`, formatting, and diff hygiene passed. Durable evidence is in `do-work/runs/work-2026-09-03-214500/REQ-483-handback.md`.

## Qualification

Passed — the exact `6a7e49d4..3eb87519` integration range contains only the two declared architecture files, all requirements trace to the implementation and tests, P-A-U is complete, and no `do-work/` path was committed by the builder.

## Testing

**Tests run:** the two focused regressions; `go test ./internal/toolboxcommands`; `go vet ./...`; full `go test ./...`; and direct `bash _dev/tests/maintainer-verify.sh`.

**Result:** All focused, module, vet, and fast canonical-gate checks passed. The direct fast gate started on merge `3eb87519`; while it ran, an unrelated report-only commit advanced main to `ed70391a`, so the green-gate record is conservatively stored at that descendant revision. The saved implementation range remains `6a7e49d4..3eb87519`.

**Red-green validation:** The old loop exceeded the 500 ms deadline on a deterministic permission refusal and the old commit path returned `GIT-INVALID-OPTIONS`; both added tests pass after the fix.

**Heavy boundary:** `architecture_test.go` matches `skills/do-work/tools/do-work-cli/**/*_test.go` from `maintainer-verify.sh --heavy-surfaces`; exact-revision heavy permission is therefore required before independent review and finalization.

## Open Questions

- [ ] Run `bash _dev/tests/maintainer-verify.sh --heavy` at `3eb87519df19a14103f407159b9f6e753b51ca7b`; did it exit 0?
  Recommended: Yes
  Also: No — report the failing lane
