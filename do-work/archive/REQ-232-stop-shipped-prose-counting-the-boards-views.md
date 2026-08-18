---
id: REQ-232
title: "Review fix: stop shipped prose from counting the board's views"
status: completed
completed_at: 2026-08-18T01:50:46Z
claimed_at: 2026-08-18T01:46:45Z
domain: general
created_at: 2026-08-18T01:17:41Z
user_request: UR-051
addendum_to: REQ-227
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: shipped-prose-hand-counts-board-views
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-18T01:44:00Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - cross-route regression gates
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
maintenance: true
write_set:
- skills/do-work-board/actions/board.md
---

# Review Fix: Stop Shipped Prose from Counting the Board's Views

## What

`skills/do-work-board/actions/board.md` describes the board by listing its views and, in one place, by numbering them. Both go stale every time a view is added, and both are stale now. Replace the counting with a description keyed on what the board does, so there is nothing to keep in step.

Done means the class cannot recur: after this REQ, adding a sixth view requires no edit to `board.md` at all. Incrementing "third" to "fifth" and appending "Timeline" to the two lists would leave the next view to break it again — that is the same fix this file has silently needed twice already.

## Context

Found by REQ-227's mandatory restatement sweep. REQ-227 added the Timeline as a fifth view; the file already said "third" when there were four, so the drift predates this REQ and was simply made larger by it.

`CLAUDE.md` § *State conditions, not lists* names this exact failure — when a rule or description applies "whenever X", key it on the condition, because a hand-maintained list goes stale as the set grows. `_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale* records the same lesson from four independent defects, and REQ-225 closed a fifth instance of it in the shipped shell guide two REQs ago.

Deliberately **not** fixed inside REQ-227: that REQ's constraints scope it to adding a view additively, its write set names the tool's own files, and the fix here is a rewrite of a sibling package's action file rather than an increment.

## Requirements

- Remove the hand-maintained counts and lists of views from `skills/do-work-board/actions/board.md`, describing the board's capability instead of enumerating its tabs. Where a specific view genuinely has to be named — the Testing view owns a write surface and an API, so it does — name that view for its own reason rather than by its position among the others.
- Do not change what the file says about the write surfaces. `CLAUDE.md` § *Kanban Board Write Surfaces* is the canonical count of those, this file's Step 6 is the Testing view's own contract, and neither is in question here.
- No tool behavior changes. This REQ is documentation only.
- Keep `_dev/tests/shipped-package-reference-contract.sh` green — this is a shipped file, so it may not cite `_dev/` paths.

## Instances

- [x] `skills/do-work-board/actions/board.md:3`: "render this repo's `do-work/` queue as a Kanban board + completion calendar" — a two-item list of what is now five views.
- [x] `skills/do-work-board/actions/board.md:14`: "linked from the Board/Calendar view toggle" — names two of five.
- [x] `skills/do-work-board/actions/board.md:93`: "The served board's **Testing** view (a third view next to Board / Calendar)" — an ordinal that was already wrong before REQ-227 and is wronger after it.

## Red-Green Proof

**RED prompt/case:** Grep `skills/do-work-board/actions/board.md` for the three instances above.
**Why RED now:** All three are present and all three misdescribe a board that has five views.
**GREEN when:** The named surfaces are absent — the file describes the board without counting or part-listing its views — and `bash _dev/tests/maintainer-verify.sh` still exits 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**, whose deletion branch is the one that applies here.

---

## Triage

**Route: B** - Medium

**Reasoning:** The three instances are named with line numbers and the intent is stated, but what each sentence should become had to be worked out from what the file actually needs to say — and the whole file had to be swept for instances the finding did not list.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The three named instances**, read in context rather than by line number alone:

- `:3` — the action's one-line summary, describing the tool's output as "a Kanban board + completion calendar". A two-item list where there are now five views.
- `:14` — a When-to-Use bullet routing the reader to the Testing view "linked from the Board/Calendar view toggle". Names two of five to describe a switcher that holds all of them.
- `:93` — Step 6's opening, "a third view next to Board / Calendar". An ordinal, and one that was already wrong before REQ-227 added the fifth.

**What must not move.** `:5` is the write-surface paragraph — `CLAUDE.md` § *Kanban Board Write Surfaces* is the canonical count of those and Step 6 is the Testing view's own contract. Neither is in question, and the REQ's requirements say so explicitly.

**One view genuinely has to be named.** The Testing view owns the only write surface and the only server API on the board, so Step 6 is about that view specifically. Naming it for *that* rather than for its position is what keeps the sentence true when a sixth view lands.

**Swept beyond the three.** Grepped every shipped markdown file for view ordinals and for `Board/Calendar` pairings; the three named instances were the only ones. No other shipped file enumerates the board's views.

*Explored by work action (inline, serial mode)*

## Scope

**Files I will touch:**
- `skills/do-work-board/actions/board.md` (modify) — the three instances.

**Files I will NOT touch:** `CLAUDE.md` (the write-surface count is canonical and not in question), `web/template.html` (the switcher itself is the source of truth for which views exist), any test.

**Acceptance criteria (restated from REQ):**
- [ ] The counts and part-lists of views are gone; the file describes the board's capability instead of enumerating its tabs.
- [ ] The Testing view is named for its own reason, not by position.
- [ ] The write-surface paragraph is unchanged.
- [ ] No tool behavior changes.
- [ ] `_dev/tests/shipped-package-reference-contract.sh` stays green.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/actions/board.md` (modified)

**What was done:** Replaced three hand-maintained descriptions of which views the board has with descriptions of what the board does. The summary line now names the switcher and what it covers rather than listing two tabs; the Testing-view route now names the switcher rather than two of its entries; and Step 6 identifies the Testing view by the property that makes Step 6 about it — it owns the board's only write surface and its only server API — rather than by its position among the others. Three lines changed, nothing else.

The `## Instances` checklist above is closed: `:3` now reads "an interactive HTML board — one page whose view switcher covers the queue's current state, its history, and its timing"; `:14` now reads "reachable from the board's view switcher"; `:93` now reads "the only one that writes anything, and the only one with a server API behind it".

## Testing

**Tests run:** `bash _dev/tests/shipped-package-reference-contract.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing — both exit 0

**Red-green validation:** the REQ's captured RED is a grep for the three named surfaces, and its GREEN is their absence — the Finding-Closure Ratchet's deletion branch.

- RED: `a third view next to Board / Calendar`, `Board/Calendar view toggle`, and `Kanban board + completion calendar` were all present.
- GREEN: all three absent. Swept wider than the finding listed — grepped every shipped markdown file for view ordinals (`a (second|third|fourth|fifth) view`) and `Board/Calendar` pairings: none remain anywhere under `skills/`.
- The write-surface paragraph appears in the diff only as context; `git diff --unified=0` carries exactly three `+`/`-` pairs and none of them touches it.

**New tests added:** none. The fix is a deletion, and its GREEN is the absence of the named surfaces — the ratchet's own deletion branch. A regex guard against re-introducing view ordinals in prose would be exactly the "add an instruction to fix a drift" move `crew-members/maintenance.md` § 1 says to reach for last, and it would misfire on any prose legitimately naming a view.

**Existing tests updated (cross-REQ impact):** none.

## Lessons Learned

**What worked:** Asking what each sentence was *for* before rewriting it. Two of the three were describing a mechanism — the view switcher — by listing some of its contents, and naming the mechanism made them shorter as well as durable. The third was harder, because Step 6 genuinely is about one specific view; the answer was to name it by the property that makes the step exist, which is a better sentence than the ordinal was even while the ordinal was still accurate.

**What didn't:** Nothing failed. Worth recording that the tempting fix was available and wrong: incrementing "third" to "fifth" and appending "Timeline" to the two lists would have satisfied this REQ's title and left the sixth view to break it again. The REQ's own What section had to rule it out in so many words, which is a sign of how natural the increment feels.

**Worth knowing:** This was the second instance of the same rule closed in one session — REQ-225 moved a shell rule out of one script's section into a shared condition, and this moved a description of the board out of a list of its tabs. Both were found the same way: something new landed, and a sweep asked what else described the thing that changed. Worth keeping as a habit rather than as a rule, because the sweep is cheap and the drift is silent.

## Review

**Overall: 96%** | 2026-08-18T01:50:46Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None. The wider sweep for the same pattern across every shipped markdown file came back empty, so there is no second instance to queue.

**Minor findings:** 1 (report only)
- Nothing prevents a future view ordinal from being written into this file again. The structural fix is that there is no longer a sentence inviting one, which is the honest limit of a deletion; a prose regex guard would misfire on legitimate mentions and is the addition `maintenance.md` § 1 says to reach for last.

**Acceptance:** Pass — all four restated criteria verified: the three named surfaces are gone and the wider sweep is clean, the Testing view is named by its write surface, the write-surface paragraph carries no `+`/`-` diff line, and both the shipped reference contract and the full baseline exit 0.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*
