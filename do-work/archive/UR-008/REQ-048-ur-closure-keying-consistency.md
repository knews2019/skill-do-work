---
id: REQ-048
title: "Make cleanup Pass 1 and work Step 8 agree on when a UR closes (key on user_request, not the requests: array)"
status: completed
route: B
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T13:23:51Z
completed_at: 2026-07-29T13:34:27Z
commit: da8ec3c
user_request: UR-008
domain: general
prime_files: []
tdd: false
depends_on: []
related: []
batch: deep-review-followups
write_set: [actions/cleanup.md, actions/capture.md]
maintenance: false
---

# UR-Closure Keying Consistency (Cleanup Pass 1 vs Work Step 8)

## What

The two UR-closure readers disagree on what "all REQs resolved" means. `actions/work.md` Step 8 (~:567) scans `queue/`, `working/`, and `archive/` for **all REQs belonging to the UR by `user_request:` frontmatter**. `actions/cleanup.md` Pass 1 reads only the UR's own `requests:` array from `input.md`. Review-spawned and addendum follow-ups carry `user_request: UR-NNN` but are never added to the array — so once the array's REQs are terminal, cleanup Pass 1 read literally archives a UR that Step 8 correctly keeps open, stranding the pending follow-ups' UR reference. This is live today: UR-007's array is `[REQ-032, REQ-033, REQ-034]` (all terminal) while six clarify-derived follow-ups (REQ-051–056, `user_request: UR-007`) pend against it. (Promoted from the prior session's CHECKPOINT.md note; originally observed with REQ-041/042, since resolved into those follow-ups.)

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `crew-members/general.md` + `coding-guardrails.md`, `actions/cleanup.md` Pass 1, `actions/work.md` Step 8 (:583), `actions/work-reference.md` § Terminal-resolved status set (:193–199), `actions/capture.md` Step 5. Approach: (1) ripgrep every reader of the `requests:` array across `actions/`; (2) rewrite cleanup Pass 1 to derive UR membership from a `user_request:` frontmatter scan over the same four locations work Step 8 names, gating on work-reference's terminal-resolved set by reference; (3) document the array in `actions/capture.md` Step 5 as capture-time record only, naming the scan as the closure predicate. `prime_files: []` — no prime to read.
- [x] **[APPLY]:** Two files touched, both in `write_set`. `actions/cleanup.md`: Pass 1 rewritten (condition paragraph + 5 numbered substeps), one Red Flag extended with the cause, one Verification Checklist item added. `actions/capture.md`: one paragraph added after Step 5 substep 4. No other file edited.
- [x] **[UNIFY]:** `git diff --numstat` → `actions/capture.md` (+2/−0), `actions/cleanup.md` (+14/−8); 2 files, no third file in the diff. Verified `actions/cleanup.md` — Pass 1 now scans `user_request:` across `queue/`, `working/`, `archive/` root, `archive/UR-NNN/`; duplicate-location guard text preserved verbatim (archive root + UR folder → flag, leave untouched); terminal-resolved set referenced to `actions/work-reference.md`, not restated or forked; `failed` explicitly named as holding the UR open. Verified `actions/capture.md` — the added paragraph is documentation only; Step 5's four substeps, the templates pointer, and the commit step are unchanged. `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` (repo's only native linter for action prose). No debug artifacts, no TODOs, no stray files: `git status --porcelain` shows exactly ` M actions/capture.md` and ` M actions/cleanup.md` and nothing else (this repo excludes `do-work/` via `.git/info/exclude`, so the REQ's own bookkeeping never appears in the diff). `actions/work.md` and `actions/work-reference.md` untouched, as required.

## Requirements

1. Make cleanup Pass 1 key on the same condition as Step 8: a UR closes only when **every REQ carrying its `user_request:` across `queue/`, `working/`, `archive/` root, and `archive/UR-NNN/`** is terminally resolved. The `requests:` array remains useful as the capture-time record but must not be the closure predicate — per the repo's Closed Enumerations Go Stale rule, state the condition, treat the array as illustrative.
2. Decide and document the array's maintenance story: either follow-up creators (work Step 8's follow-up REQs, capture addenda) append to the originating UR's `requests:` array, or the array is explicitly documented as capture-time-only and every reader is pointed at the `user_request:` scan. Either is acceptable; both at once is not required — pick one and grep for every other reader of `requests:` before finishing.

## Constraints

- Preserve Pass 1's existing duplicate-location guard (REQ found in both archive root and UR folder → flag, leave untouched).
- `failed` stays non-terminal for closure in both readers (it holds a UR open) — don't disturb the terminal-resolved set definition in `actions/work-reference.md`.

## Red-Green Proof
**RED prompt/case:** With today's tree: run cleanup Pass 1's steps literally against UR-007 — its array `[032, 033, 034]` is all-terminal, so Pass 1 archives UR-007 while queued REQ-051–056 still carry `user_request: UR-007`; Step 8's rule (~work.md:567) run on the same state keeps it open.
**Why RED now:** The two readers key on different data; review-spawned follow-ups exist only in the `user_request:` dimension.
**GREEN when:** Both readers evaluate the same predicate; the UR-007 scenario yields "still open" from both; the `requests:` array's role (maintained vs capture-time-only) is stated in its canonical home and consistent at every reader.
**Validation:** User confirmed (approved capture of the checkpoint-noted discrepancy)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** B (Explore then Build)
**Reasoning:** Clear outcome (cleanup Pass 1 must key on the same `user_request:`-scan predicate as work Step 8), but the "grep every other reader of `requests:`" requirement and the two-file coordination (cleanup.md + capture.md, where the array is authored/documented) need an exploration pass before editing — not a blind single-spot change. Route B.
**Complexity indicators:** 2 requirements; a completeness requirement (find every reader of `requests:`); a small design decision (item 2: maintained array vs capture-time-only); two constraints to preserve (Pass 1 duplicate-location guard; `failed` stays non-terminal). Live-today correctness impact (UR-007 would be prematurely archived by this loop's end-of-run cleanup without the fix).
**Rigor:** Standard independent review (main-context) + verify both readers now evaluate the identical predicate and the UR-007 scenario yields "still open" from cleanup.

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Exploration

Sweep command: `rg -n 'requests:|requests\[|requests array|input\.md' actions/` plus a narrowing pass (`rg -n '\brequests\b'`) over `actions/`, `SKILL.md`, `docs/`, and `tools/queue-kanban/*.go`, filtering out `user-requests`/`user_request` noise.

**Every reader of the UR `requests:` array, and what it uses it for:**

| Site | Uses the array as… | In write_set? |
|---|---|---|
| `actions/capture.md:205, 208` (Step 5 substeps 1 + 4) | **Writer** — creates it empty, then fills it with the REQ ids this capture produced. Nothing else ever appends. | yes |
| `actions/cleanup.md:48–49` (Pass 1) | **The UR-closure predicate** — the bug. Iterates the array and closes the UR when those ids are terminal, regardless of any other REQ carrying the UR's `user_request:`. | yes |
| `actions/forensics.md:64–66` (Check 4, Orphaned URs) | **A closure predicate** for a diagnostic Warning ("this UR should have been moved to `archive/`"). Same stale-list read as cleanup's, so it emits a false-positive warning while follow-ups pend. | **no** → `## Discovered Tasks` |
| `actions/verify-requests.md:42, 47` | **Capture-time scope** — enumerates which REQs to grade for coverage against the original input. Correct use: it *should* only see what capture created. Consistent with the chosen story; no change needed. | no (no change needed) |
| `actions/capture-reference.md:142` | **Field shape only** — the `UR input.md` template's `requests: [REQ-020]` line. Carries no semantics. | no (no change needed) |
| `actions/present-work.md:430` (Portfolio mode) | **Display** — "read each UR's `input.md` for the title and request list" over already-archived UR folders. Not a closure decision. | no (no change needed) |

**Non-readers confirmed:** `actions/version.md:150–151` already does it right (archive source reads the UR folder's actual `REQ-*.md` files; active source scans `user_request:` frontmatter across `queue/` and `working/`). `actions/ai-report.md:61`, `actions/review-work.md:60`, `actions/roadmap.md`, `actions/recap.md`, `tools/queue-kanban/` read `input.md` for verbatim input or title, never the array. `actions/work-reference.md`'s Schema Read Contract has no row for `requests` — it is a REQ-field contract, so `actions/capture.md` Step 5 is the array's canonical home and is where the maintenance story now lives.

**Cleanup Pass 1's pre-fix keying:** substep 1 parsed the array; substep 2 looked for each array id in exactly two locations (`archive/UR-NNN/`, `archive/` root) — it never looked in `queue/` or `working/` at all, so a pending REQ was invisible twice over (wrong key *and* wrong location set). Work Step 8 (`actions/work.md:583`) scans four locations by `user_request:`.

## Scope

Files I will touch:

- `actions/cleanup.md` — Pass 1 keying (+ the one Red Flag and one Verification Checklist item that state Pass 1's contract)
- `actions/capture.md` — Step 5, document the `requests:` array's role

Both are the declared `write_set`. Nothing else. `actions/forensics.md` needs the same alignment but is out of `write_set` → recorded as a Discovered Task rather than edited.

## Implementation Summary

**Files changed:**

- `actions/cleanup.md` (modified) — Pass 1 rewritten. A new "The closure condition, stated — not a stored list" paragraph states the predicate and why the array can't be it; substep 1 now collects the UR's REQs by reading `user_request:` from every `REQ-*.md` in `do-work/queue/`, `do-work/working/`, `do-work/archive/` root (non-recursive), and `do-work/archive/UR-NNN/` — the same four locations `actions/work.md` Step 8 names; substep 2 gates on `actions/work-reference.md`'s terminal-resolved set **by reference** and spells out that anything outside it holds the UR open, `failed` included; substeps 3–4 now say "collected REQs" instead of "REQs in the array"; a new report-only substep 5 flags an array id found in none of the four locations as a missing file that must never wedge closure. The duplicate-location guard survives verbatim inside substep 2. One Red Flag now names the cause under its symptom, and one Verification Checklist item asserts membership came from the scan.
- `actions/capture.md` (modified) — one paragraph after Step 5 substep 4 documenting `requests:` as the capture-time record only and never the closure predicate, listing what carries `user_request:` without being appended (review-spawned follow-ups, addendum REQs, clarify-derived REQs), naming the four-location `user_request:` scan as the answer to "is this UR finished?", and separating the array's legitimate audience (`actions/verify-requests.md`-style coverage grading) from readers deciding closure.

**What was done:** the two UR-closure readers now evaluate one condition instead of two different data sources. Per the repo's Closed Enumerations Go Stale rule, cleanup Pass 1 states the *condition* (every REQ carrying this `user_request:`, wherever it sits, is terminally resolved) and treats `input.md`'s `requests:` array as the capture-time record it actually is — a list nothing maintains after capture, which is exactly why keying on it archived URs out from under their follow-ups. The terminal-resolved set is cited, never restated, so `failed`'s non-terminal-for-closure status has a single definition in `actions/work-reference.md`. Verified against the live tree: UR-007's array is `[REQ-032, REQ-033, REQ-034]` (all `completed`), but the scan finds 17 members carrying `user_request: UR-007` — 11 `completed` in `archive/` root and REQ-051 through REQ-056 `pending` in `queue/` — so the rewritten Pass 1 reports `UR-007 still open (11/17 REQs complete)` where the old text would have archived it.

## Decisions

- **D-01 — Requirement 2: the array is documented as capture-time-only; follow-up creators do NOT append to it.** Took the recommended option. Value: zero new machinery, and it matches what the array already *is* — a record of one capture event, written once at `actions/capture.md` Step 5. Risk of the alternative: making work Step 8, capture's addendum path, and clarify each append to a sibling UR's `input.md` adds three write sites to a field no consumer would then be allowed to trust anyway (any single missed append silently re-creates this exact bug), and it means mutating a UR's `input.md` long after capture — precisely the kind of hand-maintained enumeration the Closed Enumerations Go Stale rule exists to stop. Reversible: if a maintained array is ever wanted, the closure predicate is already independent of it, so adding appends later breaks nothing.
- **D-02 — Membership comes from the `user_request:` scan alone; the array contributes nothing to the predicate.** Considered union-ing the array's ids into the membership set as a belt-and-braces measure and rejected it: an array id with no file anywhere would then hold its UR open forever with no way to clear it, and the GREEN criterion is that *both* readers evaluate the identical predicate — a union is not what work Step 8 evaluates. The array's one remaining use in Pass 1 is the report-only cross-check (substep 5), which cannot affect the outcome.
- **D-03 — Existing report strings left alone.** Substep 4 still prints `UR-NNN still open (X/Y REQs complete)` rather than "resolved", because the Reporting section's example line quotes that exact wording; renaming it would have pulled an unrelated section into the diff for no behavioral gain.

## Discovered Tasks

- **[normal] `actions/forensics.md` Check 4 (Orphaned URs, lines 63–66) is the third reader of the closure predicate and still keys on the `requests:` array.** It reads the array, checks whether those ids are all in `archive/`, and warns "this UR should have been moved to `archive/`" — so for UR-007 today it emits a false-positive Warning telling the user to archive a UR that correctly has six pending REQs. Fix is the same shape as the Pass 1 rewrite: derive membership from a `user_request:` scan across `do-work/queue/`, `do-work/working/`, `do-work/archive/` root, and `do-work/archive/UR-NNN/`, gate on `actions/work-reference.md`'s terminal-resolved set. Out of this REQ's `write_set`.
- **[normal] `actions/work.md` Step 8's archive table (line 583) counts `failed` as "finished" for UR closure, contradicting `actions/work-reference.md`'s Terminal-resolved status set.** That section (line 195) names "`actions/work.md` Step 8's UR-final check" as one of its honoring readers and (line 199) states `failed` "stays outside this set … a UR with a `failed` REQ needs those follow-ups before it can close." Step 8's row instead lists `failed` in its finished set and closes the UR, leaving the failed REQ at `archive/` root. After this REQ, cleanup Pass 1 holds such a UR open while work Step 8 would close it — the two readers agree on the key and the locations but not on this one status. Narrow in practice (Step 8's failure classification normally spawns a `pending` follow-up that holds the UR open anyway), but real when no follow-up is created (cycle detection skips it). Fix: drop `failed` from the finished set in that table row and cite the terminal-resolved set. Out of this REQ's `write_set` (the orchestrator held `actions/work.md` out of scope).

## Review

**Acceptance: Pass — overall ~95%.** Main-context review against the full 2-file diff + verification both readers evaluate the same predicate.

**Requirements (both met):**
1. cleanup Pass 1 rewritten: collect UR members by scanning `user_request:` across all four locations (queue/, working/, archive/ root, archive/UR-NNN/), gate on work-reference's terminal-resolved set BY REFERENCE (not forked); `failed` explicitly holds the UR open; duplicate-location guard preserved verbatim; new report-only cross-check (substep 5) flags a stale `requests:` entry found nowhere without ever holding the UR open. Red Flag + Verification Checklist item added.
2. capture.md documents `requests:` as the capture-time record only, names the `user_request:` scan as the closure predicate, and points verify-requests.md as a legitimate array reader (grades capture coverage). Chose the capture-time-only option (D-01) — zero new machinery.

**Verified:** the UR-007 scenario (array `[032,033,034]` all-terminal, REQ-051–056 pending) now yields "still open" from cleanup Pass 1, matching Step 8's keying + locations. D-02 correctly rejected union-ing the array into the predicate (a stale id would wedge closure forever). qualify + contract-regressions pass.

**Discovered tasks (both queued `pending-answers`):** REQ-058 (forensics.md Check 4 — same array-keying bug, live false-positive on UR-007); REQ-059 (work.md Step 8's table counts `failed` as finished, contradicting the terminal-resolved set — the residual `failed` disagreement, needs a semantics decision).

No Important/Critical findings in scope. No inline follow-ups.

## Lessons Learned
**What worked:** Sweeping ALL readers of `requests:` before editing surfaced two more — forensics Check 4 (same predicate bug) and work.md Step 8 (`failed` contradiction) — turning a two-file fix into a complete map of the closure-predicate surface.
**Worth knowing:** The `requests:` array is written once at capture and never maintained — any "is this UR done?" reader keying on it is a latent premature-archive bug. Four readers exist: capture (writer), cleanup Pass 1 (fixed), forensics Check 4 (→ REQ-058), verify-requests (correct — wants capture-time scope). Separately, work.md Step 8 and work-reference.md disagree on whether `failed` closes a UR (→ REQ-059).

## Orientation
UR closure now has one stated predicate — every REQ carrying `user_request: UR-NNN` (scanned across queue/working/archive) is terminally resolved — shared by cleanup Pass 1 and work Step 8, with the `requests:` array demoted to a capture-time record (documented in capture.md). Lives in `actions/cleanup.md` Pass 1 + `actions/capture.md` Step 5. Two readers still to align (forensics Check 4, work.md Step 8's `failed` row) are queued as REQ-058/059. No map change — hardens UR lifecycle consistency.
