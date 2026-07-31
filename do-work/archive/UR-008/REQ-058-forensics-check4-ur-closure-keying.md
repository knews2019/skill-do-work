---
id: REQ-058
title: "forensics.md Check 4 (Orphaned URs) keys UR closure on the requests: array — same stale-list bug as cleanup Pass 1"
status: completed
route: B
claimed_at: 2026-07-29T13:55:51Z
commit: 1c3b662
created_at: 2026-07-29T13:31:00Z
status_changed_at: 2026-07-29T13:51:27Z
user_request: UR-008
addendum_to: REQ-048
domain: general
prime_files: []
tdd: false
depends_on: []
write_set: [actions/forensics.md]
maintenance: false
review_generated: false
---

# forensics.md Check 4 keys UR closure on the stale requests: array

## What
Discovered during REQ-048. `actions/forensics.md` Check 4 (Orphaned URs, ~lines 63–66) is the third reader of the UR-closure predicate and still keys on the UR's `requests:` array. It warns "this UR should have been moved to archive/" when the array's ids are all archived — so for UR-007 today (array `[032,033,034]` all terminal, but REQ-051–056 pending with `user_request: UR-007`) it emits a **false-positive Warning** telling the user to archive a UR that correctly has pending follow-ups. Same shape as the REQ-048 cleanup Pass 1 fix.

## Open Questions
- [x] Fix forensics Check 4 to key on the `user_request:` scan (across `queue/`, `working/`, `archive/` root, `archive/UR-NNN/`, gating on `actions/work-reference.md`'s terminal-resolved set) instead of the `requests:` array? → Confirmed: yes — apply the same one-reader fix REQ-048 applied to cleanup Pass 1, removing the live false-positive warning against UR-007.

## Full Context
Discovered task from REQ-048. See `do-work/archive/REQ-048-ur-closure-keying-consistency.md` → `## Discovered Tasks`.

## Triage

**Route:** B
**Reasoning:** `actions/forensics.md` Check 4 rewrite to the `user_request:` scan across the four locations (same shape as REQ-048's cleanup Pass 1 fix), gating on work-reference's terminal-resolved set. Answered yes via clarify. Route B.
**Rigor:** Standard main-context review (part of the parallel disjoint-write_set batch 051/052/054/057/058; single-file, no spec-cluster overlap).

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Exploration

**Check 4's keying before this REQ** (`actions/forensics.md` :63-66) was four lines: list UR folders in `do-work/user-requests/`, read the `requests` array from `input.md` frontmatter, check whether **all ids in that array** exist somewhere under `do-work/archive/` (root or `UR-NNN/`), and Warn if so. Two independent defects in one predicate: membership came from the capture-time array, and resolution was tested by **file location** (does the file sit in `archive/`) rather than by `status`.

**The predicate I mirrored** is `actions/cleanup.md` Pass 1 (:46-64, rewritten by REQ-048). Its shape:

1. Membership is *derived*, not stored — scan the `user_request` field of every `REQ-*.md` across four locations: `do-work/queue/`, `do-work/working/`, `do-work/archive/` root (non-recursive), and `do-work/archive/UR-NNN/`.
2. Resolution is a **status** test against `actions/work-reference.md`'s Schema Read Contract → **Terminal-resolved status set** (:193-199), referenced by pointer. Any status outside it holds the UR open, `failed` included.
3. The `requests:` array is capture-time-only (`actions/capture.md` Step 5); review-spawned follow-ups, addendum REQs, and clarify-derived REQs all carry `user_request:` without ever being appended to it.

**The live false positive, verified before the fix.** `do-work/user-requests/UR-007/input.md` has `requests: [REQ-032, REQ-033, REQ-034]`, all three terminal in the archive — so old Check 4 fired its Warning. Meanwhile six REQs carry `user_request: UR-007` and are *not* resolved: REQ-053/055/056 (`status: pending`, in `do-work/queue/`) and REQ-051/052/054 (`status: claimed`, in `do-work/working/`). The new predicate collects all nine, finds six outside the terminal-resolved set, and leaves UR-007 open — no finding. (REQ-058 itself is UR-008, unaffected either way.)

## Scope

**Files I will touch:** `actions/forensics.md`

Nothing else. `docs/forensics-guide.md` describes this check in one summary-table row but is outside the declared `write_set` — see `## Discovered Tasks`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files: []` — read `crew-members/general.md` + `crew-members/coding-guardrails.md`. Approach: read cleanup Pass 1 and work-reference's Terminal-resolved status set first, confirm the UR-007 false positive against the real tree (`grep user_request: UR-007` across queue/working + statuses), then rewrite Check 4's body to cleanup Pass 1's derive-membership-then-status-test shape, pointing at work-reference for the set instead of forking it. Keying only — no new finding types, no new severities.
- [x] **[APPLY]:** One file edited, `actions/forensics.md` Check 4 (Orphaned URs). Replaced the array-keyed 4-bullet body with a stated closure condition plus 3 numbered steps: (1) collect the UR's REQs by scanning `user_request` across the four locations, (2) test each against work-reference's Terminal-resolved status set by pointer (`failed` holds open, stated), (3) Warn only when all collected REQs are terminally resolved and the UR is still in `user-requests/`. The array is explicitly named as *not* the predicate, with the reason (capture-time-only) and the co-reader pointer to `actions/cleanup.md` Pass 1 so the two can't drift. Check title, severity (Warning), and diagnostic intent unchanged.
- [x] **[UNIFY]:** `git diff --stat` — the only entry attributable to this REQ is `actions/forensics.md` (+15/-4, one hunk at Check 4); the three other entries in the working tree (`_dev/tests/contract-regressions.sh`, `actions/capture-reference.md`, `tools/queue-kanban/generate_test.go`) belong to the parallel builders in this batch, not to me. Prose-only change to a Markdown action file — no code, no debug artifacts, no build outputs, no `console.log`/`TODO`/scratch text (full diff read line-by-line above). Linters: `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` Per-file review: `actions/forensics.md` — Check 4 body only; Checks 1-3 and 5-12, Output Format, Red Flags, and Verification Checklist all byte-identical; the two file-path citations added (`actions/cleanup.md`, `actions/capture.md`, `actions/work-reference.md`) are path-form per the cross-reference convention and all three files exist.

## Implementation Summary

**Files changed:**
- `actions/forensics.md` (modified) — Check 4 (Orphaned URs) re-keyed

**What was done**

Check 4 no longer reads the UR's `requests:` array as its closure predicate. It now derives UR membership the way `actions/cleanup.md` Pass 1 does — scanning the `user_request` frontmatter field of every `REQ-*.md` in `do-work/queue/`, `do-work/working/`, `do-work/archive/` root (non-recursive), and `do-work/archive/UR-NNN/` — and tests each collected REQ's `status` against `actions/work-reference.md`'s Schema Read Contract → Terminal-resolved status set, referenced by pointer rather than restated. Any status outside that set holds the UR open, with `failed` called out explicitly. The Warning still fires for exactly the case the check was written to catch (a UR that should have been archived and wasn't) at the same severity; what changed is which URs qualify.

This also replaces the old *location* test ("do all the array's REQs exist under `archive/`") with a *status* test, matching Pass 1: archival location is Pass 0/Pass 2's business, and a REQ's `status` is what says whether it is resolved.

Effect on the live tree: UR-007 stops false-positiving. Its array `[REQ-032, REQ-033, REQ-034]` is all-terminal, but REQ-051/052/054 (`claimed`) and REQ-053/055/056 (`pending`) carry `user_request: UR-007`, so the UR is correctly reported as still open.

## Decisions

- **DECIDE & STATE — pointer-only for the terminal set, no parenthetical value list.** cleanup Pass 1 restates `(completed, completed-with-issues, or cancelled — see …)` alongside its pointer. I cite work-reference without the values, matching Check 11's existing phrasing in this same file ("that table is the canonical vocabulary … do not re-enumerate it here"). One less copy of the set to drift, and it satisfies the REQ's "reference by pointer, don't fork it" directly. `failed` is still named inline because it is the *counter*-intuitive member (terminal but non-closing) and naming it is what prevents the next reader from re-introducing the bug.
- **DECIDE & STATE — did not port cleanup Pass 1's step 5 (report-only array cross-check).** Pass 1 additionally reports `⚠ REQ-NNN listed in UR-NNN's requests: array but found nowhere`. That is a *new finding type* for forensics, not part of the keying fix, and the REQ scopes this to "no behavior beyond the keying fix." Left out deliberately; it would be a clean standalone REQ if wanted (a missing REQ file is a legitimate forensic finding).
- **DECIDE & STATE — did not port Pass 1's duplicate-REQ guard** (same id in both `archive/` root and `archive/UR-NNN/`). It exists in Pass 1 to stop a *move*; forensics is read-only, so it has nothing to protect.

## Discovered Tasks

- `docs/forensics-guide.md` (~:14) summarizes this check as "UR folders in `user-requests/` where all REQs are **archived** but UR wasn't moved" — a location-phrased summary of what is now a status-phrased predicate. It was never array-keyed, so it isn't wrong in the way Check 4 was, but a terminally-resolved REQ still sitting in `queue/` (e.g. a `cancelled` REQ awaiting Pass 0) now counts toward closure while the guide's wording implies it wouldn't. One-clause reword ("all REQs are terminally resolved"). Outside this REQ's `write_set`; not fixed inline.
