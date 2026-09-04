---
id: REQ-502
title: 'Review fix: Remove enriched checkpoint entries in cleanup mover'
status: claimed
domain: backend
created_at: 2026-09-02T14:26:49Z
user_request: UR-083
addendum_to: REQ-489
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-03T20:59:01Z
  basis:
    - trivial short-circuit
tdd: true
suggested_spec: bug-fix
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#alternate-writer-contract-drift]
sweep: true
sweep_key: checkpoint-section-blind-line-editing
dispatch_at: 2026-09-03T20:59:01Z
builder_handback_at: 2026-09-03T21:02:28Z
integration_at: 2026-09-03T21:05:36Z
status_changed_at: 2026-09-04T12:36:12Z
commit: ed692757dfc642f3ad34b171dde9f6490c857beb
heavy_verified_at: 2026-09-04T12:36:12Z
heavy_verified_revision: c0d8ce1cb44cc1830b167214c018d76ba87baffc
claimed_at: 2026-09-04T13:32:07Z
---

# Review Fix: Remove Enriched Checkpoint Entries in Cleanup Mover

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reproduce at the cleanup seam, then expose the existing request-state removal primitive instead of copying its stored-format parsing rules.
- [x] **[APPLY]:** Added the focused regression first, observed the orphaned continuation lines, then delegated cleanup removal to the shared whole-entry primitive.
- [x] **[UNIFY]:** Reviewed all three changed files; ran gofmt, focused and package tests, `go vet`, `git diff --check`, and a debug-artifact scan.

## What

Make every cleanup-package departure from `do-work/working/` remove the complete matching checkpoint entry, not only its header line. Done means the cleanup mover and canonical request-state mover share the same exact-section, whole-entry semantics, so this root cause cannot recur through an alternate departure path.

## Context

Independent review of REQ-489 found that `internal/cleanup.ownedCheckpointRemoval` still globally filters only the matching `- REQ-NNN:` header and leaves Step 10's indented detail lines orphaned. The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this root cause; REQ-489 itself is already claimed and therefore is not a fold candidate.

## Requirements

- Remove the matching own-label cleanup claim header and all immediately following nonblank indented continuation lines.
- Preserve foreign-label entries and their continuation bytes exactly.
- Locate the real `## In Progress (interrupted)` section by a whole heading line; inline or backticked mentions elsewhere are not entries.
- Reuse or align with the request-state semantics without creating a second drifting definition.

## Instances

- [ ] `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` `ownedCheckpointRemoval`: cleanup archives a terminal working REQ but leaves its enriched checkpoint continuation block unattributed (found by REQ-489 / UR-083)

## Red-Green Proof

**RED prompt/case:** Extend `TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry` with an enriched own-label entry and an enriched foreign-label entry, then build the cleanup plan.
**Why RED now:** `ownedCheckpointRemoval` skips only the own header line and appends all following continuation lines.
**GREEN when:** Cleanup removes the complete own entry, preserves the foreign header and continuation bytes exactly, and ignores inline/backticked heading mentions.
**Validation:** Independent review finding from REQ-489; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---
*Source: Important review finding from REQ-489; folded as the next same-root sweep after the claimed source REQ.*

## Triage

**Route: A** - Simple

**Reasoning:** The request names the exact cleanup helper, the incorrect header-only behavior, and the focused regression that proves whole-entry removal while preserving foreign bytes.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 — DECIDE & STATE:** Export the narrow writer-labelled removal primitive from `requeststate` and let `cleanup` import it. This keeps the stored-format rule in its existing owner and avoids a new package for one reuse site.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)

**What was done:** Cleanup now calls request-state's exact-section, whole-entry checkpoint remover. Its regression covers enriched own and foreign entries plus an inline heading mention.

## Root Cause

Cleanup retained an older global header-line filter after request-state moved to section-bounded whole-entry removal. The alternate writer therefore shared the checkpoint format but not its mutation semantics.

## Qualification

Passed — 3 files verified, 4 requirements traced, P-A-U confirmed. The merge range contains only the shared request-state helper, cleanup's delegation, and the cleanup seam regression; no debug artifacts or hollow changes were found.

## Testing

**Tests run:** `go test -count=1 ./internal/cleanup -run '^TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry$'`; `go test -count=1 ./internal/cleanup ./internal/requeststate`; `go vet ./internal/cleanup ./internal/requeststate`; `bash _dev/tests/maintainer-verify.sh`

**Result:** All focused, package, vet, and fast canonical-gate checks passed on merge `ed692757dfc642f3ad34b171dde9f6490c857beb`. Focused files completed below the 30-second file budget; the fast gate exited 0 and recorded this revision.

**Red-green validation:**
- `TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry`: failed before implementation because the owned continuation lines remained, then passed after cleanup delegated to the whole-entry remover.

**Existing tests updated:**
- `cleanup_plan_test.go` now proves whole-entry deletion, foreign-byte preservation, and real-heading selection at the cleanup seam.

*Verified by work action*

## Heavy Verification Plan

- `mode`: `historical-revalidation`
- `source_ranges`: `46d5fd48421062236bb4218a12e8f37f20098caf..ed692757dfc642f3ad34b171dde9f6490c857beb`
- `manifest_revision`: `c0d8ce1cb44cc1830b167214c018d76ba87baffc`
- `execution_revision`: `c0d8ce1cb44cc1830b167214c018d76ba87baffc`
- `manifest_path`: `_dev/tests/heavy-lanes.json`
- `forced_all`: `false`
- `uncertain`: `false`
- `uncovered_paths`: `[]`
- `changed_paths`:
  - `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go`
  - `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go`
  - `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`
- `selected_lanes`:
  - `do-work-cli-integrations`: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
  - `staged-skills`: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
  - `updater`: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
  - `installer`: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Open Questions

- [x] Run the selected heavy lane commands at execution revision `c0d8ce1cb44cc1830b167214c018d76ba87baffc`; did every command exit 0? → Confirmed: Yes
  Recommended: Yes
  Also: No — report the failing lane


## Answer Notes

- 2026-09-03 - [ ] Run `bash _dev/tests/maintainer-verify.sh --heavy` at `ed692757dfc642f3ad34b171dde9f6490c857beb`; did it exit 0?: No, exit 1: staged-skills-contract.sh took 35s and update-script-behavior.sh took 59s, exceeding the under-30s test-file budget
- 2026-09-04 - [ ] Run the selected heavy lane commands at execution revision `c0d8ce1cb44cc1830b167214c018d76ba87baffc`; did every command exit 0?: Confirmed: Yes

## Heavy Verification Result

Target revision: `ed692757dfc642f3ad34b171dde9f6490c857beb`
Execution revision: `c0d8ce1cb44cc1830b167214c018d76ba87baffc`

- do-work-cli-integrations: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`
