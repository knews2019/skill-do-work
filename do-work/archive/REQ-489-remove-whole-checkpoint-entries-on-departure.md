---
id: REQ-489
title: '[impact-critical] Remove whole checkpoint entries when a REQ leaves working'
status: completed
created_at: 2026-09-01T19:46:57Z
user_request: UR-083
addendum_to: REQ-440
domain: backend
route: A
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
sweep: true
sweep_key: checkpoint-section-blind-line-editing
status_changed_at: 2026-09-01T21:05:20Z
claimed_at: 2026-09-02T13:45:32Z
dispatch_at: 2026-09-02T14:02:54Z
builder_handback_at: 2026-09-02T14:11:03Z
integration_at: 2026-09-02T14:11:43Z
review_at: 2026-09-02T14:27:05Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-02T13:47:59Z
  basis:
    - trivial short-circuit
completed_at: 2026-09-02T14:31:16Z
release_at: 2026-09-02T14:33:40Z
commit: 6e92e536
kb_status: promoted
kb_entry: REQ-489-remove-whole-checkpoint-entries-when-a-r.md
---

# Remove Whole Checkpoint Entries When a REQ Leaves Working

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the CLI prime, full touch-required lessons, crew rules, and bug-fix spec. Isolated the two faults: global line filtering orphaned continuations and substring heading lookup targeted inline prose.
- [x] **[APPLY]:** Added both captured regressions first, then implemented exact heading-line bounds and matching-entry continuation removal without changing foreign entries.
- [x] **[UNIFY]:** Reviewed both changed files, ran focused and full request-state/module tests plus `go vet`, `git diff --check`, and debug-artifact checks; the builder branch was clean after commit.

## What

When the canonical `complete` (and, by the same code path, `fail` and the other departures from `do-work/working/`) removes this checkout's entry from `do-work/CHECKPOINT.md`'s `## In Progress (interrupted)` list, it deletes only the `- REQ-NNN: ...` header line. The indented `Last known state:`, `Key files being modified:`, and `Known issues:` lines that Step 10 enrichment adds beneath it stay behind as an unattributed orphan block.

Remove the entire entry: the header line plus every following indented continuation line up to the next `- ` entry, blank line at the section boundary, or heading.

## Context

Found during REQ-440's archive on 2026-09-01: after `complete REQ-440` the checkpoint kept REQ-440's three detail lines with no header, and the same orphan block from REQ-418's earlier archive was still present above it. The orchestrator removed both by hand at Step 10. Root cause: `checkpointWithoutClaim` in `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` drops only lines matching `- REQ-NNN:` with the writer label and keeps every following indented line.

## Requirements

- Entry removal deletes the header and all of its indented continuation lines.
- A checkpoint whose entry has no continuation lines behaves exactly as today.
- Foreign-label entries and their continuation lines are untouched.
- Insertion (`appendSectionEntry`) and removal (`checkpointWithoutClaim`) locate `## In Progress (interrupted)` by a whole heading line, never by substring; a backticked or inline mention of the heading elsewhere in the checkpoint must not attract or strip an entry.

## Instances
- [ ] `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` `checkpointWithoutClaim`: removes only the `- REQ-NNN:` header line and orphans the indented detail lines (found by REQ-440 / UR-083)
- [ ] `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` `appendSectionEntry` (`strings.Index` on the heading text): the backticked mention of `## In Progress (interrupted)` in a Session Notes bullet of `do-work/CHECKPOINT.md` captured the live REQ-453 claim entry, leaving the real section empty so the claim reads as foreign to Crash Recovery and defeats REQ-498's entry-removal attribution (found by UR-097)

## Red-Green Proof
**RED prompt/case:** A requeststate test that builds a checkpoint with an enriched own-label entry (header plus three indented detail lines), runs the departure removal, and asserts the section contains none of the four lines.
**Why RED now:** Only the header line is removed; the three indented lines remain.
**GREEN when:** The test passes and an existing removal test for a bare one-line entry still passes.
**Second RED case (UR-097):** a checkpoint whose Session Notes bullet contains the backticked heading text followed later by the real `## In Progress (interrupted)` heading; `checkpointWithClaim` must place the new entry after the real heading, and `checkpointWithoutClaim` must leave the bullet untouched. Today the entry lands under Session Notes.
**Validation:** Discovered task from REQ-440; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions
- [x] I discovered this out-of-scope task while working on REQ-440: the canonical archive command leaves orphaned detail lines in the session checkpoint when it removes a finished REQ's in-progress entry. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

<!-- D-XX counter: none used. Next decision: D-01. -->

---
*Source: discovered task recorded during REQ-440 (work action Step 8).*


## Answer Notes

- 2026-09-01 - [ ] I discovered this out-of-scope task while working on REQ-440: the canonical archive command leaves orphaned detail lines in the session checkpoint when it removes a finished REQ's in-progress entry. Should I process this as a new task?: Confirmed: Yes, add to queue
> ```
> Add this checkpoint-cleanup fix to the queue for normal implementation. No scope is excluded.
> ```

## Triage

**Route: A** - Simple

**Reasoning:** The request names the exact request-state helper, its faulty line-oriented behavior, the required boundary cases, and the regression tests that prove the fix.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2409 tokens; matches canonical request-state checkpoint mutation, but this partially slugged satellite cannot be narrowed below the 2000-token budget. Read anyway under the prime's touch-conditional Lessons rule.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)

**What was done:** Checkpoint claim insertion and removal now locate the real `## In Progress (interrupted)` heading as a whole line. Departure removes the matching own-label header and its indented continuation block while preserving bare entries, foreign entries, section boundaries, and inline/backticked heading mentions.

## Root Cause

Checkpoint mutation treated the document as undifferentiated text: departure filtered one matching line globally, while insertion selected the first heading substring. The code had neither exact section bounds nor a representation of continuation-line ownership.

## Qualification

Passed — 2 files verified, 4 requirements traced, P-A-U confirmed. The merge range contains only the request-state helper and its focused regressions; both changes are substantive, exact-section bounded, and preserve foreign content.

## Testing

**Tests run:** `go test -count=1 ./internal/requeststate -run 'TestCheckpointClaim(RemovalIncludesIndentedContinuationLines|UsesWholeInProgressHeadingLine)$'`, `go test -count=1 ./internal/requeststate`, and `bash _dev/tests/maintainer-verify.sh`

**Result:** All focused and package tests passed on the merged tree. The canonical maintainer gate exited 0 after ShellCheck, formatting, contract, queue-kanban, strict JavaScript, vet, and uncached do-work-cli suites. Its optional strict browser lane reported the repository's normal no-browser skip.

**Red-green validation:**
- `TestCheckpointClaimRemovalIncludesIndentedContinuationLines`: ✗ before implementation — all three own-label continuation lines remained orphaned → ✓ after implementation — the whole own entry is removed and the foreign entry remains.
- `TestCheckpointClaimUsesWholeInProgressHeadingLine`: ✗ before implementation — inline/backticked heading prose attracted insertion and could be removed → ✓ after implementation — only the real heading bounds insertion/removal and Session Notes remain byte-preserved.

**New tests added:** the two focused regressions above in `state_apply_test.go`.

## Review

**Overall: 75%** | 2026-09-02T14:27:05Z

| Dimension | Score |
|-----------|-------|
| Requirements | 75% |
| Code Quality | 90% |
| Test Adequacy | 75% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `internal/cleanup.ownedCheckpointRemoval` remains an alternate header-only departure writer and can orphan the same enriched continuation block — `impact-user-visible` → REQ-502 created as the next `checkpoint-section-blind-line-editing` sweep.

**Minor findings:** 0 (report only)
**Acceptance:** Partial — canonical request-state departures pass both captured regressions, but cleanup's recovery mover still carries the same root cause.
**Suggested testing:** 1 item — extend cleanup's working-archive regression with enriched own and foreign entries.
**Follow-ups created:** REQ-502; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Exact-section RED cases exposed both the orphaned continuation block and the inline-heading attraction bug with a small two-file diff.

**What didn't:** Treating the canonical request-state helper as the only departure writer missed cleanup's independent checkpoint-removal implementation.

**Worth knowing:** A stored-format contract is not closed until every alternate writer of that format is swept; package-local green tests can still leave a repository-wide lifecycle path stale.

## Orientation

Canonical request-state departures now remove complete checkpoint entries and ignore inline heading mentions; the behavior lives in the do-work CLI request-state subsystem. The alternate cleanup mover is explicitly tracked by REQ-502. Map unchanged.
