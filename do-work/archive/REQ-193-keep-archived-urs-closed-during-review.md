---
id: REQ-193
title: Keep archived URs closed during standalone review
status: completed
created_at: 2026-08-15T09:12:04Z
claimed_at: 2026-08-15T12:37:36Z
completed_at: 2026-08-15T13:07:45Z
user_request: UR-043
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-194]
batch: closed-ur-documentation-hardening
write_set: [skills/do-work/actions/review-work.md, skills/do-work/actions/work.md, _dev/tests/contract-regressions.sh]
route: B
kb_status: pending
kb_entry:
---

# Keep Archived URs Closed During Standalone Review

## What

Co-locate the archived-UR invariant with standalone review's archived-input read so a reviewer can create a follow-up without moving or reopening the already-closed User Request folder. Make the downstream archive instruction explicit enough that the completed follow-up returns to that folder in place.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Co-locate a context-only/stays-closed authority boundary with standalone review's archived-input fallback; add a success-archive override for review-generated REQs whose UR folder is already archived; and lock both seams with narrow, non-vacuous contract-regression extracts before changing the documentation.
- [x] **[APPLY]:** Added the contract probes first and captured the missing-clause RED, then co-located the mode-independent archived-input authority boundary in review and the already-archived review-follow-up success override in work, within the exact three-file Scope.
- [x] **[UNIFY]:** Reviewed the final three-file diff after correcting the ordinary-vs-closed orchestrated-review distinction. Verified `review-work.md` keys the boundary to archived-fallback use and preserves producer/metadata contracts; `work.md` requires both the review marker and existing archived folder before its in-place override; the aggregate assertions independently pin every clause. Bash syntax, focused probes, the full aggregate, canonical maintainer gate, and `git diff --check` pass with no debug artifacts.

## Why

The invariant already exists in abandon and capture, but standalone review loads neither action. The board correctly caught the resulting misplaced REQs; prevention belongs beside the review instruction that sends an agent into the archive.

## Context

- Original priority: Primary; no P-level severity was supplied.
- Verbatim claim: “Missing: any statement that the archived UR folder stays put.”
- Evidence: `skills/do-work/actions/review-work.md:60` reads archived input and `:412` queues follow-ups, while `skills/do-work/actions/abandon.md:129` and `skills/do-work/actions/capture.md:156` state the closed-folder invariant.
- Surface-cost: Earned. Commit `1323982` records a real misplaced-REQ replay case, and the long-lived addition is a co-located instruction plus one focused contract assertion.

## Detailed Requirements

- Add the invariant immediately beside Step 3's fallback read from `do-work/archive/UR-NNN/input.md`.
- State that reading archived input is context-only and grants no authority to move, reopen, or re-consolidate the archived UR folder.
- Preserve the established review-follow-up model: the new REQ keeps the same `user_request`, carries `review_generated: true`, and goes into `do-work/queue/`.
- Distinguish standalone review of a closed UR from orchestrated review while the UR is still open.
- Make the work/archive path explicit for a review-generated REQ whose UR already lives in `archive/`: move the completed REQ into that existing UR folder in place and never move the UR folder itself.
- Preserve review's existing narrow exception that permits appending post-work review metadata to an archived REQ.
- Add a co-located contract-regression assertion that fails if the archived-input block loses the stays-closed instruction.

## Constraints

- This is documentation/lifecycle hardening, not a claim that queue-kanban behaved incorrectly.
- State the rule at the review consumer that needs it; do not restate it across unrelated actions.
- Do not create a new UR for standalone-review follow-ups. The user confirmed the same-UR/stays-closed model.
- Do not edit downstream consumer repositories; suite update propagation is outside this REQ.

## Dependencies

None. REQ-194 depends on this REQ's canonical lifecycle statement.

## Builder Guidance

Firm intent. Prefer the smallest wording and archive-branch clarification that make the same-UR/stays-closed path unambiguous.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Add a `review_archived_input_block` probe to `_dev/tests/contract-regressions.sh` that extracts `review-work.md` Step 3 and requires the archived UR to stay closed while follow-ups go to the queue.
**Why RED now:** The extracted block names the archive fallback but contains no prohibition on moving or reopening that folder, so the focused assertion fails.
**GREEN when:** The same probe passes because the invariant is co-located in Step 3, the work archive text covers the already-archived review-generated case in place, and the baseline contract suites remain green.
**Validation:** User confirmed the contract assertion and same-UR/stays-closed lifecycle on 2026-08-15. `tdd: false` because the primary change is operating documentation; the named check still supplies finding-closure proof.

## Assets

None.

## Full Context

See `do-work/user-requests/UR-043/input.md` for the complete verbatim request, validated evidence, and batch constraints.

---
*Source: documentation-hardening finding validated through `do-work-toolbox validate-feedback`; original priority Primary.*

## Triage

**Route B** — the closed-UR intent and candidate files are firm, but the exact standalone-review read seam, completed-follow-up archive branch, and co-located contract extraction need a focused trace before Scope is declared.

## Exploration

- `review-work.md` Step 3 is the missing authority seam: it reads `do-work/archive/UR-NNN/input.md` when the UR is closed but does not say that the read is context-only or prohibit relocating the folder.
- Standalone review's producer is otherwise correct: a follow-up retains the same `user_request`, carries `review_generated: true`, and is created in `do-work/queue/`; the existing archived-REQ Review annotation remains the narrow post-work metadata exception.
- `work.md` Step 8 assumes every `user_request` points to an active UR. It needs a more specific success override when both `review_generated: true` and an existing `do-work/archive/UR-NNN/` prove the same-UR/stays-closed shape.
- The override moves only the completed REQ into the existing archived UR folder and bypasses the active-UR closure branch. It does not legitimize generic misplaced REQs and does not alter failure archival.
- The aggregate already has narrow Markdown-block extraction helpers and an existing `work_archive_success_block`; separate assertions at those two seams can reject deletion of any required clause without adding a new harness.

## Scope

**Files I will touch:**
- `skills/do-work/actions/review-work.md` (modify) — state the archived-input context-only/stays-closed boundary and preserve the same-UR queued follow-up model
- `skills/do-work/actions/work.md` (modify) — archive a successful review-generated follow-up directly into its already-archived UR folder without moving that folder
- `_dev/tests/contract-regressions.sh` (modify) — add co-located non-vacuous assertions for both lifecycle seams

**Files I will NOT touch:** `work-reference.md`, capture/abandon, board/forensics behavior, consumer repositories, or any generic archive scanner.

**Acceptance criteria (restated from REQ):**
- [x] The archived-input fallback explicitly grants context only and never authority to move, reopen, or re-consolidate the closed UR folder.
- [x] Standalone-review follow-ups keep the same `user_request`, carry `review_generated: true`, and enter `do-work/queue/`; ordinary orchestrated review remains the open-UR path.
- [x] A successful review-generated REQ whose UR folder already exists in archive moves into that folder in place, and the UR folder itself never moves.
- [x] The narrow archived-REQ Review metadata exception remains intact.
- [x] Focused contract assertions fail if either co-located lifecycle seam loses the required rule.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/review-work.md` (modified) — makes every archived-input fallback context-only, keeps the closed UR stationary, and distinguishes ordinary open-UR orchestration from closed-UR review follow-ups
- `skills/do-work/actions/work.md` (modified) — adds the success-only in-place archive override for review-generated REQs whose UR folder is already archived
- `_dev/tests/contract-regressions.sh` (modified) — extracts the two narrow lifecycle seams and independently asserts their conditions, placement, authority boundary, and producer contracts

**Behavior:** Standalone review can read a closed UR for context and queue same-UR fixes without reopening it. When such a review-generated fix later completes, work places the REQ into the existing archived UR folder and leaves that folder closed and in place; ordinary active-UR archival and failure behavior are unchanged.

## Testing

**Finding-closure RED:** With the initial 11 focused assertions added first, `bash _dev/tests/contract-regressions.sh` exited 1 because Step 3 lacked the context-only/provenance/mode clauses and Step 8 lacked the marker-plus-existing-folder in-place override. Remediation then expanded the final contract to 13 assertions so the archived-fallback boundary is explicitly mode-independent.

**GREEN:**
- `bash -n _dev/tests/contract-regressions.sh` — PASS
- focused archived-input/archive-override assertion probe — PASS, including absence of the false universal orchestrated-open statement
- formal review-remediation mutation probe — PASS: negating `already exists`, deleting `review_generated: true`, or replacing `and` with `or` each fails the complete marker-plus-existing-folder predicate
- `bash _dev/tests/contract-regressions.sh` — PASS after the initial implementation and after the mode-boundary remediation
- `bash _dev/tests/maintainer-verify.sh` — PASS on the final remediated diff, including ShellCheck, the aggregate, both Go modules, and the strict JavaScript lane
- `git diff --check` — PASS

## Qualification

- **Scope:** PASS — `scope-drift.sh` reports the exact three-file Implementation Summary matches Scope; foreign queue edits remain excluded.
- **Mechanical checks:** PASS — `qualify.sh` found complete P-A-U evidence, the exact modified file set, and no debug artifacts.
- **Substance and traceability:** PASS — all seven detailed lifecycle requirements map directly to the two co-located instruction seams and the narrow aggregate contract.
- **Wiring/data flow:** PASS — same-UR review follow-ups flow from the archived-input authority boundary, through the existing queued `review_generated` producer, to the success-only in-place archive override; the final canonical maintainer gate passes.

## Review

**Overall: 99%** | 2026-08-15T13:07:23Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 98% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None. The initial review's sole Important finding—the existence predicate surviving deletion/negation—was closed in the allowed remediation attempt and independently re-verified.
**Minor findings:** None.
**Acceptance:** Pass — archived-input authority and already-archived review-follow-up placement are complete and mutation-locked.
**Suggested testing:** None beyond the canonical maintainer gate and recorded predicate mutations.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Keying the authority rule to use of the archived fallback, rather than review mode, covers both standalone review and later orchestrated review of its generated fix.
- A literal full-predicate contract locks the marker, conjunction, archived-folder existence, and same-UR relation without adding another scanner.

**What didn't:**
- The first wording treated orchestrated review as universally open-UR, overlooking review-generated work whose UR deliberately stays archived.
- The first broad regex required the archive path token but let deletion or negation of `already exists` survive.

**Worth knowing:** The in-place archive override is legitimate only when both `review_generated: true` and the matching archived UR folder already exist. Generic live REQs under closed URs remain anomalies for REQ-194's detector path.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Orientation

**[MAP CHANGED]** Standalone review now has a closed-UR lifecycle from context read through completed follow-up placement: same-UR fixes queue normally, but the archived UR never reopens or moves, and the finished review-generated REQ returns to that existing folder in place.
