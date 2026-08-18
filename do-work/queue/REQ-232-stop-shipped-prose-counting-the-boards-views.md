---
id: REQ-232
title: "Review fix: stop shipped prose from counting the board's views"
status: pending
domain: general
created_at: 2026-08-18T01:17:41Z
user_request: UR-051
addendum_to: REQ-227
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: shipped-prose-hand-counts-board-views
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

- [ ] `skills/do-work-board/actions/board.md:3`: "render this repo's `do-work/` queue as a Kanban board + completion calendar" — a two-item list of what is now five views.
- [ ] `skills/do-work-board/actions/board.md:14`: "linked from the Board/Calendar view toggle" — names two of five.
- [ ] `skills/do-work-board/actions/board.md:93`: "The served board's **Testing** view (a third view next to Board / Calendar)" — an ordinal that was already wrong before REQ-227 and is wronger after it.

## Red-Green Proof

**RED prompt/case:** Grep `skills/do-work-board/actions/board.md` for the three instances above.
**Why RED now:** All three are present and all three misdescribe a board that has five views.
**GREEN when:** The named surfaces are absent — the file describes the board without counting or part-listing its views — and `bash _dev/tests/maintainer-verify.sh` still exits 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**, whose deletion branch is the one that applies here.
